package repository

import (
	"context"
	"database/sql"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
)

type EmployeePositionRepository struct {
	DB *sql.DB
}

func NewEmployeePositionRepository(db *sql.DB) *EmployeePositionRepository {
	return &EmployeePositionRepository{DB: db}
}

// CREATE
func (r *EmployeePositionRepository) Create(ctx context.Context, tx *sql.Tx, req dto.CreateEmployeePositionRequest) (*model.EmployeePosition, error) {
	query := `
        INSERT INTO employee_position (name, is_active) 
        VALUES ($1, true)
        RETURNING id, name, is_active, created_at, updated_at`

	row := tx.QueryRowContext(ctx, query, req.Name)

	var newPos model.EmployeePosition
	err := row.Scan(&newPos.ID, &newPos.Name, &newPos.IsActive, &newPos.CreatedAt, &newPos.UpdatedAt)
	return &newPos, err
}

// GET ALL
func (r *EmployeePositionRepository) FindAll() ([]model.EmployeePosition, error) {
	query := "SELECT id, name, is_active, created_at, updated_at FROM employee_position ORDER BY id ASC"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []model.EmployeePosition
	for rows.Next() {
		var p model.EmployeePosition
		err := rows.Scan(&p.ID, &p.Name, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, nil
}

// GET BY ID
func (r *EmployeePositionRepository) FindByID(id int) (*model.EmployeePosition, error) {
	query := "SELECT id, name, is_active, created_at, updated_at FROM employee_position WHERE id = $1"
	row := r.DB.QueryRow(query, id)

	var p model.EmployeePosition
	err := row.Scan(&p.ID, &p.Name, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}