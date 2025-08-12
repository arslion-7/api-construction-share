package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/gin-gonic/gin"
)

func OldRegistryRoutes(router *gin.Engine) {
	oldRegistryGroup := router.Group("/api/old-registries")
	{
		oldRegistryGroup.GET("/", controllers.GetOldRegistries)
		oldRegistryGroup.GET("/:id", controllers.GetOldRegistry)
		oldRegistryGroup.PUT("/:id", controllers.UpdateOldRegistry)
	}
}
