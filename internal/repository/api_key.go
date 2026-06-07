package repository

import (
	"github.com/kingkhan77/log-sense/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type APIKeyRepository struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

func (r *APIKeyRepository) Create(key *models.APIKey) error {
	return r.db.Create(key).Error
}

func (r *APIKeyRepository) ListActive() ([]models.APIKey, error) {
	var keys []models.APIKey
	err := r.db.Where("is_active = ?", true).Find(&keys).Error
	return keys, err
}

func (r *APIKeyRepository) ListByTenant(tenantID string) ([]models.APIKey, error) {
	var keys []models.APIKey
	err := r.db.Where("tenant_id = ?", tenantID).Order("created_at desc").Find(&keys).Error
	return keys, err
}

func (r *APIKeyRepository) FindByID(id string) (*models.APIKey, error) {
	var key models.APIKey
	err := r.db.Where("id = ?", id).First(&key).Error
	return &key, err
}

func (r *APIKeyRepository) Revoke(id string) error {
	return r.db.Model(&models.APIKey{}).Where("id = ?", id).Update("is_active", false).Error
}

// ValidateKey finds a candidate by its stored prefix, then runs a single bcrypt
// comparison. This avoids the O(n×bcrypt) full-table scan of the old approach.
func (r *APIKeyRepository) ValidateKey(rawKey string) (*models.APIKey, error) {
	const minLen = len("ls_") + 8
	if len(rawKey) < minLen {
		return nil, gorm.ErrRecordNotFound
	}
	prefix := rawKey[len("ls_") : len("ls_")+8]

	var candidates []models.APIKey
	if err := r.db.Where("key_prefix = ? AND is_active = ?", prefix, true).Find(&candidates).Error; err != nil {
		return nil, err
	}

	for i := range candidates {
		if bcrypt.CompareHashAndPassword([]byte(candidates[i].KeyHash), []byte(rawKey)) == nil {
			return &candidates[i], nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}
