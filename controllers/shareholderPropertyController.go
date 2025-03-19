package controllers

import (
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
