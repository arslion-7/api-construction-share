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
		GeneralContractorRoutes(api, "/general_contractors")
		RegistryRoutes(api, "/registries")
	}

	return r
}
