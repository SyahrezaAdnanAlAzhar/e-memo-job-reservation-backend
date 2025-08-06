package model

import "database/sql"

type Employee struct {
	NPK                string        `json:"npk"`
	DepartmentID       int           `json:"department_id"`
	AreaID             sql.NullInt64 `json:"area_id"`
	Name               string        `json:"name"`
	IsActive           bool          `json:"is_active"`
	Position           Position      `json:"position"`
}

type Position struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
