package controllers

import (
	"net/http"
	"strconv"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/arslion-7/api-construction-share/utils"
	"github.com/gin-gonic/gin"
)

func GetBuildings(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)
	search := c.Query("search")

	var data []models.Building
	query := initializers.DB.Model(&models.Building{}).Preload("Areas").
		Limit(pagination.PageSize).
		Offset(pagination.Offset)

	if search != "" {
		if tb, err := strconv.Atoi(search); err == nil {
			query = query.Where("t_b = ?", tb) // Search by integer value
		} else {
			// Handle cases where search is not a number, e.g., search by other fields if needed.
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid search parameter. t_b must be an integer."})
			return
		}
	}

	if err := query.Find(&data).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve data", "details": err.Error()})
		return
	}

	var total int64
	totalQuery := initializers.DB.Model(&models.Building{})
	if search != "" {
		if tb, err := strconv.Atoi(search); err == nil {
			totalQuery = totalQuery.Where("t_b = ?", tb)
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid search parameter. t_b must be an integer."})
			return
		}
	}
	totalQuery.Count(&total)

	utils.RespondWithPagination(c, data, pagination, total)

}

func GetBuilding(c *gin.Context) {
	id := c.Param("id")

	var building models.Building

	if err := initializers.DB.Preload("Areas").First(&building, id).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, building)
}
