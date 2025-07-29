package routes

import (
	"github.com/arslion-7/api-construction-share/controllers"
	"github.com/gin-gonic/gin"
)

func DashboardRoutes(r *gin.Engine) {
	dashboard := r.Group("/api/dashboard")
	{
		dashboard.GET("/stats", controllers.GetDashboardStats)
	}
}
