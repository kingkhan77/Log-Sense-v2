package dto

type IngestLogRequest struct {
	Level     string                 `json:"level" binding:"required,oneof=INFO WARN ERROR CRITICAL"`
	Message   string                 `json:"message" binding:"required"`
	Timestamp string                 `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}
