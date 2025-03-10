package controllers

import (
	"net/http"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/arslion-7/api-construction-share/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetShareholders(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)
	search := c.Query("search")

	var data []models.Shareholder
	query := initializers.DB.Model(&models.Shareholder{}).
		Limit(pagination.PageSize).
		Offset(pagination.Offset)

	if search != "" {
		// Search by multiple string fields using ILIKE (case-insensitive)
		// searchPattern := "%" + strings.ToLower(search) + "%"
		// query = query.Where("LOWER(citizen_status) LIKE ? OR LOWER(org_name) LIKE ? OR LOWER(department) LIKE ? OR LOWER(position) LIKE ? OR LOWER(firstname) LIKE ? OR LOWER(lastname) LIKE ? OR LOWER(patronymic) LIKE ? OR LOWER(additional_info) LIKE ?",
		// 	searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	if err := query.Find(&data).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve data", "details": err.Error()})
		return
	}

	var total int64
	totalQuery := initializers.DB.Model(&models.Shareholder{})
	if search != "" {
		// Same search logic for total count
		// searchPattern := "%" + strings.ToLower(search) + "%"
		// totalQuery = totalQuery.Where("LOWER(citizen_status) LIKE ? OR LOWER(org_name) LIKE ? OR LOWER(department) LIKE ? OR LOWER(position) LIKE ? OR LOWER(firstname) LIKE ? OR LOWER(lastname) LIKE ? OR LOWER(patronymic) LIKE ? OR LOWER(additional_info) LIKE ?",
		// 	searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}
	totalQuery.Count(&total)

	utils.RespondWithPagination(c, data, pagination, total)
}

func GetShareholder(c *gin.Context) {
	id := c.Param("id")

	var data models.Shareholder

	if err := initializers.DB.First(&data, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Shareholder not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
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
