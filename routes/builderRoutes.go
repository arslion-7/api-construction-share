package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/arslion-7/api-construction-share/middlewares"
	"github.com/gin-gonic/gin"
)

func BuilderRoutes(api *gin.RouterGroup, url string) {
	routes := api.Group(url, middlewares.RequireAuth)
	{
		routes.GET("/", controllers.GetBuilders)
		routes.GET("/:id", controllers.GetBuilder)
		// routes.POST("/", controllers.CreateBuilding)
		// routes.POST("/", controllers.CreateRegistry)
	}
}
