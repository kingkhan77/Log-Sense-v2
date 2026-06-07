package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RuleService struct {
	repo *repository.RuleRepository
}

func NewRuleService(repo *repository.RuleRepository) *RuleService {
	return &RuleService{repo: repo}
}

type RuleQueryInput struct {
	Level           string            `json:"level"`
	MessageContains string            `json:"message_contains"`
	Fields          map[string]string `json:"fields"`
}

func (s *RuleService) Create(rule *models.AlertRule) error {
	if err := validateRuleQuery(rule.Query); err != nil {
		return err
	}
	rule.ID = uuid.NewString()
	now := time.Now()
	rule.CreatedAt = now
	rule.UpdatedAt = now
	return s.repo.Create(rule)
}

func (s *RuleService) List(tenantID string) ([]models.AlertRule, error) {
	return s.repo.List(tenantID)
}

func (s *RuleService) GetByID(tenantID, id string) (*models.AlertRule, error) {
	rule, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if rule.TenantID != tenantID {
		return nil, gorm.ErrRecordNotFound
	}
	return rule, nil
}

func (s *RuleService) Update(tenantID string, rule *models.AlertRule) error {
	existing, err := s.GetByID(tenantID, rule.ID)
	if err != nil {
		return err
	}
	if err := validateRuleQuery(rule.Query); err != nil {
		return err
	}
	rule.TenantID = existing.TenantID
	rule.CreatedBy = existing.CreatedBy
	rule.CreatedAt = existing.CreatedAt
	rule.UpdatedAt = time.Now()
	return s.repo.Update(rule)
}

func (s *RuleService) Delete(tenantID, id string) error {
	if _, err := s.GetByID(tenantID, id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func validateRuleQuery(raw []byte) error {
	if len(raw) == 0 {
		return errors.New("query is required")
	}
	var q RuleQueryInput
	if err := json.Unmarshal(raw, &q); err != nil {
		return errors.New("invalid query JSON")
	}
	return nil
}
