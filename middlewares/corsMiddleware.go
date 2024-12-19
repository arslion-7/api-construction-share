package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from any origin
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// Handle preflight requests (OPTIONS)
		if c.Request.Method == "OPTIONS" {
			// Set additional headers for preflight requests
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Access-Control-Allow-Origin")

			// Respond to preflight request with an empty body and 204 status code
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Continue processing the request
		c.Next()
	}
}
