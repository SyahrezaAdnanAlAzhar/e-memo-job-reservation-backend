package dto

import "time"

type CreateTicketRequest struct {
	DepartmentTargetID  int     `json:"department_target_id" binding:"required,gt=0"`
	PhysicalLocationID  *int    `json:"physical_location_id"`
	SpecifiedLocationID *int    `json:"specified_location_id"`
	Description         string  `json:"description" binding:"required"`
	Deadline            *string `json:"deadline"` // "YYYY-MM-DD"
}

type UpdateTicketRequest struct {
	DepartmentTargetID  int     `json:"department_target_id" binding:"required"`
	Description         string  `json:"description" binding:"required"`
	PhysicalLocationID  *int    `json:"physical_location_id"`
	SpecifiedLocationID *int    `json:"specified_location_id"`
	Deadline            *string `json:"deadline"`
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

type TicketDetailResponse struct {
	// CORE INFORMATION
	TicketID           int        `json:"ticket_id"`
	Description        string     `json:"description"`
	TicketPriority     int        `json:"ticket_priority"`
	Version            int        `json:"version"`
	
	// DEPARTMENT INFORMATION
	DepartmentTargetID int        `json:"department_target_id"`
	DepartmentTargetName string   `json:"department_target_name"`

	// JOB INFOMATION
	JobID                *int       `json:"job_id"`
	JobPriority          *int       `json:"job_priority"`
	
	// LOCATION INFORMATION
	LocationName          *string    `json:"location_name"`
	SpecifiedLocationName *string    `json:"specified_location_name"`

	// TIME INFORMATION
	CreatedAt             time.Time  `json:"created_at"`
	TicketAgeDays         *int       `json:"ticket_age_days"`
	Deadline              *time.Time `json:"deadline"`
	DaysRemaining         *int       `json:"days_remaining"`

	// PEOPLE INFORMATION
	RequestorName         string     `json:"requestor_name"`
	RequestorDepartment   *string    `json:"requestor_department"`
	PicName               *string    `json:"pic_name"`
	PicAreaName           *string    `json:"pic_area_name"`

	// STATUS IFNORMATION
	CurrentStatus         *string    `json:"current_status"`
	CurrentStatusHexCode  *string    `json:"current_status_hex_code"`
	CurrentSectionName    *string    `json:"current_section_name"`
}