package controllers

import (
	"net/http"
	"time"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAdditionalAgreements handles GET request for additional agreements by registry ID
func GetAdditionalAgreements(c *gin.Context) {
	registryID := c.Param("registryId")

	var agreements []models.AdditionalAgreement

	if err := initializers.DB.Where("registry_id = ?", registryID).
		Order("agreement_date ASC").
		Find(&agreements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch additional agreements"})
		return
	}

	c.JSON(200, gin.H{
		"data": agreements,
	})
}

// GetAdditionalAgreement handles GET request for a single additional agreement
func GetAdditionalAgreement(c *gin.Context) {
	id := c.Param("id")
	var agreement models.AdditionalAgreement

	if err := initializers.DB.First(&agreement, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Additional agreement not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch additional agreement"})
		return
	}

	c.JSON(200, agreement)
}

// CreateAdditionalAgreement handles POST request to create a new additional agreement
func CreateAdditionalAgreement(c *gin.Context) {
	var requestBody struct {
		RegistryID      uint   `json:"registry_id" binding:"required"`
		AgreementNumber string `json:"agreement_number"`
		AgreementDate   string `json:"agreement_date" binding:"required"`
		Reason          string `json:"reason"`
		AdditionalInfo  string `json:"additional_info"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Verify registry exists
	var registry models.Registry
	if err := initializers.DB.First(&registry, requestBody.RegistryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registry not found"})
		return
	}

	// Parse date
	parsedDate, err := time.Parse("2006-01-02", requestBody.AgreementDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	agreement := models.AdditionalAgreement{
		RegistryID:      requestBody.RegistryID,
		AgreementNumber: requestBody.AgreementNumber,
		AgreementDate:   &parsedDate,
		Reason:          requestBody.Reason,
		AdditionalInfo:  requestBody.AdditionalInfo,
	}

	if err := initializers.DB.Create(&agreement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create additional agreement"})
		return
	}

	c.JSON(201, gin.H{
		"message": "Additional agreement created successfully",
		"data":    agreement,
	})
}

// UpdateAdditionalAgreement handles PUT request to update an additional agreement
func UpdateAdditionalAgreement(c *gin.Context) {
	id := c.Param("id")

	// Find the record first
	var agreement models.AdditionalAgreement
	if err := initializers.DB.First(&agreement, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Additional agreement not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch additional agreement"})
		return
	}

	// Parse request body
	var requestBody struct {
		AgreementNumber string `json:"agreement_number"`
		AgreementDate   string `json:"agreement_date" binding:"required"`
		Reason          string `json:"reason"`
		AdditionalInfo  string `json:"additional_info"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Parse date
	parsedDate, err := time.Parse("2006-01-02", requestBody.AgreementDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Update the record
	updates := map[string]interface{}{
		"agreement_number": requestBody.AgreementNumber,
		"agreement_date":   parsedDate,
		"reason":           requestBody.Reason,
		"additional_info":  requestBody.AdditionalInfo,
		"updated_at":       time.Now(),
	}

	if err := initializers.DB.Model(&agreement).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update additional agreement"})
		return
	}

	// Fetch the updated record
	if err := initializers.DB.First(&agreement, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated record"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Additional agreement updated successfully",
		"data":    agreement,
	})
}

// DeleteAdditionalAgreement handles DELETE request to delete an additional agreement
func DeleteAdditionalAgreement(c *gin.Context) {
	id := c.Param("id")

	// Find the record first
	var agreement models.AdditionalAgreement
	if err := initializers.DB.First(&agreement, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Additional agreement not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch additional agreement"})
		return
	}

	// Delete the record (soft delete)
	if err := initializers.DB.Delete(&agreement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete additional agreement"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Additional agreement deleted successfully",
	})
}
