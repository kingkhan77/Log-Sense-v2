package models

import "time"

type Alert struct {
	ID           string `gorm:"type:uuid;primaryKey"       json:"id"`
	TenantID     string `gorm:"type:uuid;not null"         json:"tenant_id"`
	ServiceID    string `gorm:"type:uuid;not null"         json:"service_id"`
	RuleID       string `gorm:"type:uuid;not null"         json:"rule_id"`
	Title        string `gorm:"not null"                   json:"title"`
	Description  string `                                  json:"description"`
	Severity     string `gorm:"not null"                   json:"severity"`
	Status       string `gorm:"not null"                   json:"status"`
	Threshold    int    `                                  json:"threshold"`
	CurrentCount int    `                                  json:"current_count"`

	TriggeredAt time.Time `json:"triggered_at"`

	AcknowledgedAt *time.Time `json:"acknowledged_at"`
	AcknowledgedBy *string    `json:"acknowledged_by"`

	ResolvedAt *time.Time `json:"resolved_at"`
	ResolvedBy *string    `json:"resolved_by"`

	NotificationSentAt *time.Time `json:"notification_sent_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
