package controller

import (
	"net/http"
	"strconv"

	"github.com/kingkhan77/log-sense/internal/middleware"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/gin-gonic/gin"
)

type AlertController struct {
	alertService *service.AlertService
}

func NewAlertController(alertService *service.AlertService) *AlertController {
	return &AlertController{alertService: alertService}
}

const (
	defaultAlertLimit = 50
	maxAlertLimit     = 200
)

func (c *AlertController) ListAlerts(ctx *gin.Context) {
	limit := defaultAlertLimit
	if l, err := strconv.Atoi(ctx.DefaultQuery("limit", "")); err == nil && l > 0 {
		if l > maxAlertLimit {
			l = maxAlertLimit
		}
		limit = l
	}

	offset := 0
	if o, err := strconv.Atoi(ctx.DefaultQuery("offset", "")); err == nil && o >= 0 {
		offset = o
	}

	result, err := c.alertService.List(middleware.TenantID(ctx), limit, offset)
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *AlertController) GetAlert(ctx *gin.Context) {
	alert, err := c.alertService.GetByID(middleware.TenantID(ctx), ctx.Param("id"))
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, alert)
}

func (c *AlertController) AcknowledgeAlert(ctx *gin.Context) {
	if err := c.alertService.Acknowledge(middleware.TenantID(ctx), ctx.Param("id"), middleware.UserID(ctx)); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "alert acknowledged"})
}

func (c *AlertController) ResolveAlert(ctx *gin.Context) {
	if err := c.alertService.Resolve(middleware.TenantID(ctx), ctx.Param("id"), middleware.UserID(ctx)); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "alert resolved"})
}
