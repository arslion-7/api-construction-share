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
		routes.PUT("/:id", controllers.UpdateRegistry)
		routes.PUT("/:id/mail", controllers.UpdateRegistryMail)
		routes.PUT("/:id/registry_number", controllers.UpdateRegistryNumber)
		routes.PUT("/:id/registry_dates", controllers.UpdateRegistryDates)
		routes.PUT("/:id/general_contractor", controllers.UpdateRegistryGeneralContractor)
		routes.PUT("/:id/building", controllers.UpdateRegistryBuilding)
		routes.PUT("/:id/builder", controllers.UpdateRegistryBuilder)
		routes.PUT("/:id/receiver", controllers.UpdateRegistryReceiver)
		routes.PUT("/:id/shareholder", controllers.UpdateRegistryShareholder)
		routes.PUT("/:id/contract", controllers.UpdateRegistryContract)
		routes.PUT("/:id/denial", controllers.UpdateRegistryDenial)
		routes.GET("/duplicate-tbs", controllers.GetDuplicateTBs)
		// routes.GET("/:id/shareholder_property", controllers.GetShareholderProperty)
	}
}
