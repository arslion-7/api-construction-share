package controllers

import (
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

// ============================================================================
// Parsers for the legacy mainpayly free-text formats.
// Formats were validated against the full 18.06.2021 dump (8656 rows);
// coverage per field is 92-100%. Anything unparsed stays in the old_* columns.
// ============================================================================

var (
	reDMY = regexp.MustCompile(`(\d{2})\.(\d{2})\.(\d{4})`)
	// "07.11.2019\n№13/2145\n08.11.2019\n(6 sany\n2145-tapgyr)"
	reMinHat = regexp.MustCompile(`(?s)^\s*(\d{2}\.\d{2}\.\d{4})\s*\n\s*№\s*(.*?)\s*\n\s*(\d{2}\.\d{2}\.\d{4})\s*\n\s*\(\s*(\d+)\s*sany\s*\n?\s*(\d+)\s*-\s*tapgyr\s*\)\s*$`)
	// "PGGŞ №7-38-3/12\n13.09.2019\nAşgabat ş." / "№21/19\n05.04.2019\nAşgabat ş." / "28/19\n06.06.2019"
	reContract = regexp.MustCompile(`(?s)^\s*(?:PGG[ŞS]\s*)?[N№]?\s*([^\n]*?)\s*\n\s*(\d{2}\.\d{2}\.\d{4})?\s*(?:\n\s*(.*))?\s*$`)
	// "91,06 ga" / "4050,0 m2" / "5,8 ga 25 ýyllyk möhlet bilen berlen"
	reMeydan = regexp.MustCompile(`^\s*([\d\s]+(?:[.,]\d+)?)\s*(ga|m2|m²)\s*(.*)$`)
	// "Baş potratçy:\n№25488666\n03.04.2019\n..."
	reShahadatnama = regexp.MustCompile(`Baş potratçy:\s*\n\s*№\s*(\d+)\s*(?:\n\s*(\d{2}\.\d{2}\.\d{4}))?`)
	// "Baş potratçy:\n1-22-25-6722\n27.02.2019\n27.02.2022"
	reYgtyyarnama = regexp.MustCompile(`Baş potratçy:\s*\n\s*(\S+)\s*\n\s*(\d{2}\.\d{2}\.\d{4})\s*\n\s*(\d{2}\.\d{2}\.\d{4})`)
	// "Pasport seriyasy:\nI-MR 294244"
	rePasport = regexp.MustCompile(`Pasport seri[ýy]asy:\s*\n\s*(\S+)\s+(\d+)`)
	// "Patent seriyasy:\nJ №72363"
	rePatent = regexp.MustCompile(`Patent seri[ýy]asy:\s*\n\s*(\S+)\s*№\s*(\d+)`)
	// "Şahadatnamasy:\n№25482186"
	reShahadatDoc = regexp.MustCompile(`Şahadatnamasy:\s*\n\s*№\s*(\d+)`)
	// "03.2019\n03.2022"
	reBashySongy = regexp.MustCompile(`^\s*(\d{2})\.(\d{4})\s*\n\s*(\d{2})\.(\d{4})\s*$`)
	// "<whose> 23.10.2018y. 1610 buýrugy esasynda, <rest>"
	reKep = regexp.MustCompile(`^(.*?)\s+(\d{2}\.\d{2}\.\d{4})\s*ý?y?\.\s+(\S+)\s+(buýrugy|buýrygy|karary|teswirnama)\s+esasynda[,.]?\s*(.*)$`)
	// "Döwletnama №967 29.04.2019 21.10. 2019" (tolerates space inside 2nd date)
	reKepCert = regexp.MustCompile(`(\S+)\s*№\s*(\S+)\s+(\d{2}\.\d{2}\.\d{4})(?:\s+(\d{2}\.\d{2}\.?\s*\d{4}))?`)
	// "eli 65-53-66-54" inside salgy_paychy
	rePhone    = regexp.MustCompile(`eli\s*\.?\s*([\d\-\s]{7,})`)
	reNonDigit = regexp.MustCompile(`\D`)

	// emlak_paychy lines
	reEntrance = regexp.MustCompile(`^(\d+)\s*-\s*nj[iy]\s+girelge$`)
	reHouse    = regexp.MustCompile(`^(\d+)\s*-\s*nj[iy]\s+jaý$`)
	reFloor    = regexp.MustCompile(`^(\d+)\s*-\s*nj[iy]\s+gaty?$`)
	reApt      = regexp.MustCompile(`^(\d+)\s*-\s*öýi$`)
	reRooms    = regexp.MustCompile(`^(\d+)\s+otagly$`)
	reSquare   = regexp.MustCompile(`^([\d\s]+(?:[.,]\d+)?)\s*m2$`)

	orgTypes = []string{
		"Hususy kärhana", "Hojalyk jemgyýeti", "Döwlet kärhanasy", "Raýat",
		"Telekeçi", "Açyk görnüşli paýdarlar jemgyýeti", "Paýdarlar jemgyýeti",
		"Jemgyýetçilik guramasy", "Daşary ýurt kompaniýasy", "Edara",
	}
	headPositions = []string{
		"başlygy", "baş direktoryň orunbasary", "baş direktory", "baş direktor",
		"direktory", "direktor", "drektory", "ýolbaşçysy", "müdiri", "başlyk",
	}
	reQuoted = regexp.MustCompile(`^[«"]([^"»]+)["»]\s*(\S+)?\s*,?\s*(.*)$`)
)

func normSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func normKey(parts ...string) string {
	lowered := make([]string, len(parts))
	for i, p := range parts {
		lowered[i] = strings.ToLower(normSpace(p))
	}
	return strings.Join(lowered, "|")
}

func strPtr(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

// strPtrN trims and truncates to n characters (runes) to fit varchar columns
func strPtrN(s string, n int) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if r := []rune(s); len(r) > n {
		s = string(r[:n])
	}
	return &s
}

func orgFromParsed(org parsedOrg) models.Org {
	return models.Org{
		OrgType:           strPtrN(org.OrgType, 255),
		OrgName:           strPtrN(org.OrgName, 255),
		HeadPosition:      strPtrN(org.HeadPosition, 50),
		HeadFullName:      strPtrN(org.HeadFullName, 255),
		OrgAdditionalInfo: strPtr(org.Extra),
	}
}

func parseDMY(s string) *time.Time {
	m := reDMY.FindStringSubmatch(s)
	if m == nil {
		return nil
	}
	day, _ := strconv.Atoi(m[1])
	month, _ := strconv.Atoi(m[2])
	year, _ := strconv.Atoi(m[3])
	if month < 1 || month > 12 || day < 1 || day > 31 {
		return nil
	}
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return &t
}

// "29 223 621,0" -> 29223621.0 (regular and non-breaking spaces)
func parseMoney(s string) *float64 {
	cleaned := strings.NewReplacer(" ", "", " ", "", ",", ".").Replace(strings.TrimSpace(s))
	if cleaned == "" {
		return nil
	}
	f, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return nil
	}
	return &f
}

func atoiPtr(s string) *int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &n
}

type parsedOrg struct {
	OrgType      string
	OrgName      string
	HeadPosition string
	HeadFullName string
	Extra        string
}

// parseOrg handles: 'Raýat: Name', 'Hususy kärhana: "X", Başlygy Y',
// 'TS we T birleşmesi, başlygy A.Dadaýew', '"Birleşme Gurluşyk" JGK, direktor Y'
func parseOrg(s string) parsedOrg {
	s1 := normSpace(s)
	out := parsedOrg{}
	lower := strings.ToLower(s1)
	for _, t := range orgTypes {
		tl := strings.ToLower(t)
		if strings.HasPrefix(lower, tl+":") || strings.HasPrefix(lower, tl+" :") {
			rest := strings.TrimSpace(s1[strings.Index(s1, ":")+1:])
			out.OrgType = t
			if t == "Raýat" || t == "Telekeçi" {
				out.OrgName = rest
				return out
			}
			if m := reQuoted.FindStringSubmatch(rest); m != nil {
				out.OrgName = strings.Trim(m[1], ` ",:«»`)
				parseHeadTail(strings.TrimSpace(m[2]+" "+m[3]), &out)
			} else if idx := strings.Index(rest, ","); idx >= 0 {
				out.OrgName = strings.Trim(rest[:idx], ` ",:«»`)
				parseHeadTail(strings.TrimSpace(rest[idx+1:]), &out)
			} else {
				out.OrgName = strings.Trim(rest, ` ",:«»`)
			}
			return out
		}
	}
	// no known org-type prefix
	if m := reQuoted.FindStringSubmatch(s1); m != nil {
		out.OrgName = strings.Trim(m[1], ` ",:«»`)
		out.OrgType = strings.TrimSpace(m[2])
		parseHeadTail(strings.TrimSpace(m[3]), &out)
		return out
	}
	if idx := strings.Index(s1, ","); idx >= 0 {
		out.OrgName = strings.Trim(s1[:idx], ` ",:«»`)
		parseHeadTail(strings.TrimSpace(s1[idx+1:]), &out)
	} else {
		out.OrgName = strings.Trim(s1, ` ",:«»`)
	}
	return out
}

// parseHeadTail splits "başlygy Hudaýberdiýew Döwran" into position + name
func parseHeadTail(tail string, out *parsedOrg) {
	if tail == "" {
		return
	}
	tl := strings.ToLower(tail)
	for _, p := range headPositions {
		if idx := strings.Index(tl, p); idx >= 0 {
			out.HeadPosition = strings.TrimSpace(tail[:idx+len(p)])
			out.HeadFullName = strings.TrimSpace(tail[idx+len(p):])
			return
		}
	}
	out.Extra = tail
}

// ============================================================================
// Entity resolution with caching: find existing records or create new ones.
// ============================================================================

type oldMigrator struct {
	tx           *gorm.DB
	areasByName  map[string]*models.Area
	builders     map[string]*uint
	gcs          map[string]*uint
	buildings    map[string]*uint
	receivers    map[string]*uint
	shareholders map[string]*uint
	created      map[string]int
}

func newOldMigrator(tx *gorm.DB) (*oldMigrator, error) {
	m := &oldMigrator{
		tx:           tx,
		areasByName:  map[string]*models.Area{},
		builders:     map[string]*uint{},
		gcs:          map[string]*uint{},
		buildings:    map[string]*uint{},
		receivers:    map[string]*uint{},
		shareholders: map[string]*uint{},
		created:      map[string]int{},
	}
	var areas []models.Area
	if err := tx.Find(&areas).Error; err != nil {
		return nil, err
	}
	for i := range areas {
		m.areasByName[strings.ToLower(normSpace(areas[i].Name))] = &areas[i]
	}
	return m, nil
}

// splitAddress matches leading comma-separated segments against known areas;
// the unmatched remainder becomes the street/address text.
func (m *oldMigrator) splitAddress(s string) ([]models.Area, *string) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	segments := strings.Split(s, ",")
	var areas []models.Area
	i := 0
	for ; i < len(segments); i++ {
		key := strings.ToLower(normSpace(segments[i]))
		a, ok := m.areasByName[key]
		if !ok {
			break
		}
		areas = append(areas, *a)
	}
	rest := normSpace(strings.Join(segments[i:], ","))
	return areas, strPtrN(rest, 510)
}

func (m *oldMigrator) resolveBuilder(gurujy, salgy string) (*uint, error) {
	if strings.TrimSpace(gurujy) == "" {
		return nil, nil
	}
	key := normKey(gurujy, salgy)
	if id, ok := m.builders[key]; ok {
		return id, nil
	}
	org := parseOrg(gurujy)
	var existing models.Builder
	err := m.tx.Where("LOWER(org_name) = ?", strings.ToLower(org.OrgName)).First(&existing).Error
	if err == nil {
		m.builders[key] = &existing.ID
		return &existing.ID, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	areas, street := m.splitAddress(salgy)
	b := models.Builder{
		Org: orgFromParsed(org),
		BuilderAddress:  models.BuilderAddress{Areas: areas, Address: street},
		FromOldRegistry: true,
	}
	if err := m.tx.Create(&b).Error; err != nil {
		return nil, err
	}
	m.created["builders"]++
	m.builders[key] = &b.ID
	return &b.ID, nil
}

func (m *oldMigrator) resolveGeneralContractor(bashPotr, shahadatnama, ygtyyarnama string) (*uint, error) {
	name := normSpace(bashPotr)
	if name == "" || strings.EqualFold(name, "Ýok") || strings.EqualFold(name, "Yok") {
		return nil, nil
	}
	key := normKey(bashPotr, shahadatnama, ygtyyarnama)
	if id, ok := m.gcs[key]; ok {
		return id, nil
	}
	org := parseOrg(bashPotr)
	gc := models.GeneralContractor{
		FromOldRegistry: true,
		Contractor: models.Contractor{
			Org: orgFromParsed(org),
		},
	}
	if cm := reShahadatnama.FindStringSubmatch(shahadatnama); cm != nil {
		if n, err := strconv.Atoi(cm[1]); err == nil && n != 0 {
			gc.CertNumber = &n
			gc.CertDate = parseDMY(cm[2])
		}
	}
	if ym := reYgtyyarnama.FindStringSubmatch(ygtyyarnama); ym != nil {
		gc.ResolutionCode = strPtr(ym[1])
		gc.ResolutionBeginDate = parseDMY(ym[2])
		gc.ResolutionEndDate = parseDMY(ym[3])
	}
	var existing models.GeneralContractor
	q := m.tx.Where("LOWER(org_name) = ?", strings.ToLower(org.OrgName))
	if gc.CertNumber != nil {
		q = q.Where("cert_number = ?", *gc.CertNumber)
	}
	if err := q.First(&existing).Error; err == nil {
		m.gcs[key] = &existing.ID
		return &existing.ID, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err := m.tx.Create(&gc).Error; err != nil {
		return nil, err
	}
	m.created["general_contractors"]++
	m.gcs[key] = &gc.ID
	return &gc.ID, nil
}

func (m *oldMigrator) resolveBuilding(old models.OldRegistry) (*uint, error) {
	desga := normSpace(strOrEmpty(old.Desga))
	if desga == "" {
		return nil, nil
	}
	key := normKey(desga, strOrEmpty(old.SalgyDesga), strOrEmpty(old.BahaUmumy))
	if id, ok := m.buildings[key]; ok {
		return id, nil
	}
	areas, street := m.splitAddress(strOrEmpty(old.SalgyDesga))
	b := models.Building{
		FromOldRegistry: true,
		BuildingAddress: models.BuildingAddress{Areas: areas, Street: street},
		BuildingMain: models.BuildingMain{
			Kind:  strPtrN(desga, 510),
			Price: parseMoney(strOrEmpty(old.BahaUmumy)),
		},
	}
	if mm := reMeydan.FindStringSubmatch(strOrEmpty(old.MeydanUmumy)); mm != nil {
		b.Square1 = parseMoney(mm[1])
		b.Square1Name = strPtrN(mm[2], 5)
		b.SquareAdditionalInfo = strPtr(mm[3])
	} else if v := strPtr(strOrEmpty(old.MeydanUmumy)); v != nil {
		b.SquareAdditionalInfo = v
	}
	if sm := reBashySongy.FindStringSubmatch(strOrEmpty(old.SeneBashySongy)); sm != nil {
		startMonth, _ := strconv.Atoi(sm[1])
		startYear, _ := strconv.Atoi(sm[2])
		endMonth, _ := strconv.Atoi(sm[3])
		endYear, _ := strconv.Atoi(sm[4])
		if startMonth >= 1 && startMonth <= 12 {
			t := time.Date(startYear, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
			b.StartDate = &t
		}
		if endMonth >= 1 && endMonth <= 12 {
			t := time.Date(endYear, time.Month(endMonth), 1, 0, 0, 0, 0, time.UTC)
			b.EndDate = &t
		}
	}
	kep := normSpace(strOrEmpty(old.KepResminama))
	if km := reKep.FindStringSubmatch(kep); km != nil {
		b.OrderWhoseWhat = strPtrN(km[1]+" "+km[4], 255)
		b.OrderDate = parseDMY(km[2])
		b.OrderCode = strPtrN(km[3], 50)
		rest := strings.TrimSpace(km[5])
		if cm := reKepCert.FindStringSubmatch(rest); cm != nil {
			b.CertName = strPtrN(cm[1], 50)
			b.Cert1Code = strPtrN(cm[2], 50)
			b.Cert1Date = parseDMY(cm[3])
			if cm[4] != "" {
				b.Cert2Date = parseDMY(strings.ReplaceAll(cm[4], " ", ""))
			}
		} else {
			b.OrderAdditionalInfo = strPtr(rest)
		}
	} else if kep != "" {
		b.OrderAdditionalInfo = strPtr(kep)
	}
	if err := m.tx.Create(&b).Error; err != nil {
		return nil, err
	}
	m.created["buildings"]++
	m.buildings[key] = &b.ID
	return &b.ID, nil
}

func (m *oldMigrator) resolveShareholder(old models.OldRegistry) (*uint, error) {
	paychy := normSpace(strOrEmpty(old.Paychy))
	if paychy == "" {
		return nil, nil
	}
	docs := strOrEmpty(old.PatentPasport)
	var passportSeries, passportNumber string
	if pm := rePasport.FindStringSubmatch(docs); pm != nil {
		passportSeries, passportNumber = pm[1], pm[2]
	}
	key := normKey(paychy, strOrEmpty(old.SalgyPaychy))
	if passportNumber != "" {
		key = normKey("pass", passportSeries, passportNumber)
	}
	if id, ok := m.shareholders[key]; ok {
		return id, nil
	}
	if passportNumber != "" {
		var existing models.Shareholder
		err := m.tx.Where("passport_series = ? AND passport_number = ?", passportSeries, passportNumber).First(&existing).Error
		if err == nil {
			m.shareholders[key] = &existing.ID
			return &existing.ID, nil
		}
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	org := parseOrg(paychy)
	areas, address := m.splitAddress(strOrEmpty(old.SalgyPaychy))
	sh := models.Shareholder{
		Org: orgFromParsed(org),
		ShareholderAddress: models.ShareholderAddress{Areas: areas, Address: address},
		FromOldRegistry:    true,
	}
	if passportSeries != "" && len(passportSeries) <= 6 && len(passportNumber) <= 6 {
		sh.PassportSeries = strPtr(passportSeries)
		sh.PassportNumber = strPtr(passportNumber)
	}
	if pm := rePatent.FindStringSubmatch(docs); pm != nil {
		if len(pm[1]) <= 2 {
			sh.PatentSeries = strPtr(pm[1])
		}
		if n, err := strconv.ParseUint(pm[2], 10, 32); err == nil {
			u := uint(n)
			sh.CertNumber = nil
			sh.PatentNumber = &u
		}
	}
	if cm := reShahadatDoc.FindStringSubmatch(docs); cm != nil {
		if n, err := strconv.ParseUint(cm[1], 10, 32); err == nil {
			u := uint(n)
			sh.CertNumber = &u
		}
	}
	if err := m.tx.Create(&sh).Error; err != nil {
		return nil, err
	}
	m.created["shareholders"]++
	// phones like "eli 65-53-66-54" hidden in the address text
	for _, pm := range rePhone.FindAllStringSubmatch(strOrEmpty(old.SalgyPaychy), -1) {
		number := reNonDigit.ReplaceAllString(pm[1], "")
		if len(number) < 5 || len(number) > 12 {
			continue
		}
		kind := "el"
		phone := models.Phone{Kind: &kind, Number: &number, ShareholderID: sh.ID}
		if err := m.tx.Create(&phone).Error; err != nil {
			return nil, err
		}
		m.created["phones"]++
	}
	m.shareholders[key] = &sh.ID
	return &sh.ID, nil
}

func (m *oldMigrator) resolveReceiver(wezipe, ady string) (*uint, error) {
	if strings.TrimSpace(ady) == "" && strings.TrimSpace(wezipe) == "" {
		return nil, nil
	}
	key := normKey(wezipe, ady)
	if id, ok := m.receivers[key]; ok {
		return id, nil
	}
	org := parseOrg(wezipe)
	r := models.Receiver{
		FromOldRegistry: true,
		CitizenStatus:   org.OrgType,
		OrgName:       org.OrgName,
		Position:      normSpace(strings.TrimSpace(org.HeadPosition + " " + org.Extra)),
	}
	nameParts := strings.Fields(normSpace(ady))
	switch {
	case len(nameParts) >= 3:
		r.Lastname = nameParts[0]
		r.Firstname = nameParts[1]
		r.Patronymic = strings.Join(nameParts[2:], " ")
	case len(nameParts) == 2:
		r.Lastname = nameParts[0]
		r.Firstname = nameParts[1]
	case len(nameParts) == 1:
		r.Lastname = nameParts[0]
	}
	var existing models.Receiver
	err := m.tx.Where("org_name = ? AND lastname = ? AND firstname = ?", r.OrgName, r.Lastname, r.Firstname).First(&existing).Error
	if err == nil {
		m.receivers[key] = &existing.ID
		return &existing.ID, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err := m.tx.Create(&r).Error; err != nil {
		return nil, err
	}
	m.created["receivers"]++
	m.receivers[key] = &r.ID
	return &r.ID, nil
}

// createShareholderProperty parses emlak_paychy into a structured property row
func (m *oldMigrator) createShareholderProperty(old models.OldRegistry, registryID uint, tb int) error {
	emlak := strOrEmpty(old.EmlakPaychy)
	if strings.TrimSpace(emlak) == "" {
		return nil
	}
	lines := strings.Split(emlak, "\n")
	sp := models.ShareholderProperty{
		TB:         &tb,
		RegistryID: registryID,
		Price:      parseMoney(strOrEmpty(old.BahaPaychy)),
		Price1m2:   parseMoney(strOrEmpty(old.Baha1m2Paychy)),
	}
	first := normSpace(strings.ReplaceAll(lines[0], ", ,", ""))
	sp.BuildingType = strPtr(strings.Trim(first, " ,"))
	var extras []string
	for _, l := range lines[1:] {
		l2 := strings.Trim(normSpace(l), " ,.")
		switch {
		case l2 == "":
		case reEntrance.MatchString(l2):
			sp.Entrance = atoiPtr(reEntrance.FindStringSubmatch(l2)[1])
		case reHouse.MatchString(l2):
			sp.Building = atoiPtr(reHouse.FindStringSubmatch(l2)[1])
		case reFloor.MatchString(l2):
			sp.Floor = atoiPtr(reFloor.FindStringSubmatch(l2)[1])
		case reApt.MatchString(l2):
			sp.Apartment = atoiPtr(reApt.FindStringSubmatch(l2)[1])
		case reRooms.MatchString(l2):
			sp.RoomCount = atoiPtr(reRooms.FindStringSubmatch(l2)[1])
		case reSquare.MatchString(l2):
			sp.Square = parseMoney(reSquare.FindStringSubmatch(l2)[1])
		default:
			extras = append(extras, l2)
		}
	}
	if len(extras) > 0 {
		sp.AdditionalInfo = strPtr(strings.Join(extras, "; "))
	}
	if err := m.tx.Create(&sp).Error; err != nil {
		return err
	}
	m.created["shareholder_properties"]++
	return nil
}

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// buildRegistry maps one old registry row into a fully adapted Registry
func (m *oldMigrator) buildRegistry(old models.OldRegistry, userID uint) (*models.Registry, error) {
	oldID := old.ID
	oldTB := old.TB
	tb := int(old.TB)

	reg := models.Registry{
		TB:     &tb,
		UserID: &userID,
		RegistryDates: models.RegistryDates{
			ReviewedAt: old.SeneSeredilen,
		},
		RegistryMail: models.RegistryMail{
			MinToMudDate: old.SeneHatMinToMud,
		},
		RegistryOldData: models.RegistryOldData{
			OldRegistryID:           &oldID,
			OldTB:                   &oldTB,
			OldMinHat:               old.MinHat,
			OldSeneHatMinToMud:      old.SeneHatMinToMud,
			OldGurujy:               old.Gurujy,
			OldPaychy:               old.Paychy,
			OldSertnamaGurujyPaychy: old.SertnamaGurujyPaychy,
			OldDesga:                old.Desga,
			OldBahaUmumy:            old.BahaUmumy,
			OldMeydanUmumy:          old.MeydanUmumy,
			OldKepResminama:         old.KepResminama,
			OldEmlakPaychy:          old.EmlakPaychy,
			OldBahaPaychy:           old.BahaPaychy,
			OldBaha1m2Paychy:        old.Baha1m2Paychy,
			OldSalgyDesga:           old.SalgyDesga,
			OldSalgyGurujy:          old.SalgyGurujy,
			OldSalgyPaychy:          old.SalgyPaychy,
			OldBashPotr:             old.BashPotr,
			OldSertnamaGurPotr:      old.SertnamaGurPotr,
			OldPotratchyKomek:       old.PotratchyKomek,
			OldShahadatnama:         old.Shahadatnama,
			OldYgtyyarnama:          old.Ygtyyarnama,
			OldPatentPasport:        old.PatentPasport,
			OldSeneBashySongy:       old.SeneBashySongy,
			OldSeneSeredilen:        old.SeneSeredilen,
			OldSeneHasabaAlnan:      old.SeneHasabaAlnan,
			OldWezipeAlanAdam:       old.WezipeAlanAdam,
			OldAdyAlanAdam:          old.AdyAlanAdam,
			OldSeneSanSertnama:      old.SeneSanSertnama,
			OldAdyPaychyAlan:        old.AdyPaychyAlan,
			OldSenePaychyAlan:       old.SenePaychyAlan,
			OldLogin:                old.Login,
		},
	}

	// mail: "date \n №num \n date2 \n (N sany \n M-tapgyr)"
	if mm := reMinHat.FindStringSubmatch(strOrEmpty(old.MinHat)); mm != nil {
		reg.MailDate = parseDMY(mm[1])
		reg.MailNumber = strPtr(mm[2])
		reg.DeliveryDate = parseDMY(mm[3])
		reg.Count = atoiPtr(mm[4])
		reg.Queue = atoiPtr(mm[5])
	}

	// registered_at from "PGGŞ \n 1-1/11026 \n 28.10.2019"
	if hm := strOrEmpty(old.SeneHasabaAlnan); hm != "" {
		reg.RegisteredAt = parseDMY(hm)
	}

	// contracts
	if cm := reContract.FindStringSubmatch(strOrEmpty(old.SertnamaGurujyPaychy)); cm != nil {
		reg.BuilderShareholderNumber = strPtr(cm[1])
		reg.BuilderShareholderDate = parseDMY(cm[2])
		reg.BuilderShareholderAdditionalInfo = strPtr(normSpace(cm[3]))
	}
	if cm := reContract.FindStringSubmatch(strOrEmpty(old.SertnamaGurPotr)); cm != nil {
		reg.BuilderContractorNumber = strPtr(cm[1])
		reg.BuilderContractorDate = parseDMY(cm[2])
		reg.BuilderContractorAdditionalInfo = strPtr(normSpace(cm[3]))
	}

	// denial (OTKAZ rows)
	if reason := strPtr(strOrEmpty(old.AdyPaychyAlan)); reason != nil {
		reg.DenialReason = reason
		if info := strOrEmpty(old.SenePaychyAlan); info != "" {
			reg.DenialDate = parseDMY(info)
			reg.DenialAdditionalInfo = strPtr(normSpace(info))
		}
	}

	// linked entities
	var err error
	if reg.BuilderID, err = m.resolveBuilder(strOrEmpty(old.Gurujy), strOrEmpty(old.SalgyGurujy)); err != nil {
		return nil, err
	}
	if reg.GeneralContractorID, err = m.resolveGeneralContractor(
		strOrEmpty(old.BashPotr), strOrEmpty(old.Shahadatnama), strOrEmpty(old.Ygtyyarnama)); err != nil {
		return nil, err
	}
	if reg.BuildingID, err = m.resolveBuilding(old); err != nil {
		return nil, err
	}
	if reg.ShareholderID, err = m.resolveShareholder(old); err != nil {
		return nil, err
	}
	if reg.ReceiverID, err = m.resolveReceiver(strOrEmpty(old.WezipeAlanAdam), strOrEmpty(old.AdyAlanAdam)); err != nil {
		return nil, err
	}

	return &reg, nil
}

// MigrateOldRegistries handles POST request to migrate every not-yet-migrated
// old registry into the registries table, adapting the legacy free-text data
// into the normalized schema (builders, shareholders, buildings, general
// contractors, receivers, shareholder properties) while keeping verbatim
// copies in the old_* columns. Idempotent via old_registry_id.
func MigrateOldRegistries(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	typedUser := user.(models.User)

	var total int64
	if err := initializers.DB.Model(&models.OldRegistry{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count old registries"})
		return
	}

	migratedSubQuery := initializers.DB.Model(&models.Registry{}).
		Select("old_registry_id").
		Where("old_registry_id IS NOT NULL")

	var oldRegistries []models.OldRegistry
	if err := initializers.DB.
		Where("id NOT IN (?)", migratedSubQuery).
		Order("t_b").
		Find(&oldRegistries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch old registries"})
		return
	}

	if len(oldRegistries) == 0 {
		c.JSON(200, gin.H{
			"message":  "Nothing to migrate, all old registries are already migrated",
			"total":    total,
			"migrated": 0,
			"skipped":  total,
		})
		return
	}

	migrated := 0
	var createdCounts map[string]int
	err := initializers.DB.Transaction(func(tx *gorm.DB) error {
		m, err := newOldMigrator(tx)
		if err != nil {
			return err
		}
		for _, old := range oldRegistries {
			reg, err := m.buildRegistry(old, typedUser.ID)
			if err != nil {
				return err
			}
			if err := tx.Create(reg).Error; err != nil {
				return err
			}
			if err := m.createShareholderProperty(old, reg.ID, int(old.TB)); err != nil {
				return err
			}
			migrated++
		}
		createdCounts = m.created
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to migrate old registries", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":  "Old registries migrated successfully",
		"total":    total,
		"migrated": migrated,
		"skipped":  total - int64(migrated),
		"created":  createdCounts,
	})
}

// RollbackOldRegistriesMigration handles POST request to remove everything the
// old-registries migration created: migrated registries (old_registry_id set),
// their shareholder properties and additional agreements, and every entity
// flagged from_old_registry that is no longer referenced by any remaining
// registry. Manually created records and matched pre-existing entities are
// never touched. The old_registries table itself is left intact, so the
// migration can be re-run afterwards.
func RollbackOldRegistriesMigration(c *gin.Context) {
	if _, exists := c.Get("user"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	deleted := map[string]int64{}
	err := initializers.DB.Transaction(func(tx *gorm.DB) error {
		migratedRegistries := tx.Model(&models.Registry{}).Unscoped().
			Select("id").Where("old_registry_id IS NOT NULL")

		res := tx.Unscoped().Where("registry_id IN (?)", migratedRegistries).Delete(&models.ShareholderProperty{})
		if res.Error != nil {
			return res.Error
		}
		deleted["shareholder_properties"] = res.RowsAffected

		res = tx.Unscoped().Where("registry_id IN (?)", migratedRegistries).Delete(&models.AdditionalAgreement{})
		if res.Error != nil {
			return res.Error
		}
		deleted["additional_agreements"] = res.RowsAffected

		res = tx.Unscoped().Where("old_registry_id IS NOT NULL").Delete(&models.Registry{})
		if res.Error != nil {
			return res.Error
		}
		deleted["registries"] = res.RowsAffected

		// phones of migration-created shareholders
		res = tx.Unscoped().
			Where("shareholder_id IN (?)", tx.Model(&models.Shareholder{}).Unscoped().Select("id").Where("from_old_registry = true")).
			Delete(&models.Phone{})
		if res.Error != nil {
			return res.Error
		}
		deleted["phones"] = res.RowsAffected

		// area join rows of migration-created entities, then the entities
		// themselves, guarded so records still referenced by remaining
		// registries survive
		type entity struct {
			name      string
			joinTable string
			joinCol   string
			fkCol     string
			model     interface{}
		}
		entities := []entity{
			{"shareholders", "shareholder_areas", "shareholder_id", "shareholder_id", &models.Shareholder{}},
			{"builders", "builder_areas", "builder_id", "builder_id", &models.Builder{}},
			{"buildings", "building_areas", "building_id", "building_id", &models.Building{}},
			{"general_contractors", "", "", "general_contractor_id", &models.GeneralContractor{}},
			{"receivers", "", "", "receiver_id", &models.Receiver{}},
		}
		for _, e := range entities {
			stillReferenced := tx.Model(&models.Registry{}).Unscoped().
				Select(e.fkCol).Where(e.fkCol + " IS NOT NULL")
			if e.joinTable != "" {
				if err := tx.Exec(
					"DELETE FROM "+e.joinTable+" WHERE "+e.joinCol+" IN (SELECT id FROM "+e.name+" WHERE from_old_registry = true AND id NOT IN (?))",
					stillReferenced,
				).Error; err != nil {
					return err
				}
			}
			res = tx.Unscoped().
				Where("from_old_registry = true").
				Where("id NOT IN (?)", stillReferenced).
				Delete(e.model)
			if res.Error != nil {
				return res.Error
			}
			deleted[e.name] = res.RowsAffected
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rollback migration", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Old registries migration rolled back successfully",
		"deleted": deleted,
	})
}
