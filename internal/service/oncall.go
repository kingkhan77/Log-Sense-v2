package service

import (
	"errors"
	"time"

	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OnCallService struct {
	repo        *repository.OnCallRepository
	serviceRepo *repository.ServiceRepository
}

func NewOnCallService(
	repo *repository.OnCallRepository,
	serviceRepo *repository.ServiceRepository,
) *OnCallService {
	return &OnCallService{repo: repo, serviceRepo: serviceRepo}
}

func (s *OnCallService) Create(schedule *models.OnCallSchedule) error {
	if schedule.EndTime.Before(schedule.StartTime) || schedule.EndTime.Equal(schedule.StartTime) {
		return errors.New("end_time must be after start_time")
	}
	svc, err := s.serviceRepo.FindByID(schedule.ServiceID)
	if err != nil || svc.TenantID != schedule.TenantID {
		return errors.New("service not found for tenant")
	}
	overlap, err := s.repo.HasOverlap(schedule.ServiceID, "", schedule.StartTime, schedule.EndTime)
	if err != nil {
		return err
	}
	if overlap {
		return errors.New("schedule overlaps with an existing on-call period for this service")
	}
	schedule.ID = uuid.NewString()
	now := time.Now()
	schedule.CreatedAt = now
	schedule.UpdatedAt = now
	return s.repo.Create(schedule)
}

func (s *OnCallService) Update(tenantID string, schedule *models.OnCallSchedule) error {
	existing, err := s.repo.FindByID(schedule.ID)
	if err != nil {
		return err
	}
	if existing.TenantID != tenantID {
		return gorm.ErrRecordNotFound
	}
	if schedule.EndTime.Before(schedule.StartTime) {
		return errors.New("end_time must be after start_time")
	}
	overlap, err := s.repo.HasOverlap(schedule.ServiceID, schedule.ID, schedule.StartTime, schedule.EndTime)
	if err != nil {
		return err
	}
	if overlap {
		return errors.New("schedule overlaps with an existing on-call period for this service")
	}
	schedule.TenantID = existing.TenantID
	schedule.CreatedAt = existing.CreatedAt
	schedule.UpdatedAt = time.Now()
	return s.repo.Update(schedule)
}

func (s *OnCallService) GetByID(tenantID, id string) (*models.OnCallSchedule, error) {
	schedule, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if schedule.TenantID != tenantID {
		return nil, gorm.ErrRecordNotFound
	}
	return schedule, nil
}

func (s *OnCallService) List(tenantID string) ([]models.OnCallSchedule, error) {
	return s.repo.List(tenantID)
}

func (s *OnCallService) Delete(tenantID, id string) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing.TenantID != tenantID {
		return gorm.ErrRecordNotFound
	}
	return s.repo.Delete(id)
}

func (s *OnCallService) GetCurrentOnCall(tenantID, serviceID string) (*models.OnCallSchedule, error) {
	svc, err := s.serviceRepo.FindByID(serviceID)
	if err != nil || svc.TenantID != tenantID {
		return nil, gorm.ErrRecordNotFound
	}
	return s.repo.GetCurrentOnCall(serviceID)
}
