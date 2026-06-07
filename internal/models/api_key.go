package models

import "time"

type APIKey struct {
	ID        string `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID  string `gorm:"type:uuid;not null"   json:"tenant_id"`
	ServiceID string `gorm:"type:uuid;not null"   json:"service_id"`
	KeyHash   string `gorm:"not null;unique"      json:"-"`
	KeyPrefix string `gorm:"size:20;index"        json:"-"`
	Name      string `                            json:"name"`
	IsActive  bool   `gorm:"default:true"         json:"is_active"`
	CreatedAt time.Time `                         json:"created_at"`
}

func (APIKey) TableName() string {
	return "api_keys"
}
