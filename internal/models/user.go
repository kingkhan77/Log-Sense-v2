package models

import "time"

type User struct {
	ID           string `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID     string `gorm:"type:uuid;not null"   json:"tenant_id"`
	Name         string `gorm:"not null"             json:"name"`
	Email        string `gorm:"not null"             json:"email"`
	PasswordHash string `gorm:"column:password_hash;not null" json:"-"`
	Role         string `gorm:"not null"             json:"role"`
	IsActive     bool   `gorm:"default:true"         json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
