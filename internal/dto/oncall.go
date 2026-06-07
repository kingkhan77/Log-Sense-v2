package dto

type CreateScheduleRequest struct {
	ServiceID string `json:"service_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`

	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}