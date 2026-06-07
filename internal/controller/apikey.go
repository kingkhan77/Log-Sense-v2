package controller

import (
	"net/http"

	"github.com/kingkhan77/log-sense/internal/middleware"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/gin-gonic/gin"
)

type APIKeyController struct {
	svc *service.APIKeyService
}

func NewAPIKeyController(svc *service.APIKeyService) *APIKeyController {
	return &APIKeyController{svc: svc}
}

func (c *APIKeyController) ListKeys(ctx *gin.Context) {
	keys, err := c.svc.List(middleware.TenantID(ctx))
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, keys)
}

func (c *APIKeyController) CreateKey(ctx *gin.Context) {
	var req struct {
		ServiceID string `json:"service_id" binding:"required"`
		Name      string `json:"name"       binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}

	raw, err := c.svc.Create(middleware.TenantID(ctx), req.ServiceID, req.Name)
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"key": raw, "message": "Store this key — it will not be shown again."})
}

func (c *APIKeyController) RevokeKey(ctx *gin.Context) {
	if err := c.svc.Revoke(middleware.TenantID(ctx), ctx.Param("id")); handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "api key revoked"})
}
