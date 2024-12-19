package controllers

import (
	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/gin-gonic/gin"
)

func GetGeneralContractors(c *gin.Context) {
	// show := c.Query("show")
	var generalContractors []models.GeneralContractor

	if err := initializers.DB.Find(&generalContractors).Error; err != nil {
		c.AbortWithStatus(400)
		return
	}
	c.JSON(200, generalContractors)
}
