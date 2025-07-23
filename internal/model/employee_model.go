package model

import "database/sql"

type Employee struct {
	NPK                string        `json:"npk"`
	DepartmentID       sql.NullInt64 `json:"department_id"`
	AreaID             sql.NullInt64 `json:"area_id"`
	Name               string        `json:"name"`
	IsActive           bool          `json:"is_active"`
	EmployeePositionID int           `json:"employee_position_id"`
}
