package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/arslion-7/api-construction-share/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

type BuildingAddressInput struct {
	// TB    *int   `json:"t_b"`
	Areas  []uint  `json:"areas"` // IDs of the associated Areas
	Street *string `json:"street"`
}

// CreateBuilding creates a new building with associated areas.
func CreateBuilding(c *gin.Context) {
	// Parse the input JSON
	var input BuildingAddressInput
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

	// Create a new building
	newBuilding := models.Building{
		// TB:    input.TB,
		Areas:  areas, // Set the association
		Street: input.Street,
	}

	if err := initializers.DB.Create(&newBuilding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create building"})
		return
	}

	c.JSON(http.StatusCreated, newBuilding)
}

func UpdateBuildingAddress(c *gin.Context) {
	// Extract building ID from the URL parameters
	id := c.Param("id")

	// Parse the request body (array of area IDs)
	var input BuildingAddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Find the building by ID
	var building models.Building
	if err := initializers.DB.Preload("Areas").Where("id = ?", id).First(&building).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Building not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch building"})
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

	// Update the building's areas (many-to-many relationship)
	if err := initializers.DB.Model(&building).Association("Areas").Replace(&areas); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update areas"})
		return
	}

	fmt.Println("input.Street", input.Street)

	building.Street = input.Street
	if err := initializers.DB.Save(&building).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update street"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Building address updated successfully"})
}
