package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/arslion-7/api-construction-share/middlewares"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(api *gin.RouterGroup, url string) {
	auth := api.Group(url)
	{
		// auth.POST("/sign-up", controllers.SignUp)
		auth.POST("/sign-in", controllers.SignIn)
		auth.GET("/validate", middlewares.RequireAuth, controllers.Validate)
		auth.PUT("/reset-password/:id", middlewares.RequireAuth, controllers.ResetPassword)
		auth.GET("/me", middlewares.RequireAuth, controllers.GetMe)
	}
}
