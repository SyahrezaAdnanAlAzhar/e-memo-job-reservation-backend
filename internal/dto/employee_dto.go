package dto

type EmployeeFilter struct {
	DepartmentID       int    `form:"department_id"`
	AreaID             int    `form:"area_id"`
	EmployeePositionID int    `form:"employee_position_id"`
	Name               string `form:"name"`
	NPK                string `form:"npk"`
	IsActive           *bool  `form:"is_active"`
}