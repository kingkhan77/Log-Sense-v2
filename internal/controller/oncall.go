package controller

import (
	"net/http"
	"time"

	"github.com/kingkhan77/log-sense/internal/dto"
	"github.com/kingkhan77/log-sense/internal/middleware"
	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/gin-gonic/gin"
)

type OnCallController struct {
	oncallService *service.OnCallService
}

func NewOnCallController(oncallService *service.OnCallService) *OnCallController {
	return &OnCallController{oncallService: oncallService}
}

func (c *OnCallController) CreateSchedule(ctx *gin.Context) {
	var req dto.CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}

	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		abortBadRequest(ctx, "invalid start_time")
		return
	}
	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		abortBadRequest(ctx, "invalid end_time")
		return
	}

	schedule := &models.OnCallSchedule{
		TenantID:  middleware.TenantID(ctx),
		ServiceID: req.ServiceID,
		UserID:    req.UserID,
		StartTime: start,
		EndTime:   end,
	}

	if err := c.oncallService.Create(schedule); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, schedule)
}

func (c *OnCallController) UpdateSchedule(ctx *gin.Context) {
	found, err := c.oncallService.GetByID(middleware.TenantID(ctx), ctx.Param("id"))
	if handleError(ctx, err) {
		return
	}

	var req dto.CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}

	if req.ServiceID != "" {
		found.ServiceID = req.ServiceID
	}
	if req.UserID != "" {
		found.UserID = req.UserID
	}
	if req.StartTime != "" {
		t, err := time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			abortBadRequest(ctx, "invalid start_time")
			return
		}
		found.StartTime = t
	}
	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			abortBadRequest(ctx, "invalid end_time")
			return
		}
		found.EndTime = t
	}

	if err := c.oncallService.Update(middleware.TenantID(ctx), found); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, found)
}

func (c *OnCallController) ListSchedules(ctx *gin.Context) {
	schedules, err := c.oncallService.List(middleware.TenantID(ctx))
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, schedules)
}

func (c *OnCallController) DeleteSchedule(ctx *gin.Context) {
	if err := c.oncallService.Delete(middleware.TenantID(ctx), ctx.Param("id")); err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "schedule deleted"})
}

func (c *OnCallController) GetCurrentOnCall(ctx *gin.Context) {
	schedule, err := c.oncallService.GetCurrentOnCall(middleware.TenantID(ctx), ctx.Param("serviceId"))
	if handleError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, schedule)
}
