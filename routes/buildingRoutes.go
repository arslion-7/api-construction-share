package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/arslion-7/api-construction-share/middlewares"
	"github.com/gin-gonic/gin"
)

func BuildingRoutes(api *gin.RouterGroup, url string) {
	routes := api.Group(url, middlewares.RequireAuth)
	{
		routes.GET("/", controllers.GetBuildings)
		routes.GET("/:id", controllers.GetBuilding)
		routes.POST("/", controllers.CreateBuilding)
		routes.PUT("/:id/update_address", controllers.UpdateBuildingAddress)
		routes.PUT("/:id/update_main", controllers.UpdateBuildingMain)
		routes.PUT("/:id/update_order", controllers.UpdateBuildingOrder)
		routes.PUT("/:id/update_cert", controllers.UpdateBuildingCert)
		routes.PUT("/:id/update_square", controllers.UpdateBuildingSquare)
		// routes.POST("/", controllers.CreateRegistry)
	}
}
