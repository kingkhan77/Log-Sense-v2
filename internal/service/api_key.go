package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/kingkhan77/log-sense/internal/models"
	"github.com/kingkhan77/log-sense/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type APIKeyService struct {
	repo        *repository.APIKeyRepository
	serviceRepo *repository.ServiceRepository
}

func NewAPIKeyService(
	repo *repository.APIKeyRepository,
	serviceRepo *repository.ServiceRepository,
) *APIKeyService {
	return &APIKeyService{repo: repo, serviceRepo: serviceRepo}
}

// keyPrefixLen is the number of base64 chars (after the "ls_" marker) stored
// in plain text for O(1) candidate lookup before bcrypt comparison.
const keyPrefixLen = 8

func extractPrefix(raw string) string {
	const marker = "ls_"
	return raw[len(marker) : len(marker)+keyPrefixLen]
}

func (s *APIKeyService) Create(tenantID, serviceID, name string) (string, error) {
	svc, err := s.serviceRepo.FindByID(serviceID)
	if err != nil || svc.TenantID != tenantID {
		return "", gorm.ErrRecordNotFound
	}

	raw, err := generateAPIKey()
	if err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	key := &models.APIKey{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		ServiceID: serviceID,
		KeyHash:   string(hash),
		KeyPrefix: extractPrefix(raw),
		Name:      name,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(key); err != nil {
		return "", err
	}

	return raw, nil
}

func (s *APIKeyService) List(tenantID string) ([]models.APIKey, error) {
	return s.repo.ListByTenant(tenantID)
}

func (s *APIKeyService) Revoke(tenantID, id string) error {
	key, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if key.TenantID != tenantID {
		return gorm.ErrRecordNotFound
	}
	return s.repo.Revoke(id)
}

func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("ls_%s", base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)), nil
}
