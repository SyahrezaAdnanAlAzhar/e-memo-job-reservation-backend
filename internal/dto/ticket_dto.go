package dto

type CreateTicketRequest struct {
	DepartmentTargetID  int    `json:"department_target_id" binding:"required,gt=0"`
	PhysicalLocationID  *int   `json:"physical_location_id"`
	SpecifiedLocationID *int   `json:"specified_location_id"`
	Description         string `json:"description" binding:"required"`
}

type UpdateTicketRequest struct {
	DepartmentTargetID  int    `json:"department_target_id" binding:"required"`
	Description         string `json:"description" binding:"required"`
	PhysicalLocationID  *int   `json:"physical_location_id"`
	SpecifiedLocationID *int   `json:"specified_location_id"`
}

type ReorderTicketsRequest struct {
	DepartmentTargetID int                 `json:"department_target_id" binding:"required"`
	Items              []ReorderTicketItem `json:"items" binding:"required,min=1"`
}

type ChangeTicketStatusRequest struct {
	TargetStatusID int `json:"target_status_id" binding:"required"`
}

type RejectTicketRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type ExecuteActionRequest struct {
	ActionName string `json:"action_name" binding:"required"`
	Reason     string `json:"reason"`
}

type ReorderTicketItem struct {
	TicketID int `json:"ticket_id" binding:"required"`
	Version  int `json:"version" binding:"required"`
}
