package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc{
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		roleStr, ok := role.(string)
		if !ok || roleStr != "admin"{
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden - Admin only!",
			})
			c.Abort()
			return
		}

		c.Next()
	}

}