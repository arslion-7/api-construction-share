package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/arslion-7/api-construction-share/middlewares"
	"github.com/gin-gonic/gin"
)

func RegistryRoutes(api *gin.RouterGroup, url string) {
	routes := api.Group(url, middlewares.RequireAuth)
	{
		routes.GET("/", controllers.GetRegistries)
		routes.GET("/:id", controllers.GetRegistry)
		routes.POST("/", controllers.CreateRegistry)
		routes.PUT("/:id/update_registry_number", controllers.UpdateRegistryNumber)
		routes.PUT("/:id/update_general_contractor", controllers.UpdateRegistryGeneralContractor)
		routes.PUT("/:id/update_building", controllers.UpdateRegistryBuilding)
		routes.PUT("/:id/update_builder", controllers.UpdateRegistryBuilder)
		routes.PUT("/:id/update_receiver", controllers.UpdateRegistryReceiver)
	}
}
