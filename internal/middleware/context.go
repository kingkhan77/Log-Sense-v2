package middleware

import "github.com/gin-gonic/gin"

func TenantID(c *gin.Context) string {
	v, _ := c.Get("tenant_id")
	s, _ := v.(string)
	return s
}

func UserID(c *gin.Context) string {
	v, _ := c.Get("user_id")
	s, _ := v.(string)
	return s
}

func Role(c *gin.Context) string {
	v, _ := c.Get("role")
	s, _ := v.(string)
	return s
}

func ServiceID(c *gin.Context) string {
	v, _ := c.Get("service_id")
	s, _ := v.(string)
	return s
}
