package controllers

import (
	"fmt"
	"net/http"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/arslion-7/api-construction-share/utils"
	"github.com/gin-gonic/gin"
)

func GetRegistries(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)
	// search := c.Query("search")

	var data []models.Registry
	query := initializers.DB.Model(&models.Registry{}).
		Preload("GeneralContractor").
		Limit(pagination.PageSize).
		Offset(pagination.Offset)

	if err := query.Find(&data).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve data", "details": err.Error()})
		return
	}

	var total int64
	totalQuery := initializers.DB.Model(&models.Registry{})

	totalQuery.Count(&total)

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

	fmt.Println("input.GeneralContractorID", input.GeneralContractorID)

	registry.GeneralContractorID = input.GeneralContractorID

	initializers.DB.Save(&registry)
	c.JSON(200, registry)

}
