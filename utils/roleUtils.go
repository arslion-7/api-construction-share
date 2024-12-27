package utils

import (
	"github.com/arslion-7/api-construction-share/models"
	"github.com/gin-gonic/gin"
)

// Predefined roles for validation and consistency
type Roles struct {
	Admin string
	User  string
}

// Global instance of Roles
var Role = Roles{
	Admin: "admin",
	User:  "user",
}

func CheckUserRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Example: Retrieve user role from request header or JWT token
		user, _ := c.Get("user")
		typedUser := user.(models.User)

		if *typedUser.Role == role {
			c.Next()
		} else {
			c.JSON(403, gin.H{
				"error": "Siziň bu hereketi ýerine ýetirmäge mümkinçiligiňiz ýok",
			})
			c.Abort()
		}
	}
}
