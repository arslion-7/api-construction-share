package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/arslion-7/api-construction-share/middlewares"
	"github.com/arslion-7/api-construction-share/utils"
	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup, url string) {
	users := api.Group(url, middlewares.RequireAuth, utils.CheckUserRole(utils.Role.Admin))
	{
		users.GET("/", controllers.GetUsers)
		users.GET("/:id", controllers.GetUser)
		users.POST("", controllers.CreateUser)
		users.PUT("/:id", controllers.UpdateUser)
		users.DELETE("/:id", controllers.DeleteUser)
		users.PUT("/:id/restore", controllers.RestoreUser)
	}
}
