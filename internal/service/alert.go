package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kingkhan77/log-sense/internal/constants"
	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/internal/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AlertService struct {
	repo  *repository.AlertRepository
	redis *redis.Client
}

func NewAlertService(
	repo *repository.AlertRepository,
	redis *redis.Client,
) *AlertService {
	return &AlertService{repo: repo, redis: redis}
}

type AlertListResult struct {
	Total  int64          `json:"total"`
	Alerts []models.Alert `json:"alerts"`
}

func (s *AlertService) List(tenantID string, limit, offset int) (*AlertListResult, error) {
	total, err := s.repo.Count(tenantID)
	if err != nil {
		return nil, err
	}
	alerts, err := s.repo.List(tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	return &AlertListResult{Total: total, Alerts: alerts}, nil
}

func (s *AlertService) GetByID(tenantID, id string) (*models.Alert, error) {
	alert, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if alert.TenantID != tenantID {
		return nil, gorm.ErrRecordNotFound
	}
	return alert, nil
}

func (s *AlertService) Acknowledge(tenantID, alertID, userID string) error {
	alert, err := s.GetByID(tenantID, alertID)
	if err != nil {
		return err
	}
	if alert.Status != constants.AlertOpen {
		return errors.New("alert is not open")
	}

	now := time.Now()
	alert.Status = constants.AlertAcknowledged
	alert.AcknowledgedAt = &now
	alert.AcknowledgedBy = &userID
	alert.UpdatedAt = now

	if err := s.repo.Update(alert); err != nil {
		return err
	}
	s.clearRuleCache(alert.RuleID)
	return nil
}

func (s *AlertService) Resolve(tenantID, alertID, userID string) error {
	alert, err := s.GetByID(tenantID, alertID)
	if err != nil {
		return err
	}
	if alert.Status == constants.AlertResolved {
		return errors.New("alert already resolved")
	}

	now := time.Now()
	alert.Status = constants.AlertResolved
	alert.ResolvedAt = &now
	alert.ResolvedBy = &userID
	alert.UpdatedAt = now

	if err := s.repo.Update(alert); err != nil {
		return err
	}
	s.clearRuleCache(alert.RuleID)
	return nil
}

func (s *AlertService) clearRuleCache(ruleID string) {
	if s.redis == nil {
		return
	}
	_ = s.redis.Del(context.Background(), fmt.Sprintf("alert:rule:%s", ruleID)).Err()
}
