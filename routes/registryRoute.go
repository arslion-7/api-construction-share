package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/gin-gonic/gin"
)

func RegistryRoutes(api *gin.RouterGroup, url string) {
	routes := api.Group(url)
	{
		routes.GET("/", controllers.GetRegistries)
		routes.GET("/:id", controllers.GetRegistry)
		routes.POST("/", controllers.CreateRegistry)
		routes.PUT("/:id/update_general_contractor", controllers.UpdateRegistryGeneralContractor)
	}
}
