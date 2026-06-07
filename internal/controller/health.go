package controller

import (
	"context"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/kingkhan77/log-sense/pkg"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type HealthController struct {
	db          *gorm.DB
	redis       *redis.Client
	kafkaClient sarama.Client
}

// NewHealthController creates a persistent Kafka client once so that each
// health check reuses the existing connection instead of opening a new one.
func NewHealthController(db *gorm.DB, redis *redis.Client, cfg *pkg.Config) *HealthController {
	kafkaClient, err := sarama.NewClient(cfg.Kafka.Brokers, sarama.NewConfig())
	if err != nil {
		log.Warn().Err(err).Msg("health controller: kafka client unavailable at startup")
	}
	return &HealthController{db: db, redis: redis, kafkaClient: kafkaClient}
}

func (h *HealthController) Health(ctx *gin.Context) {
	status := "UP"
	checks := gin.H{}

	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		status = "DOWN"
		checks["postgres"] = "DOWN"
	} else {
		checks["postgres"] = "UP"
	}

	if err := h.redis.Ping(context.Background()).Err(); err != nil {
		status = "DOWN"
		checks["redis"] = "DOWN"
	} else {
		checks["redis"] = "UP"
	}

	if h.kafkaClient == nil || h.kafkaClient.Closed() {
		status = "DOWN"
		checks["kafka"] = "DOWN"
	} else if err := h.kafkaClient.RefreshMetadata(); err != nil {
		status = "DOWN"
		checks["kafka"] = "DOWN"
	} else {
		checks["kafka"] = "UP"
	}

	code := http.StatusOK
	if status == "DOWN" {
		code = http.StatusServiceUnavailable
	}

	ctx.JSON(code, gin.H{"status": status, "checks": checks})
}
