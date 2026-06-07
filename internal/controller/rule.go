package controller

import (
	"net/http"

	"github.com/kingkhan77/log-sense/internal/dto"
	"github.com/kingkhan77/log-sense/internal/middleware"
	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type RuleController struct {
	ruleService *service.RuleService
}

func NewRuleController(ruleService *service.RuleService) *RuleController {
	return &RuleController{ruleService: ruleService}
}

func (c *RuleController) CreateRule(ctx *gin.Context) {
	var req dto.CreateRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}

	enabled := true
	if req.IsEnabled != nil {
		enabled = *req.IsEnabled
	}

	rule := &models.AlertRule{
		TenantID:      middleware.TenantID(ctx),
		ServiceID:     req.ServiceID,
		CreatedBy:     middleware.UserID(ctx),
		Name:          req.Name,
		Description:   req.Description,
		Severity:      req.Severity,
		Query:         datatypes.JSON(req.Query),
		Threshold:     req.Threshold,
		WindowMinutes: req.WindowMinutes,
		IsEnabled:     enabled,
	}

	if err := c.ruleService.Create(rule); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, rule)
}

func (c *RuleController) ListRules(ctx *gin.Context) {
	rules, err := c.ruleService.List(middleware.TenantID(ctx))
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, rules)
}

func (c *RuleController) GetRule(ctx *gin.Context) {
	rule, err := c.ruleService.GetByID(middleware.TenantID(ctx), ctx.Param("id"))
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, rule)
}

func (c *RuleController) UpdateRule(ctx *gin.Context) {
	existing, err := c.ruleService.GetByID(middleware.TenantID(ctx), ctx.Param("id"))
	if handleError(ctx, err) {
		return
	}

	var req dto.UpdateRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Severity != "" {
		existing.Severity = req.Severity
	}
	if len(req.Query) > 0 {
		existing.Query = datatypes.JSON(req.Query)
	}
	if req.Threshold > 0 {
		existing.Threshold = req.Threshold
	}
	if req.WindowMinutes > 0 {
		existing.WindowMinutes = req.WindowMinutes
	}
	if req.IsEnabled != nil {
		existing.IsEnabled = *req.IsEnabled
	}

	if err := c.ruleService.Update(middleware.TenantID(ctx), existing); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, existing)
}

func (c *RuleController) DeleteRule(ctx *gin.Context) {
	if handleError(ctx, c.ruleService.Delete(middleware.TenantID(ctx), ctx.Param("id"))) {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "rule deleted"})
}
