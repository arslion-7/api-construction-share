package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetShareholderProperty(c *gin.Context) {
	registryId := c.Query("registryId")

	var data models.ShareholderProperty

	if err := initializers.DB.Where("registry_id", registryId).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Shareholder property not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)

}

func CreateShareholderProperty(c *gin.Context) {
	registryId := c.Query("registryId")

	registryID, _ := strconv.Atoi(registryId)

	var input models.ShareholderProperty
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	newShareholderProperty := models.ShareholderProperty{
		BuildingType:   input.BuildingType,
		Part:           input.Part,
		Building:       input.Building,
		Entrance:       input.Entrance,
		Floor:          input.Floor,
		Apartment:      input.Apartment,
		RoomCount:      input.RoomCount,
		Square:         input.Square,
		Price:          input.Price,
		Price1m2:       input.Price1m2,
		AdditionalInfo: input.AdditionalInfo,
		RegistryID:     uint(registryID),
	}

	if err := initializers.DB.Save(&newShareholderProperty).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, newShareholderProperty)
}

func UpdateShareholderProperty(c *gin.Context) {
	id := c.Param("id")

	var input models.ShareholderProperty
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Println("input", input)

	var shareholderProperty models.ShareholderProperty
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&shareholderProperty).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	shareholderProperty.BuildingType = input.BuildingType
	shareholderProperty.Part = input.Part
	shareholderProperty.Building = input.Building
	shareholderProperty.Entrance = input.Entrance
	shareholderProperty.Floor = input.Floor
	shareholderProperty.Apartment = input.Apartment
	shareholderProperty.RoomCount = input.RoomCount
	shareholderProperty.Square = input.Square
	shareholderProperty.Price = input.Price
	shareholderProperty.Price1m2 = input.Price1m2
	shareholderProperty.AdditionalInfo = input.AdditionalInfo

	if err := initializers.DB.Save(&shareholderProperty).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, shareholderProperty)
}
