package controller

import (
	"net/http"

	"github.com/kingkhan77/log-sense/internal/dto"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}

	token, err := c.authService.Login(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.LoginResponse{Token: token})
}
