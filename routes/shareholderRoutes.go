package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/arslion-7/api-construction-share/middlewares"
	"github.com/gin-gonic/gin"
)

func ShareholderRoutes(api *gin.RouterGroup, url string) {
	routes := api.Group(url, middlewares.RequireAuth)
	{
		routes.GET("/", controllers.GetShareholders)
		routes.GET("/:id", controllers.GetShareholder)
		routes.POST("/", controllers.CreateShareholder)
		routes.PUT("/:id/update_address", controllers.UpdateShareholderAddress)
		routes.PUT("/:id/update_docs", controllers.UpdateShareholderDocs)
	}
}
