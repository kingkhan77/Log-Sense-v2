package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Admin() gin.HandlerFunc {

	return func(c *gin.Context) {

		role, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}