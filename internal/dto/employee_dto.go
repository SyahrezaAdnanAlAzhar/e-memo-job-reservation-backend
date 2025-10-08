package dto

import "github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"

type EmployeeFilter struct {
	DepartmentID       int    `form:"department_id"`
	AreaID             int    `form:"area_id"`
	EmployeePositionID int    `form:"employee_position_id"`
	Name               string `form:"name"`
	NPK                string `form:"npk"`
	IsActive           *bool  `form:"is_active"`
	// PAGINATION
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type EmployeeOptionsFilter struct {
	Role               string `form:"role"`
	SectionID          int    `form:"section_id"`
	DepartmentTargetID int    `form:"department_target_id"`
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	PageSize    int   `json:"page_size"`
}

type PaginatedEmployeeResponse struct {
	Data       []model.Employee `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

type EmployeeOptionResponse struct {
	NPK  string `json:"npk"`
	Name string `json:"name"`
}

type CreateEmployeeRequest struct {
	NPK                string `json:"npk" binding:"required"`
	Name               string `json:"name" binding:"required"`
	DepartmentID       int    `json:"department_id" binding:"required,gt=0"`
	AreaID             *int   `json:"area_id"`
	EmployeePositionID int    `json:"employee_position_id" binding:"required,gt=0"`
}

type UpdateEmployeeRequest struct {
	Name               string `json:"name" binding:"required"`
	DepartmentID       int    `json:"department_id" binding:"required,gt=0"`
	AreaID             *int   `json:"area_id"`
	EmployeePositionID int    `json:"employee_position_id" binding:"required,gt=0"`
}

type UpdateEmployeeStatusRequest struct {
	IsActive bool `json:"is_active"`
}
