package dto

import "time"

type AssignPICRequest struct {
	PicJobNPK string `json:"pic_job_npk" binding:"required"`
}

type ReorderJobsRequest struct {
	DepartmentTargetID int              `json:"department_target_id" binding:"required"`
	Items              []ReorderJobItem `json:"items" binding:"required,min=1"`
}

type ReorderJobItem struct {
	JobID   int `json:"job_id" binding:"required"`
	Version int `json:"version" binding:"required"`
}

type JobDetailResponse struct {
	// CORE INFORMATION
	JobID          int    `json:"job_id"`
	TicketID       int    `json:"ticket_id"`
	Description    string `json:"description"`
	JobPriority    int    `json:"job_priority"`
	TicketPriority int    `json:"ticket_priority"`
	Version        int    `json:"version"`

	// DEPARTMENT INFORMATION
	AssignedDepartmentID   int    `json:"assigned_department_id"`
	AssignedDepartmentName string `json:"assigned_department_name"`

	// STATUS INFORMATION
	CurrentStatus        *string `json:"current_status"`
	CurrentStatusHexCode *string `json:"current_status_hex_code"`

	// PEOPLE INFORMATION
	PicName       *string `json:"pic_name"`
	RequestorName string  `json:"requestor_name"`

	// TIME INFORMATION
	TicketAgeDays *int       `json:"ticket_age_days"`
	Deadline      *time.Time `json:"deadline"`
	DaysRemaining *int       `json:"days_remaining"`
}

type JobFilter struct {
	AssignedDepartmentID int    `form:"assigned_department_id"`
	StatusID             int    `form:"status_id"`
	PicNPK               string `form:"pic_npk"`
	SearchQuery          string `form:"search"`
	SortBy               string `form:"sort_by"`
}

type AvailableActionResponse struct {
	Name        string `json:"name"`
}