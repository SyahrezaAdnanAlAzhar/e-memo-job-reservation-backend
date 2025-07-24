package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
)

type EmployeeRepository struct {
	DB *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

func (r *EmployeeRepository) GetAllEmployees() ([]model.Employee, error) {
	rows, err := r.DB.Query("SELECT npk, name, employee_position_id FROM employee WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []model.Employee
	for rows.Next() {
		var e model.Employee
		if err := rows.Scan(&e.NPK, &e.Name, &e.EmployeePositionID); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func (r *EmployeeRepository) FindByNPK(npk string) (*model.Employee, error) {
	query := "SELECT npk, name, employee_position_id, is_active FROM employee WHERE npk = $1"
	row := r.DB.QueryRow(query, npk)

	var e model.Employee
	err := row.Scan(
		&e.NPK, &e.DepartmentID, &e.AreaID, &e.Name, &e.IsActive, &e.EmployeePositionID,
		&e.Position.ID, &e.Position.Name,
	)
	if err != nil {
		return nil, err
	}
	if !e.IsActive {
		return nil, errors.New("user is not active")
	}
	return &e, nil
}

func (r *EmployeeRepository) GetEmployeePositionID(ctx context.Context, npk string) (int, error) {
	var positionID int
	query := "SELECT employee_position_id FROM employee WHERE npk = $1"
	err := r.DB.QueryRowContext(ctx, query, npk).Scan(&positionID)
	return positionID, err
}
