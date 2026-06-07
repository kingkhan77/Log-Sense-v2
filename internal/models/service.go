package models

import "time"

type Service struct {
	ID          string `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID    string `gorm:"type:uuid;not null"   json:"tenant_id"`
	Name        string `gorm:"not null"             json:"name"`
	Description string `                            json:"description"`
	IsActive    bool   `gorm:"default:true"         json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
}
