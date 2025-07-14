package repository

import (
	"database/sql"
	"errors"
)

type Employee struct {
	NPK          string        `json:"npk"`
	DepartmentID sql.NullInt64 `json:"department_id"`
	AreaID       sql.NullInt64 `json:"area_id"`
	Name         string        `json:"name"`
	IsActive     bool          `json:"is_active"`
	PositionID   int           `json:"position_id"`
}

type EmployeeRepository struct {
	DB *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

func (r *EmployeeRepository) GetAllEmployees() ([]Employee, error) {
	rows, err := r.DB.Query("SELECT npk, name, position_id FROM employee WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.NPK, &e.Name, &e.PositionID); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func (r *EmployeeRepository) FindByNPK(npk string) (*Employee, error) {
	query := "SELECT npk, name, position_id, is_active FROM employee WHERE npk = $1"
	row := r.DB.QueryRow(query, npk)

	var e Employee
	err := row.Scan(&e.NPK, &e.Name, &e.PositionID, &e.IsActive)
	if err != nil {
		return nil, err
	}
	if !e.IsActive {
		return nil, errors.New("user is not active")
	}
	return &e, nil
}
