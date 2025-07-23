package dto

type CreateWorkflowRequest struct {
	Name            string `json:"name" binding:"required"`
	StatusTicketIDs []int  `json:"status_ticket_ids" binding:"required,min=1"`
}