package controllers

import (
	"strings"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/arslion-7/api-construction-share/utils"
	"github.com/gin-gonic/gin"
)

func GetRegistries(c *gin.Context) {
	// Get pagination parameters
	pagination := utils.GetPaginationParams(c)
	search := c.Query("search")

	// Fetch data
	var data []models.Registry
	query := initializers.DB.Model(&models.Registry{}).
		Limit(pagination.PageSize).
		Offset(pagination.Offset)

	// Apply search filter if a search term is provided
	if search != "" {
		query = query.Where("LOWER(org_name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	// Fetch the data
	if err := query.Find(&data).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Failed to retrieve data"})
		return
	}

	// Get total count with the same search condition
	var total int64
	totalQuery := initializers.DB.Model(&models.Registry{})
	if search != "" {
		totalQuery = totalQuery.Where("LOWER(org_name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}
	totalQuery.Count(&total)

	// Respond with paginated data
	utils.RespondWithPagination(c, data, pagination, total)
}

func GetRegistry(c *gin.Context) {
	id := c.Param("id")

	var registry models.Registry

	if err := initializers.DB.First(&registry, id).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, registry)
}
