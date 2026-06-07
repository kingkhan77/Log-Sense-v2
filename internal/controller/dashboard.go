package controller

import (
	"net/http"

	"github.com/kingkhan77/log-sense/internal/middleware"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	dashboardService *service.DashboardService
}

func NewDashboardController(dashboardService *service.DashboardService) *DashboardController {
	return &DashboardController{dashboardService: dashboardService}
}

func (c *DashboardController) Summary(ctx *gin.Context) {
	summary, err := c.dashboardService.Summary(middleware.TenantID(ctx))
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, summary)
}
