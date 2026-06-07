package service

import (
	"time"

	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantService struct {
	repo *repository.ServiceRepository
}

func NewServiceService(repo *repository.ServiceRepository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) Create(tenantID, name, description string) (*models.Service, error) {
	svc := &models.Service{
		ID:          uuid.NewString(),
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.Create(svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *TenantService) List(tenantID string) ([]models.Service, error) {
	return s.repo.List(tenantID)
}

func (s *TenantService) GetByID(tenantID, id string) (*models.Service, error) {
	svc, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if svc.TenantID != tenantID {
		return nil, gorm.ErrRecordNotFound
	}
	return svc, nil
}

func (s *TenantService) Update(tenantID, id, name, description string) error {
	svc, err := s.GetByID(tenantID, id)
	if err != nil {
		return err
	}
	if name != "" {
		svc.Name = name
	}
	if description != "" {
		svc.Description = description
	}
	return s.repo.Update(svc)
}

func (s *TenantService) Delete(tenantID, id string) error {
	if _, err := s.GetByID(tenantID, id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
