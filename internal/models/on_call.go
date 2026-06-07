package models

import "time"

type OnCallSchedule struct {
	ID        string    `gorm:"type:uuid;primaryKey"  json:"id"`
	TenantID  string    `gorm:"type:uuid;not null"    json:"tenant_id"`
	ServiceID string    `gorm:"type:uuid;not null"    json:"service_id"`
	UserID    string    `gorm:"type:uuid;not null"    json:"user_id"`
	StartTime time.Time `                             json:"start_time"`
	EndTime   time.Time `                             json:"end_time"`
	CreatedAt time.Time `                             json:"created_at"`
	UpdatedAt time.Time `                             json:"updated_at"`
}

func (OnCallSchedule) TableName() string {
	return "oncall_rotations"
}
