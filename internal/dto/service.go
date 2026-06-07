package dto

type CreateServiceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateServiceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}