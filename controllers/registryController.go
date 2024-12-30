package controllers

import (
	"net/http"
	"strconv"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/arslion-7/api-construction-share/utils"
	"github.com/gin-gonic/gin"
)

func GetRegistries(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)
	search := c.Query("search")

	var data []models.Registry
	query := initializers.DB.Model(&models.Registry{}).
		Preload("GeneralContractor").
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
	totalQuery := initializers.DB.Model(&models.Registry{})
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

func GetRegistry(c *gin.Context) {
	id := c.Param("id")

	var registry models.Registry

	if err := initializers.DB.Preload("GeneralContractor").
		First(&registry, id).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, registry)
}

func CreateRegistry(c *gin.Context) {
	var registry models.Registry

	if err := initializers.DB.Save(&registry).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, registry)
}

// type UpdateRegistryGeneralContractorInput struct {
// 	GeneralContractorId `json:general_contractor_id`
// }

func UpdateRegistryGeneralContractor(c *gin.Context) {
	var input models.Registry

	var registry models.Registry
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&registry).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.BindJSON(&input)

	registry.GeneralContractorID = input.GeneralContractorID

	initializers.DB.Save(&registry)
	c.JSON(200, registry)

}

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
