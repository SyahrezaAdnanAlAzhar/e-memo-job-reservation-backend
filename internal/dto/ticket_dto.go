package dto

import "time"

type CreateTicketRequest struct {
	DepartmentTargetID  int     `form:"department_target_id" binding:"required,gt=0"`
	PhysicalLocationID  *int    `form:"physical_location_id"`
	SpecifiedLocationID *int    `form:"specified_location_id"`
	Description         string  `form:"description" binding:"required"`
	Deadline            *string `form:"deadline"` // "YYYY-MM-DD"
}

type UpdateTicketRequest struct {
	DepartmentTargetID  int     `json:"department_target_id" binding:"required"`
	Description         string  `json:"description" binding:"required"`
	PhysicalLocationID  *int    `json:"physical_location_id"`
	SpecifiedLocationID *int    `json:"specified_location_id"`
	Deadline            *string `json:"deadline"`
	Version             int     `json:"version" binding:"required,gte=1"`
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
	TicketID       int    `json:"ticket_id"`
	Description    string `json:"description"`
	TicketPriority int    `json:"ticket_priority"`
	Version        int    `json:"version"`

	// DEPARTMENT INFORMATION
	DepartmentTargetID   int    `json:"department_target_id"`
	DepartmentTargetName string `json:"department_target_name"`

	// JOB INFOMATION
	JobID       *int `json:"job_id"`
	JobPriority *int `json:"job_priority"`

	// LOCATION INFORMATION
	LocationName          *string `json:"location_name"`
	SpecifiedLocationName *string `json:"specified_location_name"`

	// TIME INFORMATION
	CreatedAt     time.Time  `json:"created_at"`
	TicketAgeDays *int       `json:"ticket_age_days"`
	Deadline      *time.Time `json:"deadline"`
	DaysRemaining *int       `json:"days_remaining"`

	// PEOPLE INFORMATION
	RequestorName       string  `json:"requestor_name"`
	RequestorDepartment *string `json:"requestor_department"`
	PicName             *string `json:"pic_name"`
	PicAreaName         *string `json:"pic_area_name"`

	// STATUS IFNORMATION
	CurrentStatus        *string `json:"current_status"`
	CurrentStatusHexCode *string `json:"current_status_hex_code"`
	CurrentSectionName   *string `json:"current_section_name"`
}

type TicketFilter struct {
	// FILTER BY ID
	SectionID          int    `form:"section_id"`
	StatusID           int    `form:"status_id"`
	DepartmentTargetID int    `form:"department_target_id"`
	Requestor          string `form:"requestor"`
	PicNPK             string `form:"pic_npk"`

	// FILTER BY SEARCH QUERY
	SearchQuery string `form:"search"`

	// SORTING OPTION
	SortBy string `form:"sort_by"`
}

type DeleteFilesRequest struct {
	FilePathsToDelete []string `json:"file_paths_to_delete" binding:"required,min=1"`
}

type TicketSummaryFilter struct {
	DepartmentID int `form:"department_id"`
	SectionID    int `form:"section_id"`
	Year         int `form:"year"`
	Month        int `form:"month"`
}

type TicketSummaryResponse struct {
	StatusID   int    `json:"status_id"`
	StatusName string `json:"status_name"`
	HexCode    string `json:"hex_code"`
	Total      int64  `json:"total"`
}

type OldestTicketResponse struct {
	TicketID  int       `json:"ticket_id"`
	CreatedAt time.Time `json:"created_at"`
}