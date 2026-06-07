package controller

import (
	"net/http"

	"github.com/kingkhan77/log-sense/internal/dto"
	"github.com/kingkhan77/log-sense/internal/middleware"
	"github.com/kingkhan77/log-sense/internal/service"
	"github.com/gin-gonic/gin"
)

type IngestionController struct {
	ingestionService *service.IngestionService
}

func NewIngestionController(ingestionService *service.IngestionService) *IngestionController {
	return &IngestionController{ingestionService: ingestionService}
}

func (c *IngestionController) IngestLog(ctx *gin.Context) {
	var req dto.IngestLogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		abortBadRequest(ctx, err.Error())
		return
	}

	err := c.ingestionService.PublishLog(
		middleware.TenantID(ctx),
		middleware.ServiceID(ctx),
		service.LogIngestInput{
			Level:     req.Level,
			Message:   req.Message,
			Timestamp: req.Timestamp,
			Metadata:  req.Metadata,
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to ingest log"})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "log ingested"})
}
