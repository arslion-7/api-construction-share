package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/gin-gonic/gin"
)

func GeneralContractorRoutes(api *gin.RouterGroup, url string) {
	routes := api.Group(url)
	{
		routes.GET("/", controllers.GetGeneralContractors)
		routes.GET("/:id", controllers.GetGeneralContractor)
		routes.POST("/", controllers.CreateGeneralContractor)
		routes.PUT("/:id", controllers.UpdateGeneralContractorOrg)
		routes.PUT("/:id/cert", controllers.UpdateGeneralContractorCert)
		routes.PUT("/:id/resolution", controllers.UpdateGeneralContractorResolution)
	}
}
