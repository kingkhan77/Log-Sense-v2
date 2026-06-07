package repository

import (
	"time"

	"github.com/kingkhan77/log-sense/internal/models"
	"gorm.io/gorm"
)

type AlertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) Create(alert *models.Alert) error {
	return r.db.Create(alert).Error
}

func (r *AlertRepository) FindByID(id string) (*models.Alert, error) {
	var alert models.Alert
	err := r.db.Where("id = ?", id).First(&alert).Error
	return &alert, err
}

func (r *AlertRepository) List(tenantID string, limit, offset int) ([]models.Alert, error) {
	var alerts []models.Alert
	err := r.db.Where("tenant_id = ?", tenantID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&alerts).Error
	return alerts, err
}

func (r *AlertRepository) Update(alert *models.Alert) error {
	return r.db.Save(alert).Error
}

func (r *AlertRepository) GetOpenAlertByRule(ruleID string) (*models.Alert, error) {
	var alert models.Alert
	err := r.db.Where(
		"rule_id = ? AND status IN ?",
		ruleID,
		[]string{"OPEN", "ACKNOWLEDGED"},
	).First(&alert).Error
	return &alert, err
}

func (r *AlertRepository) GetOpenAlerts() ([]models.Alert, error) {
	var alerts []models.Alert
	err := r.db.Where("status IN ?", []string{"OPEN", "ACKNOWLEDGED"}).Find(&alerts).Error
	return alerts, err
}

func (r *AlertRepository) ExistsOpenAlert(ruleID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Alert{}).Where(
		"rule_id = ? AND status IN ?",
		ruleID,
		[]string{"OPEN", "ACKNOWLEDGED"},
	).Count(&count).Error
	return count > 0, err
}

func (r *AlertRepository) CountOpenAlerts(tenantID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Alert{}).Where(
		"tenant_id = ? AND status IN ?",
		tenantID,
		[]string{"OPEN", "ACKNOWLEDGED"},
	).Count(&count).Error
	return count, err
}

func (r *AlertRepository) CountSince(tenantID string, since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.Alert{}).Where(
		"tenant_id = ? AND created_at >= ?",
		tenantID,
		since,
	).Count(&count).Error
	return count, err
}

func (r *AlertRepository) Count(tenantID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Alert{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return count, err
}

func (r *AlertRepository) GetOpenAlertsByTenant(tenantID string) ([]models.Alert, error) {
	var alerts []models.Alert
	err := r.db.Where(
		"tenant_id = ? AND status IN ?",
		tenantID,
		[]string{"OPEN", "ACKNOWLEDGED"},
	).Order("created_at desc").Find(&alerts).Error
	return alerts, err
}
