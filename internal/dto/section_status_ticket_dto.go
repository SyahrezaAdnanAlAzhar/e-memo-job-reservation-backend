package dto

type CreateSectionStatusTicketRequest struct {
	Name     string `json:"name" binding:"required"`
	Sequence int    `json:"sequence" binding:"required"`
}

type UpdateSectionStatusTicketStatusRequest struct {
	IsActive bool `json:"is_active"`
}