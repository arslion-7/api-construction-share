package controllers

import (
	"net/http"
	"strings"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/arslion-7/api-construction-share/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetReceivers(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)
	search := c.Query("search")

	var data []models.Receiver
	query := initializers.DB.Model(&models.Receiver{}).
		Order("id").
		Limit(pagination.PageSize).
		Offset(pagination.Offset)

	if search != "" {
		// Search by multiple string fields using ILIKE (case-insensitive)
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(citizen_status) LIKE ? OR LOWER(org_name) LIKE ? OR LOWER(department) LIKE ? OR LOWER(position) LIKE ? OR LOWER(firstname) LIKE ? OR LOWER(lastname) LIKE ? OR LOWER(patronymic) LIKE ? OR LOWER(additional_info) LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	if err := query.Find(&data).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve data", "details": err.Error()})
		return
	}

	var total int64
	totalQuery := initializers.DB.Model(&models.Receiver{})
	if search != "" {
		// Same search logic for total count
		searchPattern := "%" + strings.ToLower(search) + "%"
		totalQuery = totalQuery.Where("LOWER(citizen_status) LIKE ? OR LOWER(org_name) LIKE ? OR LOWER(department) LIKE ? OR LOWER(position) LIKE ? OR LOWER(firstname) LIKE ? OR LOWER(lastname) LIKE ? OR LOWER(patronymic) LIKE ? OR LOWER(additional_info) LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}
	totalQuery.Count(&total)

	utils.RespondWithPagination(c, data, pagination, total)
}

func GetReceiver(c *gin.Context) {
	id := c.Param("id")

	var receiver models.Receiver

	if err := initializers.DB.First(&receiver, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Receiver not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receiver)
}

// CreateReceiverInput struct to validate input data
type CreateReceiverInput struct {
	CitizenStatus  string `json:"citizen_status" binding:"required"`
	OrgName        string `json:"org_name" binding:"required"`
	Department     string `json:"department"`
	Position       string `json:"position"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	Patronymic     string `json:"patronymic"`
	AdditionalInfo string `json:"additional_info"`
}

func CreateReceiver(c *gin.Context) {
	// Validate input
	var input CreateReceiverInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}
	// Create Receiver
	receiver := models.Receiver{
		CitizenStatus:  input.CitizenStatus,
		OrgName:        input.OrgName,
		Department:     input.Department,
		Position:       input.Position,
		Firstname:      input.Firstname,
		Lastname:       input.Lastname,
		Patronymic:     input.Patronymic,
		AdditionalInfo: input.AdditionalInfo,
	}

	if err := initializers.DB.Create(&receiver).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to create receiver", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, receiver)
}

// UpdateReceiverInput struct to validate input data
type UpdateReceiverInput struct {
	CitizenStatus  *string `json:"citizen_status"`
	OrgName        *string `json:"org_name"`
	Department     *string `json:"department"`
	Position       *string `json:"position"`
	Firstname      *string `json:"firstname"`
	Lastname       *string `json:"lastname"`
	Patronymic     *string `json:"patronymic"`
	AdditionalInfo *string `json:"additional_info"`
}

func UpdateReceiver(c *gin.Context) {
	// Get id from url
	id := c.Param("id")

	// Get data from db
	var receiver models.Receiver
	if err := initializers.DB.First(&receiver, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Receiver not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve receiver", "details": err.Error()})
		return
	}

	// Validate input
	var input UpdateReceiverInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Update receiver
	if input.CitizenStatus != nil {
		receiver.CitizenStatus = *input.CitizenStatus
	}
	if input.OrgName != nil {
		receiver.OrgName = *input.OrgName
	}
	if input.Department != nil {
		receiver.Department = *input.Department
	}
	if input.Position != nil {
		receiver.Position = *input.Position
	}
	if input.Firstname != nil {
		receiver.Firstname = *input.Firstname
	}
	if input.Lastname != nil {
		receiver.Lastname = *input.Lastname
	}
	if input.Patronymic != nil {
		receiver.Patronymic = *input.Patronymic
	}
	if input.AdditionalInfo != nil {
		receiver.AdditionalInfo = *input.AdditionalInfo
	}

	if err := initializers.DB.Save(&receiver).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to update receiver", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receiver)
}

func DeleteReceiver(c *gin.Context) {
	id := c.Param("id")

	var receiver models.Receiver

	if err := initializers.DB.First(&receiver, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Receiver not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve receiver", "details": err.Error()})
		return
	}

	if err := initializers.DB.Delete(&receiver).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to delete receiver", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Receiver deleted successfully"})
}
