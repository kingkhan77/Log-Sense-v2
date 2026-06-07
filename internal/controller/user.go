package controller

import (
	"net/http"

	"github.com/kingkhan77/log-sense/internal/dto"
	"github.com/kingkhan77/log-sense/internal/middleware"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) CreateDeveloper(ctx *gin.Context) {
	var req dto.CreateDeveloperRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}

	tenantID := middleware.TenantID(ctx)
	if err := c.userService.CreateDeveloper(tenantID, req.Name, req.Email, req.Password); err != nil {
		if handleError(ctx, err) {
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "developer created"})
}

func (c *UserController) ListDevelopers(ctx *gin.Context) {
	tenantID := middleware.TenantID(ctx)
	users, err := c.userService.ListDevelopers(tenantID)
	if handleError(ctx, err) {
		return
	}

	resp := make([]dto.DeveloperResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, dto.DeveloperResponse{
			ID: u.ID, Name: u.Name, Email: u.Email, Role: u.Role,
		})
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *UserController) UpdateDeveloper(ctx *gin.Context) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}

	tenantID := middleware.TenantID(ctx)
	if err := c.userService.UpdateDeveloper(tenantID, ctx.Param("id"), req.Name, req.Email); handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "developer updated"})
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password"     binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}
	if err := c.userService.ChangePassword(middleware.UserID(ctx), req.CurrentPassword, req.NewPassword); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

func (c *UserController) DeleteDeveloper(ctx *gin.Context) {
	tenantID := middleware.TenantID(ctx)
	if err := c.userService.DeactivateDeveloper(tenantID, ctx.Param("id")); handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "developer deactivated"})
}
