package controllers

import (
	"fmt"
	"strings"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/arslion-7/api-construction-share/utils"
	"github.com/gin-gonic/gin"
)

func GetGeneralContractors(c *gin.Context) {
	// Get pagination parameters
	pagination := utils.GetPaginationParams(c)
	search := c.Query("search")

	fmt.Println("search", search)

	// Fetch data
	var data []models.GeneralContractor
	query := initializers.DB.Model(&models.GeneralContractor{}).
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
	totalQuery := initializers.DB.Model(&models.GeneralContractor{})
	if search != "" {
		totalQuery = totalQuery.Where("LOWER(org_name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}
	totalQuery.Count(&total)

	// Respond with paginated data
	utils.RespondWithPagination(c, data, pagination, total)
}
