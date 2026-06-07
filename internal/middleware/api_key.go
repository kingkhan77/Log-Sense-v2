package middleware

import (
	"net/http"

	"github.com/kingkhan77/log-sense/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func APIKeyAuth(apiKeyRepo *repository.APIKeyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "api key missing"})
			c.Abort()
			return
		}

		key, err := apiKeyRepo.ValidateKey(apiKey)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
			c.Abort()
			return
		}

		c.Set("tenant_id", key.TenantID)
		c.Set("service_id", key.ServiceID)
		c.Next()
	}
}
