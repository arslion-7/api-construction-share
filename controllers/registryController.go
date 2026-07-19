package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/arslion-7/api-construction-share/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// generateShareholderDescription creates a description from shareholder org info and docs additional info
func generateShareholderDescription(shareholder *models.Shareholder) string {
	if shareholder == nil {
		return ""
	}

	var parts []string

	// Add organization info with translated titles
	if shareholder.OrgName != nil && *shareholder.OrgName != "" {
		orgNameLabel := "Guramanyň ady"
		if shareholder.OrgType != nil && *shareholder.OrgType == "Raýat" {
			orgNameLabel = "Raýat A.F.Aa"
		}
		parts = append(parts, fmt.Sprintf("%s: %s", orgNameLabel, *shareholder.OrgName))
	}

	if shareholder.OrgType != nil && *shareholder.OrgType != "" {
		parts = append(parts, fmt.Sprintf("Gurama görnüşi: %s", *shareholder.OrgType))
	}

	if shareholder.HeadFullName != nil && *shareholder.HeadFullName != "" {
		parts = append(parts, fmt.Sprintf("Ýolbaşçynyň A.F.Aa: %s", *shareholder.HeadFullName))
	}

	if shareholder.HeadPosition != nil && *shareholder.HeadPosition != "" {
		parts = append(parts, fmt.Sprintf("Ýolbaşçynyň wezipesi: %s", *shareholder.HeadPosition))
	}

	if shareholder.OrgAdditionalInfo != nil && *shareholder.OrgAdditionalInfo != "" {
		parts = append(parts, fmt.Sprintf("Goşmaça maglumat: %s", *shareholder.OrgAdditionalInfo))
	}

	// Add documents additional info
	if shareholder.DocsAdditionalInfo != nil && *shareholder.DocsAdditionalInfo != "" {
		parts = append(parts, fmt.Sprintf("Goşmaça maglumaty: %s", *shareholder.DocsAdditionalInfo))
	}

	return strings.Join(parts, "\n")
}

// generateGeneralContractorDescription creates a description from general contractor org info, cert, and resolution
func generateGeneralContractorDescription(contractor *models.GeneralContractor) string {
	if contractor == nil {
		return ""
	}

	var parts []string

	// Add organization info with translated titles
	if contractor.OrgName != nil && *contractor.OrgName != "" {
		orgNameLabel := "Guramanyň ady"
		if contractor.OrgType != nil && *contractor.OrgType == "Raýat" {
			orgNameLabel = "Raýat A.F.Aa"
		}
		parts = append(parts, fmt.Sprintf("%s: %s", orgNameLabel, *contractor.OrgName))
	}

	if contractor.OrgType != nil && *contractor.OrgType != "" {
		parts = append(parts, fmt.Sprintf("Gurama görnüşi: %s", *contractor.OrgType))
	}

	if contractor.HeadFullName != nil && *contractor.HeadFullName != "" {
		parts = append(parts, fmt.Sprintf("Ýolbaşçynyň A.F.Aa: %s", *contractor.HeadFullName))
	}

	if contractor.HeadPosition != nil && *contractor.HeadPosition != "" {
		parts = append(parts, fmt.Sprintf("Ýolbaşçynyň wezipesi: %s", *contractor.HeadPosition))
	}

	// Add certificate info
	if contractor.CertNumber != nil {
		parts = append(parts, fmt.Sprintf("Sertifikat belgisi: %d", *contractor.CertNumber))
	}

	if contractor.CertDate != nil {
		parts = append(parts, fmt.Sprintf("Sertifikat senesi: %s", contractor.CertDate.Format("02.01.2006")))
	}

	// Add resolution info
	if contractor.ResolutionCode != nil && *contractor.ResolutionCode != "" {
		parts = append(parts, fmt.Sprintf("Ygtyýarnama belgisi: %s", *contractor.ResolutionCode))
	}

	if contractor.ResolutionBeginDate != nil {
		parts = append(parts, fmt.Sprintf("Ygtyýarnama başy: %s", contractor.ResolutionBeginDate.Format("02.01.2006")))
	}

	if contractor.ResolutionEndDate != nil {
		parts = append(parts, fmt.Sprintf("Ygtyýarnama soňy: %s", contractor.ResolutionEndDate.Format("02.01.2006")))
	}

	if contractor.OrgAdditionalInfo != nil && *contractor.OrgAdditionalInfo != "" {
		parts = append(parts, fmt.Sprintf("Goşmaça maglumat: %s", *contractor.OrgAdditionalInfo))
	}

	return strings.Join(parts, "\n")
}

// generateBuildingDescription creates a description from building info
func generateBuildingDescription(building *models.Building) string {
	if building == nil {
		return ""
	}

	var parts []string

	if building.TB != nil {
		parts = append(parts, fmt.Sprintf("Tertip belgisi: %d", *building.TB))
	}

	if building.IdentNumber != nil {
		parts = append(parts, fmt.Sprintf("Desga belgisi: %d", *building.IdentNumber))
	}

	if building.Kind != nil && *building.Kind != "" {
		parts = append(parts, fmt.Sprintf("Desga görnüşi: %s", *building.Kind))
	}

	// Add address info
	if building.Areas != nil && len(building.Areas) > 0 {
		areaNames := make([]string, len(building.Areas))
		for i, area := range building.Areas {
			areaNames[i] = area.Name
		}
		parts = append(parts, fmt.Sprintf("Adres: %s", strings.Join(areaNames, ", ")))
	}

	if building.Street != nil && *building.Street != "" {
		parts = append(parts, fmt.Sprintf("Köçe: %s", *building.Street))
	}

	// Add price and percentage
	if building.Price != nil {
		parts = append(parts, fmt.Sprintf("Bahasy: %d man.", *building.Price))
	}

	if building.Percentage != nil {
		parts = append(parts, fmt.Sprintf("Göterimi: %d%%", *building.Percentage))
	}

	// Add order info
	if building.OrderWhoseWhat != nil && *building.OrderWhoseWhat != "" {
		parts = append(parts, fmt.Sprintf("Karar näme: %s", *building.OrderWhoseWhat))
	}

	if building.OrderCode != nil && *building.OrderCode != "" {
		parts = append(parts, fmt.Sprintf("Karar belgisi: %s", *building.OrderCode))
	}

	if building.OrderDate != nil {
		parts = append(parts, fmt.Sprintf("Karar senesi: %s", building.OrderDate.Format("02.01.2006")))
	}

	// Add square info
	if building.Square1 != nil && building.Square1Name != nil {
		parts = append(parts, fmt.Sprintf("%s: %dm²", *building.Square1Name, *building.Square1))
	}

	if building.Square2 != nil && building.Square2Name != nil {
		parts = append(parts, fmt.Sprintf("%s: %dm²", *building.Square2Name, *building.Square2))
	}

	if building.Square3 != nil && building.Square3Name != nil {
		parts = append(parts, fmt.Sprintf("%s: %dm²", *building.Square3Name, *building.Square3))
	}

	// Add certificate info
	if building.CertName != nil && *building.CertName != "" {
		parts = append(parts, fmt.Sprintf("Sertifikat ady: %s", *building.CertName))
	}

	if building.Cert1Code != nil && *building.Cert1Code != "" {
		parts = append(parts, fmt.Sprintf("1-nji sertifikat: %s", *building.Cert1Code))
	}

	if building.Cert2Code != nil && *building.Cert2Code != "" {
		parts = append(parts, fmt.Sprintf("2-nji sertifikat: %s", *building.Cert2Code))
	}

	if building.OrderAdditionalInfo != nil && *building.OrderAdditionalInfo != "" {
		parts = append(parts, fmt.Sprintf("Karar goşmaça: %s", *building.OrderAdditionalInfo))
	}

	if building.SquareAdditionalInfo != nil && *building.SquareAdditionalInfo != "" {
		parts = append(parts, fmt.Sprintf("Meýdan goşmaça: %s", *building.SquareAdditionalInfo))
	}

	return strings.Join(parts, "\n")
}

// generateBuilderDescription creates a description from builder org info
func generateBuilderDescription(builder *models.Builder) string {
	if builder == nil {
		return ""
	}

	var parts []string

	// Add organization info with translated titles
	if builder.OrgName != nil && *builder.OrgName != "" {
		orgNameLabel := "Guramanyň ady"
		if builder.OrgType != nil && *builder.OrgType == "Raýat" {
			orgNameLabel = "Raýat A.F.Aa"
		}
		parts = append(parts, fmt.Sprintf("%s: %s", orgNameLabel, *builder.OrgName))
	}

	if builder.OrgType != nil && *builder.OrgType != "" {
		parts = append(parts, fmt.Sprintf("Gurama görnüşi: %s", *builder.OrgType))
	}

	if builder.HeadFullName != nil && *builder.HeadFullName != "" {
		parts = append(parts, fmt.Sprintf("Ýolbaşçynyň A.F.Aa: %s", *builder.HeadFullName))
	}

	if builder.HeadPosition != nil && *builder.HeadPosition != "" {
		parts = append(parts, fmt.Sprintf("Ýolbaşçynyň wezipesi: %s", *builder.HeadPosition))
	}

	// Add address info
	if builder.Areas != nil && len(builder.Areas) > 0 {
		areaNames := make([]string, len(builder.Areas))
		for i, area := range builder.Areas {
			areaNames[i] = area.Name
		}
		parts = append(parts, fmt.Sprintf("Adres: %s", strings.Join(areaNames, ", ")))
	}

	if builder.Address != nil && *builder.Address != "" {
		parts = append(parts, fmt.Sprintf("Köçe: %s", *builder.Address))
	}

	if builder.OrgAdditionalInfo != nil && *builder.OrgAdditionalInfo != "" {
		parts = append(parts, fmt.Sprintf("Goşmaça maglumat: %s", *builder.OrgAdditionalInfo))
	}

	return strings.Join(parts, "\n")
}

// applyRegistryListFilters applies source/search/column filters shared by the
// data and count queries of GetRegistries
func applyRegistryListFilters(db *gorm.DB, c *gin.Context) *gorm.DB {
	shareholdersJoined := false
	joinShareholders := func() {
		if !shareholdersJoined {
			db = db.Joins("LEFT JOIN shareholders ON shareholders.id = registries.shareholder_id")
			shareholdersJoined = true
		}
	}

	// source=old -> only registries migrated from old_registries,
	// source=new -> only manually created ones
	switch c.Query("source") {
	case "old":
		db = db.Where("registries.old_registry_id IS NOT NULL")
	case "new":
		db = db.Where("registries.old_registry_id IS NULL")
	}

	if search := c.Query("search"); search != "" {
		if tb, err := strconv.Atoi(search); err == nil {
			db = db.Where("registries.t_b = ?", tb)
		} else {
			joinShareholders()
			like := "%" + strings.ToLower(search) + "%"
			db = db.Where("LOWER(shareholders.org_name) LIKE ? OR LOWER(shareholders.head_full_name) LIKE ? OR LOWER(shareholders.docs_additional_info) LIKE ?", like, like, like)
		}
	}

	// column filters
	if tbStr := strings.TrimSpace(c.Query("t_b")); tbStr != "" {
		if tb, err := strconv.Atoi(tbStr); err == nil {
			db = db.Where("registries.t_b = ?", tb)
		}
	}
	if ids := parseIDList(c.Query("gc_ids")); len(ids) > 0 {
		db = db.Where("registries.general_contractor_id IN ?", ids)
	}
	if ids := parseIDList(c.Query("user_ids")); len(ids) > 0 {
		db = db.Where("registries.user_id IN ?", ids)
	}
	if q := strings.TrimSpace(c.Query("shareholder_q")); q != "" {
		joinShareholders()
		like := "%" + strings.ToLower(q) + "%"
		db = db.Where("LOWER(shareholders.org_name) LIKE ? OR LOWER(shareholders.head_full_name) LIKE ? OR LOWER(shareholders.docs_additional_info) LIKE ?", like, like, like)
	}
	if q := strings.TrimSpace(c.Query("builder_q")); q != "" {
		like := "%" + strings.ToLower(q) + "%"
		db = db.Joins("LEFT JOIN builders ON builders.id = registries.builder_id").
			Where("LOWER(builders.org_name) LIKE ? OR LOWER(builders.head_full_name) LIKE ?", like, like)
	}
	if q := strings.TrimSpace(c.Query("building_q")); q != "" {
		like := "%" + strings.ToLower(q) + "%"
		db = db.Joins("LEFT JOIN buildings ON buildings.id = registries.building_id").
			Where("LOWER(buildings.kind) LIKE ? OR LOWER(buildings.street) LIKE ?", like, like)
	}
	return db
}

func parseIDList(s string) []int {
	var ids []int
	for _, p := range strings.Split(s, ",") {
		if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			ids = append(ids, n)
		}
	}
	return ids
}

func GetRegistries(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)

	var data []models.Registry
	query := applyRegistryListFilters(
		initializers.DB.Model(&models.Registry{}).
			Preload("User").
			Preload("GeneralContractor").
			Preload("Building").
			Preload("Builder").
			Preload("Receiver").
			Preload("Shareholder"), c).
		Order("registries.t_b").
		Limit(pagination.PageSize).
		Offset(pagination.Offset)

	if err := query.Find(&data).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve data", "details": err.Error()})
		return
	}

	// Generate shareholder description for each registry
	for i := range data {
		data[i].ShareholderDescription = generateShareholderDescription(data[i].Shareholder)
		data[i].GeneralContractorDescription = generateGeneralContractorDescription(data[i].GeneralContractor)
		data[i].BuildingDescription = generateBuildingDescription(data[i].Building)
		data[i].BuilderDescription = generateBuilderDescription(data[i].Builder)
	}

	var total int64
	applyRegistryListFilters(initializers.DB.Model(&models.Registry{}), c).Count(&total)

	utils.RespondWithPagination(c, data, pagination, total)
}

type RegistryFilterOption struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// GetRegistryFilterOptions returns the distinct general contractors and users
// referenced by registries, for the table column filter dropdowns
func GetRegistryFilterOptions(c *gin.Context) {
	var gcs []RegistryFilterOption
	initializers.DB.Raw(`SELECT DISTINCT g.id, COALESCE(g.org_name, '') AS name
		FROM general_contractors g
		JOIN registries r ON r.general_contractor_id = g.id
		WHERE r.deleted_at IS NULL
		ORDER BY name`).Scan(&gcs)

	var users []RegistryFilterOption
	initializers.DB.Raw(`SELECT DISTINCT u.id, u.email AS name
		FROM users u
		JOIN registries r ON r.user_id = u.id
		WHERE r.deleted_at IS NULL
		ORDER BY name`).Scan(&users)

	c.JSON(200, gin.H{"general_contractors": gcs, "users": users})
}

func GetRegistry(c *gin.Context) {
	id := c.Param("id")

	var registry models.Registry

	if err := initializers.DB.
		Preload("User").
		Preload("GeneralContractor").
		Preload("Building").
		Preload("Builder").
		Preload("Receiver").
		Preload("Shareholder").
		// Preload("ShareholderProperty").
		First(&registry, id).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	// Generate shareholder description
	registry.ShareholderDescription = generateShareholderDescription(registry.Shareholder)
	registry.GeneralContractorDescription = generateGeneralContractorDescription(registry.GeneralContractor)
	registry.BuildingDescription = generateBuildingDescription(registry.Building)
	registry.BuilderDescription = generateBuilderDescription(registry.Builder)

	c.JSON(200, registry)
}

type MainRegistryInput struct {
	TB           *int       `gorm:"column:t_b" json:"t_b"`
	ReviewedAt   *time.Time `json:"reviewed_at"`
	RegisteredAt *time.Time `json:"registered_at"`
}

func CreateRegistry(c *gin.Context) {
	var input MainRegistryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var registry models.Registry

	user, exists := c.Get("user")
	fmt.Println("user", user)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	typedUser := user.(models.User)

	registry.TB = input.TB
	registry.ReviewedAt = input.ReviewedAt
	registry.RegisteredAt = input.RegisteredAt
	registry.UserID = &typedUser.ID

	if err := initializers.DB.Save(&registry).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, registry)
}

func UpdateRegistryMail(c *gin.Context) {
	var input models.RegistryMail

	var registry models.Registry
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&registry).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.BindJSON(&input)

	registry.MailDate = input.MailDate
	registry.MailNumber = input.MailNumber
	registry.DeliveryDate = input.DeliveryDate
	registry.Count = input.Count
	registry.Queue = input.Queue
	registry.MinToMudDate = input.MinToMudDate

	initializers.DB.Save(&registry)
	c.JSON(200, registry)
}

func UpdateRegistry(c *gin.Context) {
	var input MainRegistryInput

	var registry models.Registry
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&registry).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.BindJSON(&input)

	registry.TB = input.TB
	registry.ReviewedAt = input.ReviewedAt
	registry.RegisteredAt = input.RegisteredAt

	initializers.DB.Save(&registry)
	c.JSON(200, registry)
}

// type UpdateRegistryGeneralContractorInput struct {
// 	GeneralContractorId `json:general_contractor_id`
// }

type UpdateRegistryNumberInput struct {
	TB int `json:"t_b"`
}

func UpdateRegistryNumber(c *gin.Context) {
	var input UpdateRegistryNumberInput

	var registry models.Registry
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&registry).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.BindJSON(&input)

	registry.TB = &input.TB

	initializers.DB.Save(&registry)
	c.JSON(200, registry)
}

func UpdateRegistryDates(c *gin.Context) {
	var input models.RegistryDates

	var registry models.Registry
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&registry).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.BindJSON(&input)

	registry.ReviewedAt = input.ReviewedAt
	registry.RegisteredAt = input.RegisteredAt

	initializers.DB.Save(&registry)
	c.JSON(200, registry)
}

func UpdateRegistryGeneralContractor(c *gin.Context) {
	var input models.Registry

	var registry models.Registry
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&registry).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.BindJSON(&input)

	initializers.DB.Model(&models.Registry{}).Where("id = ?", id).Update("general_contractor_id", input.GeneralContractorID)

	c.JSON(200, registry)

}

func UpdateRegistryBuilding(c *gin.Context) {
	var input models.Registry

	var registry models.Registry
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&registry).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.BindJSON(&input)

	initializers.DB.Model(&models.Registry{}).Where("id = ?", id).Update("building_id", input.BuildingID)

	c.JSON(200, registry)

}

func UpdateRegistryBuilder(c *gin.Context) {
	var input models.Registry

	var registry models.Registry
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&registry).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.BindJSON(&input)

	initializers.DB.Model(&models.Registry{}).Where("id = ?", id).Update("builder_id", input.BuilderID)

	c.JSON(200, registry)

}

func UpdateRegistryReceiver(c *gin.Context) {
	var input models.Registry

	var registry models.Registry
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&registry).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.BindJSON(&input)

	initializers.DB.Model(&models.Registry{}).Where("id = ?", id).Update("receiver_id", input.ReceiverID)

	c.JSON(200, registry)
}

func UpdateRegistryShareholder(c *gin.Context) {
	var input models.Registry

	var registry models.Registry
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&registry).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.BindJSON(&input)

	initializers.DB.Model(&models.Registry{}).Where("id = ?", id).Update("shareholder_id", input.ShareholderID)

	c.JSON(200, registry)
}

func UpdateRegistryContract(c *gin.Context) {
	var input models.RegistryContract

	var registry models.Registry
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&registry).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.BindJSON(&input)

	// registry.ContractBuilderShareholderAddress = input.ContractBuilderContractorAreas
	// registry.ContractBuilderContractorAddress = input.ContractBuilderContractorAddress

	registry.BuilderContractorNumber = input.BuilderContractorNumber
	registry.BuilderContractorDate = input.BuilderContractorDate
	registry.BuilderShareholderNumber = input.BuilderShareholderNumber
	registry.BuilderShareholderDate = input.BuilderShareholderDate
	registry.BuilderContractorAdditionalInfo = input.BuilderContractorAdditionalInfo
	registry.BuilderShareholderAdditionalInfo = input.BuilderShareholderAdditionalInfo

	// areaIDs := input.ContractBuilderShareholderAreas

	// Fetch the areas to associate
	// var contractBuilderContractorAreas []models.Area
	// if len(areaIDs) > 0 {
	// 	if err := initializers.DB.Where("code IN ?", areaIDs).Find(&contractBuilderContractorAreas).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch areas"})
	// 		return
	// 	}
	// }

	initializers.DB.Save(&registry)
	c.JSON(200, registry)

	c.JSON(200, registry)
}

type DuplicateTBResult struct {
	ID uint `json:"id" gorm:"column:id"`
	TB int  `json:"t_b" gorm:"column:t_b"`
}

func GetDuplicateTBs(c *gin.Context) {
	var results []DuplicateTBResult

	// Find registries where t_b appears more than once
	// First, get all t_b values that have duplicates
	subQuery := initializers.DB.Model(&models.Registry{}).
		Select("t_b").
		Where("t_b IS NOT NULL").
		Group("t_b").
		Having("COUNT(*) > 1")

	// Then get all registries with those t_b values
	if err := initializers.DB.Model(&models.Registry{}).
		Select("id, t_b").
		Where("t_b IN (?)", subQuery).
		Order("t_b, id").
		Find(&results).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func UpdateRegistryDenial(c *gin.Context) {
	var input models.RegistryDenial

	var registry models.Registry
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&registry).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.BindJSON(&input)

	registry.DenialReason = input.DenialReason
	registry.DenialDate = input.DenialDate
	registry.DenialAdditionalInfo = input.DenialAdditionalInfo

	initializers.DB.Save(&registry)
	c.JSON(200, registry)
}
