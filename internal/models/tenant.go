package models

import "time"

type Tenant struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Name      string `gorm:"not null"`
	Status    string `gorm:"not null;default:ACTIVE"`
	APIKey    string `gorm:"column:api_key;unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Tenant) TableName() string {
	return "tenants"
}
