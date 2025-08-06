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
	query := `
        SELECT 
            e.npk, e.department_id, e.area_id, e.name, e.is_active,
            ep.id as position_id, ep.name as position_name
        FROM employee e
        JOIN employee_position ep ON e.employee_position_id = ep.id
        WHERE e.is_active = true
        ORDER BY e.name ASC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []model.Employee
	for rows.Next() {
		var e model.Employee
		err := rows.Scan(
			&e.NPK,
			&e.DepartmentID,
			&e.AreaID,
			&e.Name,
			&e.IsActive,
			&e.Position.ID,
			&e.Position.Name,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func (r *EmployeeRepository) FindByNPK(npk string) (*model.Employee, error) {
	query := `
        SELECT 
            e.npk, e.department_id, e.area_id, e.name, e.is_active,
            ep.id as position_id, ep.name as position_name
        FROM employee e
        JOIN employee_position ep ON e.employee_position_id = ep.id
        WHERE e.npk = $1`
	row := r.DB.QueryRow(query, npk)

	var e model.Employee
	err := row.Scan(
		&e.NPK,
		&e.DepartmentID,
		&e.AreaID,
		&e.Name,
		&e.IsActive,
		&e.Position.ID,
		&e.Position.Name,
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
