package repository

import (
	"time"

	"github.com/kingkhan77/log-sense/internal/models"
	"gorm.io/gorm"
)

type OnCallRepository struct {
	db *gorm.DB
}

func NewOnCallRepository(db *gorm.DB) *OnCallRepository {
	return &OnCallRepository{
		db: db,
	}
}

func (r *OnCallRepository) FindByID(id string) (*models.OnCallSchedule, error) {
	var schedule models.OnCallSchedule
	err := r.db.Where("id = ?", id).First(&schedule).Error
	return &schedule, err
}

func (r *OnCallRepository) Create(
	schedule *models.OnCallSchedule,
) error {

	return r.db.Create(schedule).Error
}

func (r *OnCallRepository) Update(
	schedule *models.OnCallSchedule,
) error {

	return r.db.Save(schedule).Error
}

func (r *OnCallRepository) List(
	tenantID string,
) ([]models.OnCallSchedule, error) {

	var schedules []models.OnCallSchedule

	err := r.db.
		Where("tenant_id = ?", tenantID).
		Order("start_time asc").
		Find(&schedules).
		Error

	return schedules, err
}

func (r *OnCallRepository) Delete(id string) error {
	return r.db.Delete(&models.OnCallSchedule{}, "id = ?", id).Error
}

// HasOverlap returns true if any schedule for the same service overlaps [start, end).
// Pass excludeID = "" when creating; pass the existing ID when updating.
func (r *OnCallRepository) HasOverlap(serviceID, excludeID string, start, end time.Time) (bool, error) {
	var count int64
	q := r.db.Model(&models.OnCallSchedule{}).
		Where("service_id = ? AND start_time < ? AND end_time > ?", serviceID, end, start)
	if excludeID != "" {
		q = q.Where("id != ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

func (r *OnCallRepository) GetCurrentOnCall(
	serviceID string,
) (*models.OnCallSchedule, error) {

	var schedule models.OnCallSchedule

	now := time.Now()

	err := r.db.
		Where(
			"service_id = ? AND start_time <= ? AND end_time >= ?",
			serviceID,
			now,
			now,
		).
		First(&schedule).
		Error

	return &schedule, err
}

func (r *OnCallRepository) GetAllCurrentOnCall(
	serviceID string,
) ([]models.OnCallSchedule, error) {

	var schedules []models.OnCallSchedule

	now := time.Now()

	err := r.db.
		Where(
			"service_id = ? AND start_time <= ? AND end_time >= ?",
			serviceID,
			now,
			now,
		).
		Find(&schedules).
		Error

	return schedules, err
}