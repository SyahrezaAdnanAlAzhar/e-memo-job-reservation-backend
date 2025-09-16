package dto

type EmployeeFilter struct {
	DepartmentID       int    `form:"department_id"`
	AreaID             int    `form:"area_id"`
	EmployeePositionID int    `form:"employee_position_id"`
	Name               string `form:"name"`
	NPK                string `form:"npk"`
	IsActive           *bool  `form:"is_active"`
}

type EmployeeOptionsFilter struct {
	Role               string `form:"role"`
	SectionID          int    `form:"section_id"`
	DepartmentTargetID int    `form:"department_target_id"`
}

type EmployeeOptionResponse struct {
	NPK  string `json:"npk"`
	Name string `json:"name"`
}
