package controllers

import (
	"net/http"
	"strconv"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/arslion-7/api-construction-share/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetBuilders(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)
	search := c.Query("search")

	var data []models.Builder
	query := initializers.DB.Model(&models.Builder{}).Preload("Areas").
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
	totalQuery := initializers.DB.Model(&models.Builder{})
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

func GetBuilder(c *gin.Context) {
	id := c.Param("id")

	var builder models.Builder

	if err := initializers.DB.Preload("Areas").First(&builder, id).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, builder)
}

type BuilderAddressInput struct {
	// TB    *int   `json:"t_b"`
	Areas                 []uint  `json:"areas"` // IDs of the associated Areas
	Address               *string `gorm:"type:varchar(510);index" json:"address"`
	AddressAdditionalInfo *string `gorm:"type:varchar(510);index" json:"address_additional_info"`
}

// CreateBuilder creates a new builder with associated areas.
func CreateBuilder(c *gin.Context) {
	// Parse the input JSON
	var input BuilderAddressInput
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

	// Create a new builder
	newBuilder := models.Builder{
		// TB:    input.TB,
		BuilderAddress: models.BuilderAddress{
			Areas:                 areas, // Set the association
			Address:               input.Address,
			AddressAdditionalInfo: input.AddressAdditionalInfo,
		},
	}

	if err := initializers.DB.Create(&newBuilder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create builder"})
		return
	}

	c.JSON(http.StatusCreated, newBuilder)
}

func UpdateBuilderAddress(c *gin.Context) {
	// Extract builder ID from the URL parameters
	id := c.Param("id")

	// Parse the request body (array of area IDs)
	var input BuilderAddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Find the builder by ID
	var builder models.Builder
	if err := initializers.DB.Preload("Areas").Where("id = ?", id).First(&builder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Builder not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch builder"})
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

	// Update the builder's areas (many-to-many relationship)
	if err := initializers.DB.Model(&builder).Association("Areas").Replace(&areas); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update areas"})
		return
	}

	builder.Address = input.Address
	builder.AddressAdditionalInfo = input.AddressAdditionalInfo
	if err := initializers.DB.Save(&builder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update street"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Builder address updated successfully"})
}

// // func UpdateBuilderMain(c *gin.Context) {
// // 	var builder models.Builder
// // 	id := c.Params.ByName("id")
// // 	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&builder).Error; err != nil {
// // 		c.AbortWithStatus(404)
// // 		return
// // 	}

// // 	var input models.BuilderMain

// // 	c.BindJSON(&input)

// // 	builder.TB = input.TB
// // 	builder.IdentNumber = input.IdentNumber
// // 	builder.Kind = input.Kind
// // 	builder.Price = input.Price
// // 	builder.Percentage = input.Percentage
// // 	builder.StartDate = input.StartDate
// // 	builder.EndDate = input.EndDate

// // 	initializers.DB.Save(&builder)
// // 	c.JSON(200, builder)
// // }
