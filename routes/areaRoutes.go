package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/arslion-7/api-construction-share/middlewares"
	"github.com/gin-gonic/gin"
)

func AreaRoutes(api *gin.RouterGroup, url string) {
	routes := api.Group(url, middlewares.RequireAuth)
	{
		routes.GET("/area_hierarchy", controllers.GetAreaHierarchy)
		routes.GET("/area_hierarchy/:code", controllers.FetchAreaHierarchy)
		// routes.POST("/", controllers.CreateRegistry)
	}
}
