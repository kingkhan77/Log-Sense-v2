package repository

import (
	"github.com/kingkhan77/log-sense/internal/models"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(service *models.Service) error {
	return r.db.Create(service).Error
}

func (r *ServiceRepository) FindByID(id string) (*models.Service, error) {
	var service models.Service
	err := r.db.Where("id = ?", id).First(&service).Error
	return &service, err
}

func (r *ServiceRepository) List(tenantID string) ([]models.Service, error) {
	var services []models.Service
	err := r.db.Where("tenant_id = ?", tenantID).Find(&services).Error
	return services, err
}

func (r *ServiceRepository) Update(service *models.Service) error {
	return r.db.Save(service).Error
}

func (r *ServiceRepository) Delete(id string) error {
	return r.db.Delete(&models.Service{}, "id = ?", id).Error
}
