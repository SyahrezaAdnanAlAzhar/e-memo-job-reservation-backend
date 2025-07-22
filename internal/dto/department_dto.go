package dto

type CreateDepartmentRequest struct {
	Name       string `json:"name" binding:"required"`
	ReceiveJob bool   `json:"receive_job"`
}

type UpdateDepartmentRequest struct {
	Name       string `json:"name" binding:"required"`
	ReceiveJob bool   `json:"receive_job"`
	IsActive   bool   `json:"is_active"`
}

type UpdateStatusRequest struct {
	IsActive bool `json:"is_active"`
}