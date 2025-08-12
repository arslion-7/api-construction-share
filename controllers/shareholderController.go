package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/arslion-7/api-construction-share/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetShareholders(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)
	search := c.Query("search")
	lowerSearch := strings.ToLower(search)

	var data []models.Shareholder
	query := initializers.DB.Model(&models.Shareholder{}).
		Preload("Areas").
		Preload("Phones").
		Order("id").
		Limit(pagination.PageSize).
		Offset(pagination.Offset)

	if search != "" {
		// Search by organization fields and document fields
		searchPattern := "%" + lowerSearch + "%"
		query = query.Where(`
			LOWER(org_type) LIKE ? OR 
			LOWER(org_name) LIKE ? OR 
			LOWER(head_position) LIKE ? OR 
			LOWER(head_full_name) LIKE ? OR 
			LOWER(org_additional_info) LIKE ? OR
			LOWER(passport_series) LIKE ? OR
			CAST(passport_number AS TEXT) LIKE ? OR
			LOWER(patent_series) LIKE ? OR
			CAST(patent_number AS TEXT) LIKE ? OR
			CAST(cert_number AS TEXT) LIKE ? OR
			LOWER(docs_additional_info) LIKE ?
		`,
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	if err := query.Find(&data).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve data", "details": err.Error()})
		return
	}

	var total int64
	totalQuery := initializers.DB.Model(&models.Shareholder{})
	if search != "" {
		// Same search logic for total count
		searchPattern := "%" + lowerSearch + "%"
		totalQuery = totalQuery.Where(`
			LOWER(org_type) LIKE ? OR 
			LOWER(org_name) LIKE ? OR 
			LOWER(head_position) LIKE ? OR 
			LOWER(head_full_name) LIKE ? OR 
			LOWER(org_additional_info) LIKE ? OR
			LOWER(passport_series) LIKE ? OR
			CAST(passport_number AS TEXT) LIKE ? OR
			LOWER(patent_series) LIKE ? OR
			CAST(patent_number AS TEXT) LIKE ? OR
			CAST(cert_number AS TEXT) LIKE ? OR
			LOWER(docs_additional_info) LIKE ?
		`,
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}
	totalQuery.Count(&total)

	utils.RespondWithPagination(c, data, pagination, total)
}

func GetShareholder(c *gin.Context) {
	id := c.Param("id")

	var data models.Shareholder

	if err := initializers.DB.Preload("Areas").Preload("Phones").First(&data, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Shareholder not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// CreateShareholder creates a new shareholder with associated areas.
func CreateShareholder(c *gin.Context) {
	// Parse the input JSON
	var input AddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch the areas to associate
	var areas []models.Area
	if len(input.Areas) > 0 {
		if err := initializers.DB.Where("code IN ?", input.Areas).Find(&areas).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch areas"})
			return
		}
	}

	// Create a new shareholder
	newShareholder := models.Shareholder{
		// TB:    input.TB,
		ShareholderAddress: models.ShareholderAddress{
			Areas:                 areas, // Set the association
			Address:               input.Address,
			AddressAdditionalInfo: input.AddressAdditionalInfo,
		},
	}

	if err := initializers.DB.Create(&newShareholder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shareholder"})
		return
	}

	c.JSON(http.StatusCreated, newShareholder)
}

func UpdateShareholderAddress(c *gin.Context) {
	// Extract shareholder ID from the URL parameters
	id := c.Param("id")

	// Parse the request body (array of area IDs)
	var input AddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Find the shareholder by ID
	var shareholder models.Shareholder
	if err := initializers.DB.Preload("Areas").Where("id = ?", id).First(&shareholder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shareholder not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shareholder"})
		return
	}

	areaIDs := input.Areas

	// Fetch the areas to associate
	var areas []models.Area
	if len(areaIDs) > 0 {
		if err := initializers.DB.Where("code IN ?", areaIDs).Find(&areas).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch areas"})
			return
		}
	}

	// Update the shareholder's areas (many-to-many relationship)
	if err := initializers.DB.Model(&shareholder).Association("Areas").Replace(&areas); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update areas"})
		return
	}

	shareholder.Address = input.Address
	shareholder.AddressAdditionalInfo = input.AddressAdditionalInfo
	if err := initializers.DB.Save(&shareholder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update street"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shareholder address updated successfully"})
}

func UpdateShareholderDocs(c *gin.Context) {
	var data models.Shareholder
	id := c.Param("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&data).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	var input models.ShareholderDocs

	c.BindJSON(&input)

	fmt.Println("input", input)

	data.PassportSeries = input.PassportSeries
	data.PassportNumber = input.PassportNumber
	data.PatentSeries = input.PatentSeries
	data.PatentNumber = input.PatentNumber
	data.CertNumber = input.CertNumber
	data.DocsAdditionalInfo = input.DocsAdditionalInfo

	initializers.DB.Save(&data)
	c.JSON(200, data)
}

func UpdateShareholderOrg(c *gin.Context) {
	var data models.Shareholder
	id := c.Param("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&data).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	var input models.Org

	c.BindJSON(&input)

	data.OrgType = input.OrgType
	data.OrgName = input.OrgName
	data.HeadPosition = input.HeadPosition
	data.HeadFullName = input.HeadFullName
	data.OrgAdditionalInfo = input.OrgAdditionalInfo

	initializers.DB.Save(&data)
	c.JSON(200, data)
}

type PhoneInput struct {
	Kind   *string `json:"kind" gorm:"size:3"`
	Number *string `json:"number" gorm:"size:12"`
	Owner  *string `json:"owner" gorm:"size:125"`
}

type PhoneInputWrapper struct {
	Phones []PhoneInput `json:"phones"` // Ensure JSON tag matches the expected input structure
}

func UpdateShareholderPhones(c *gin.Context) {
	var input PhoneInputWrapper
	id := c.Param("id")

	// Find existing phone records
	var existingPhones []models.Phone
	if err := initializers.DB.Where("shareholder_id = ?", id).Find(&existingPhones).Error; err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "Phones not found"})
		return
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	uintID, _ := strconv.Atoi(id)

	// Update or create phone records
	var updatedPhones []models.Phone
	for _, phone := range input.Phones {
		updatedPhones = append(updatedPhones, models.Phone{
			Kind:          phone.Kind,
			Number:        phone.Number,
			Owner:         phone.Owner,
			ShareholderID: uint(uintID),
		})
	}

	// Delete old records and insert new ones
	tx := initializers.DB.Begin()
	if err := tx.Where("shareholder_id = ?", id).Delete(&models.Phone{}).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(500, gin.H{"error": "Failed to delete old phone records"})
		return
	}
	if err := tx.Create(&updatedPhones).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(500, gin.H{"error": "Failed to update phone records"})
		return
	}
	tx.Commit()

	c.JSON(200, updatedPhones)
}

func DeleteShareholder(c *gin.Context) {
	id := c.Param("id")

	var data models.Shareholder

	if err := initializers.DB.First(&data, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Shareholder not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve shareholder", "details": err.Error()})
		return
	}

	if err := initializers.DB.Delete(&data).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to delete shareholder", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shareholder deleted successfully"})
}
