package repository

import (
	"github.com/kingkhan77/log-sense/internal/models"
	"gorm.io/gorm"
)

type RuleRepository struct {
	db *gorm.DB
}

func NewRuleRepository(db *gorm.DB) *RuleRepository {
	return &RuleRepository{db: db}
}

func (r *RuleRepository) Create(rule *models.AlertRule) error {
	return r.db.Create(rule).Error
}

func (r *RuleRepository) FindByID(id string) (*models.AlertRule, error) {
	var rule models.AlertRule
	err := r.db.Where("id = ?", id).First(&rule).Error
	return &rule, err
}

func (r *RuleRepository) List(tenantID string) ([]models.AlertRule, error) {
	var rules []models.AlertRule
	err := r.db.Where("tenant_id = ?", tenantID).Find(&rules).Error
	return rules, err
}

func (r *RuleRepository) ListEnabled() ([]models.AlertRule, error) {
	var rules []models.AlertRule
	err := r.db.Where("is_enabled = ?", true).Find(&rules).Error
	return rules, err
}

func (r *RuleRepository) Update(rule *models.AlertRule) error {
	return r.db.Save(rule).Error
}

func (r *RuleRepository) Delete(id string) error {
	return r.db.Delete(&models.AlertRule{}, "id = ?", id).Error
}
