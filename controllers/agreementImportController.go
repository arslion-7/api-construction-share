package controllers

import (
	"encoding/csv"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Import of goşmaça şertnama (additional agreement) CSVs produced by
// migrations/g2_extract from a legacy paylynew dump. The CSV has no direct
// registry link, so each agreement is matched to a registry through the
// shareholder name in `paychy` (or familiya/ady/otchestvo when present).

var (
	reISODate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)
	// "PGGŞ 1-2/9055" / "PGGŞ\n1-1/1776" inside sene_sered_hasaba_alnan
	rePGGSNumber = regexp.MustCompile(`PGG[ŞS]\s*\n?\s*([^\s,]+)`)
	// "Hasaba alnan: 17.05.2021"
	reHasabaAlnan = regexp.MustCompile(`Hasaba alnan:\s*(\d{2}\.\d{2}\.\d{4})`)
	// "Seredilen: 09.04.2021"
	reSeredilen = regexp.MustCompile(`Seredilen:\s*(\d{2}\.\d{2}\.\d{4})`)
	// leading org-type words in paychy, e.g. "Raýat Babaýew ..." / "Telekeçi ..."
	rePaychyPrefix = regexp.MustCompile(`^(Raýat|Rayat|Telekeçi|Telekeci|Hususy kärhana|Hojalyk jemgyýeti|Döwlet kärhanasy)[:\s]+`)
)

func parseAnyDate(s string) *time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if reISODate.MatchString(s) {
		if t, err := time.Parse("2006-01-02", s[:10]); err == nil {
			return &t
		}
	}
	return parseDMY(s)
}

// agreementImportRow is one CSV record keyed by header name
type agreementImportRow map[string]string

func (r agreementImportRow) get(key string) string {
	return strings.TrimSpace(r[key])
}

// shareholderKey normalizes a person/org name for matching
func shareholderKey(s string) string {
	s = normSpace(s)
	s = rePaychyPrefix.ReplaceAllString(s, "")
	s = strings.Trim(s, ` ",:«»`)
	return strings.ToLower(normSpace(s))
}

// ImportAgreementsCSV handles POST multipart upload of a g2_extract CSV and
// creates additional agreements matched to registries by shareholder name.
// Idempotent via old_g2_tb: rows already imported are skipped.
func ImportAgreementsCSV(c *gin.Context) {
	if _, exists := c.Get("user"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file is required (field 'file')"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse CSV", "details": err.Error()})
		return
	}
	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV has no data rows"})
		return
	}

	header := records[0]
	colIdx := map[string]int{}
	for i, name := range header {
		colIdx[strings.TrimSpace(name)] = i
	}
	for _, required := range []string{"t_b", "paychy", "sebabi"} {
		if _, ok := colIdx[required]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CSV is missing required column: " + required})
			return
		}
	}

	rows := make([]agreementImportRow, 0, len(records)-1)
	for _, rec := range records[1:] {
		row := agreementImportRow{}
		for name, idx := range colIdx {
			if idx < len(rec) {
				row[name] = rec[idx]
			}
		}
		rows = append(rows, row)
	}

	// registryID lookup by normalized shareholder name (org_name and head name)
	type regShareholder struct {
		ID           uint
		OrgName      *string
		HeadFullName *string
	}
	var regs []regShareholder
	if err := initializers.DB.Model(&models.Registry{}).
		Select("registries.id, shareholders.org_name, shareholders.head_full_name").
		Joins("JOIN shareholders ON shareholders.id = registries.shareholder_id").
		Scan(&regs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load registries", "details": err.Error()})
		return
	}
	registriesByName := map[string][]uint{}
	addKey := func(name *string, id uint) {
		if name == nil {
			return
		}
		key := shareholderKey(*name)
		if key == "" {
			return
		}
		ids := registriesByName[key]
		for _, existing := range ids {
			if existing == id {
				return
			}
		}
		registriesByName[key] = append(ids, id)
	}
	for _, r := range regs {
		addKey(r.OrgName, r.ID)
		addKey(r.HeadFullName, r.ID)
	}

	// already imported g2 ids
	var importedTBs []int
	initializers.DB.Model(&models.AdditionalAgreement{}).
		Where("old_g2_tb IS NOT NULL").Pluck("old_g2_tb", &importedTBs)
	alreadyImported := map[int]bool{}
	for _, tb := range importedTBs {
		alreadyImported[tb] = true
	}

	type problemRow struct {
		TB          int    `json:"t_b"`
		Paychy      string `json:"paychy"`
		RegistryIDs []uint `json:"registry_ids,omitempty"`
	}
	var unmatched, ambiguous []problemRow
	imported := 0
	skippedExisting := 0

	err = initializers.DB.Transaction(func(tx *gorm.DB) error {
		for _, row := range rows {
			tb, err := strconv.Atoi(row.get("t_b"))
			if err != nil {
				continue
			}
			if alreadyImported[tb] {
				skippedExisting++
				continue
			}

			// shareholder name: prefer structured familiya/ady/otchestvo
			name := normSpace(strings.TrimSpace(
				row.get("familiya") + " " + row.get("ady") + " " + row.get("otchestvo")))
			if name == "" {
				name = row.get("gurama_ady")
			}
			if name == "" {
				name = row.get("paychy")
			}
			key := shareholderKey(name)

			regIDs := registriesByName[key]
			switch {
			case len(regIDs) == 0:
				unmatched = append(unmatched, problemRow{TB: tb, Paychy: row.get("paychy")})
				continue
			case len(regIDs) > 1:
				ambiguous = append(ambiguous, problemRow{TB: tb, Paychy: row.get("paychy"), RegistryIDs: regIDs})
				continue
			}

			seredText := row.get("sene_sered_hasaba_alnan")

			number := row.get("hasaba_alnan_belgisi")
			if number == "" {
				if m := rePGGSNumber.FindStringSubmatch(seredText); m != nil {
					number = m[1]
				}
			}

			date := parseAnyDate(row.get("hasaba_alnan_date"))
			if date == nil {
				if m := reHasabaAlnan.FindStringSubmatch(seredText); m != nil {
					date = parseDMY(m[1])
				}
			}

			reason := row.get("sebabi_esasy")
			if reason == "" {
				reason = row.get("sebabi")
			}
			if r := []rune(reason); len(r) > 255 {
				reason = string(r[:255])
			}

			var infoParts []string
			if p := row.get("paychy"); p != "" {
				infoParts = append(infoParts, "Paýçy: "+p)
			}
			if s := row.get("sebabi"); s != "" && s != reason {
				infoParts = append(infoParts, "Sebäbi: "+s)
			}
			if m := row.get("min_hat"); m != "" {
				infoParts = append(infoParts, "Min hat: "+normSpace(m))
			}
			seredilen := row.get("seredilen_date")
			if seredilen == "" {
				if m := reSeredilen.FindStringSubmatch(seredText); m != nil {
					seredilen = m[1]
				}
			}
			if seredilen != "" {
				infoParts = append(infoParts, "Seredilen: "+seredilen)
			}
			if l := row.get("login"); l != "" {
				infoParts = append(infoParts, "Login: "+l)
			}

			tbCopy := tb
			agreement := models.AdditionalAgreement{
				RegistryID:      regIDs[0],
				AgreementNumber: number,
				AgreementDate:   date,
				Reason:          reason,
				AdditionalInfo:  strings.Join(infoParts, "\n"),
				OldG2TB:         &tbCopy,
				FromOldRegistry: true,
			}
			if err := tx.Create(&agreement).Error; err != nil {
				return err
			}
			imported++
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import agreements", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":          "Agreements imported",
		"total":            len(rows),
		"imported":         imported,
		"skipped_existing": skippedExisting,
		"unmatched":        unmatched,
		"ambiguous":        ambiguous,
	})
}

// RollbackImportedAgreements deletes every agreement created by the CSV import
// (old_g2_tb set). Manually created agreements are untouched.
func RollbackImportedAgreements(c *gin.Context) {
	if _, exists := c.Get("user"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	res := initializers.DB.Unscoped().
		Where("old_g2_tb IS NOT NULL").
		Delete(&models.AdditionalAgreement{})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rollback imported agreements", "details": res.Error.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "Imported agreements removed",
		"deleted": res.RowsAffected,
	})
}
