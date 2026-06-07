package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func abortNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func abortBadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
}

func handleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if err == gorm.ErrRecordNotFound {
		abortNotFound(c)
		return true
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	return true
}
