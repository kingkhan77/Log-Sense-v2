package dto

import "encoding/json"

type CreateRuleRequest struct {
	ServiceID     string          `json:"service_id" binding:"required"`
	Name          string          `json:"name" binding:"required"`
	Description   string          `json:"description"`
	Severity      string          `json:"severity" binding:"required"`
	Query         json.RawMessage `json:"query" binding:"required"`
	Threshold     int             `json:"threshold" binding:"required"`
	WindowMinutes int             `json:"window_minutes" binding:"required"`
	IsEnabled     *bool           `json:"is_enabled"`
}

type UpdateRuleRequest struct {
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Severity      string          `json:"severity"`
	Query         json.RawMessage `json:"query"`
	Threshold     int             `json:"threshold"`
	WindowMinutes int             `json:"window_minutes"`
	IsEnabled     *bool           `json:"is_enabled"`
}
