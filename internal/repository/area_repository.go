package repository

import (
	"database/sql"
	"time"
)

type Area struct {
	ID           int       `json:"id"`
	DepartmentID int       `json:"department_id"`
	Name         string    `json:"name"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateAreaRequest struct {
	DepartmentID int    `json:"department_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
}

type AreaRepository struct {
	DB *sql.DB
}

func NewAreaRepository(db *sql.DB) *AreaRepository {
	return &AreaRepository{DB: db}
}

// CREATE
func (r *AreaRepository) Create(req CreateAreaRequest) (*Area, error) {
	query := `
        INSERT INTO area (department_id, name, is_active) 
        VALUES ($1, $2, false)
        RETURNING id, department_id, name, is_active, created_at, updated_at`

	row := r.DB.QueryRow(query, req.DepartmentID, req.Name)

	var newArea Area
	err := row.Scan(
		&newArea.ID, &newArea.DepartmentID, &newArea.Name,
		&newArea.IsActive, &newArea.CreatedAt, &newArea.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &newArea, nil
}
