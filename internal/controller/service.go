package controller

import (
	"net/http"

	"github.com/kingkhan77/log-sense/internal/dto"
	"github.com/kingkhan77/log-sense/internal/middleware"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/gin-gonic/gin"
)

type TenantServiceController struct {
	serviceService *service.TenantService
}

func NewServiceController(serviceService *service.TenantService) *TenantServiceController {
	return &TenantServiceController{serviceService: serviceService}
}

func (c *TenantServiceController) CreateService(ctx *gin.Context) {
	var req dto.CreateServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}

	svc, err := c.serviceService.Create(middleware.TenantID(ctx), req.Name, req.Description)
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusCreated, svc)
}

func (c *TenantServiceController) ListServices(ctx *gin.Context) {
	list, err := c.serviceService.List(middleware.TenantID(ctx))
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, list)
}

func (c *TenantServiceController) GetService(ctx *gin.Context) {
	svc, err := c.serviceService.GetByID(middleware.TenantID(ctx), ctx.Param("id"))
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, svc)
}

func (c *TenantServiceController) UpdateService(ctx *gin.Context) {
	var req dto.UpdateServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}

	err := c.serviceService.Update(middleware.TenantID(ctx), ctx.Param("id"), req.Name, req.Description)
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "service updated"})
}

func (c *TenantServiceController) DeleteService(ctx *gin.Context) {
	if handleError(ctx, c.serviceService.Delete(middleware.TenantID(ctx), ctx.Param("id"))) {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "service deleted"})
}
