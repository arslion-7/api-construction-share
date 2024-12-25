package controllers

import (
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

func GetGeneralContractor(c *gin.Context) {
	id := c.Param("id")

	var generalContractor models.GeneralContractor

	if err := initializers.DB.First(&generalContractor, id).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, generalContractor)
}

// type CreateGeneralContractorInput struct {
// 	TB      *int    `json:"t_b" binding:"required"`
// 	OrgType *string `json:"org_type" binding:"required"`
// 	OrgName *string `json:"org_name" binding:"required"`
// }

func CreateGeneralContractor(c *gin.Context) {
	// var gc CreateGeneralContractorInput
	var input models.Org
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	newGeneralContractor := models.GeneralContractor{
		Org: models.Org{
			TB:                input.TB,
			OrgType:           input.OrgType,
			OrgName:           input.OrgName,
			HeadPosition:      input.HeadPosition,
			HeadFullName:      input.HeadFullName,
			OrgAdditionalInfo: input.OrgAdditionalInfo,
		},
	}

	if err := initializers.DB.Create(&newGeneralContractor).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
	}

	c.JSON(201, newGeneralContractor)

}
