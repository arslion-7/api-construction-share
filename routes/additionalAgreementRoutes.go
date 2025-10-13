package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/gin-gonic/gin"
)

func AdditionalAgreementRoutes(api *gin.RouterGroup) {
	// Routes under /api/registries/:id/
	// These will be added to existing registry routes

	// Individual agreement routes under /api/additional-agreements
	agreementRoutes := api.Group("/additional-agreements")
	{
		agreementRoutes.GET("/registry/:registryId", controllers.GetAdditionalAgreements)
		agreementRoutes.POST("/", controllers.CreateAdditionalAgreement)
		agreementRoutes.GET("/:id", controllers.GetAdditionalAgreement)
		agreementRoutes.PUT("/:id", controllers.UpdateAdditionalAgreement)
		agreementRoutes.DELETE("/:id", controllers.DeleteAdditionalAgreement)
	}
}
