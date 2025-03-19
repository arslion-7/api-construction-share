package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/arslion-7/api-construction-share/middlewares"
	"github.com/gin-gonic/gin"
)

func ShareholderPropertyRoutes(api *gin.RouterGroup, url string) {
	routes := api.Group(url, middlewares.RequireAuth)
	{
		routes.GET("/", controllers.GetShareholderProperty)
		routes.POST("/", controllers.CreateShareholderProperty)
		// routes.PUT("/:id/update_address", controllers.UpdateBuilderAddress)
		// routes.POST("/", controllers.CreateRegistry)
	}
}
