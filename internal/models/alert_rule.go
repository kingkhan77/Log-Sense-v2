package models

import (
	"time"

	"gorm.io/datatypes"
)

type AlertRule struct {
	ID            string         `gorm:"type:uuid;primaryKey"  json:"id"`
	TenantID      string         `gorm:"type:uuid;not null"    json:"tenant_id"`
	ServiceID     string         `gorm:"type:uuid;not null"    json:"service_id"`
	CreatedBy     string         `gorm:"type:uuid;not null"    json:"created_by"`
	Name          string         `gorm:"not null"              json:"name"`
	Description   string         `                             json:"description"`
	Severity      string         `gorm:"not null"              json:"severity"`
	Query         datatypes.JSON `gorm:"type:jsonb"            json:"query"`
	Threshold     int            `gorm:"not null"              json:"threshold"`
	WindowMinutes int            `gorm:"not null"              json:"window_minutes"`
	IsEnabled     bool           `gorm:"default:true"          json:"is_enabled"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
