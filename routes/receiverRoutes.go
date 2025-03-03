package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/arslion-7/api-construction-share/middlewares"
	"github.com/gin-gonic/gin"
)

func ReceiverRoutes(api *gin.RouterGroup, url string) {
	routes := api.Group(url, middlewares.RequireAuth)
	{
		routes.GET("/", controllers.GetReceivers)
		routes.GET("/:id", controllers.GetReceiver)
		routes.POST("/", controllers.CreateReceiver)
		routes.PUT("/:id", controllers.UpdateReceiver)
	}
}
