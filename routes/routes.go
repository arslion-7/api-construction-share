package routes

import (
	"github.com/arslion-7/api-construction-share/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// r.MaxMultipartMemory = 8 << 20 // 8 MiB

	r.Use(middlewares.CorsMiddleware())

	api := r.Group("/api")
	{
		api.Static("/uploads", "./uploads")
		AuthRoutes(api, "/auth")
		UserRoutes(api, "/users")
		AreaRoutes(api, "/areas")
		GeneralContractorRoutes(api, "/general_contractors")
		BuildingRoutes(api, "/buildings")
		BuilderRoutes(api, "/builders")
		ReceiverRoutes(api, "/receivers")
		ShareholderRoutes(api, "/shareholders")
		ShareholderPropertyRoutes(api, "/shareholder_properties")
		RegistryRoutes(api, "/registries")
	}

	return r
}
