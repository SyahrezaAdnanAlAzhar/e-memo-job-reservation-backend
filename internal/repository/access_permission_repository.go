package repository

import (
	"database/sql"
	"time"
)

type AccessPermission struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAccessPermissionRequest struct {
	Name string `json:"name" binding:"required"`
}

type AccessPermissionRepository struct {
	DB *sql.DB
}

func NewAccessPermissionRepository(db *sql.DB) *AccessPermissionRepository {
	return &AccessPermissionRepository{DB: db}
}

// CREATE
func (r *AccessPermissionRepository) Create(req CreateAccessPermissionRequest) (*AccessPermission, error) {
	query := `
        INSERT INTO access_permission (name, is_active) 
        VALUES ($1, false)
        RETURNING id, name, is_active, created_at, updated_at`

	row := r.DB.QueryRow(query, req.Name)

	var newPermission AccessPermission
	err := row.Scan(
		&newPermission.ID, &newPermission.Name, &newPermission.IsActive,
		&newPermission.CreatedAt, &newPermission.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &newPermission, nil
}

// GET ALL
func (r *AccessPermissionRepository) FindAll() ([]AccessPermission, error) {
	query := "SELECT id, name, is_active, created_at, updated_at FROM access_permission ORDER BY id ASC"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []AccessPermission
	for rows.Next() {
		var p AccessPermission
		err := rows.Scan(&p.ID, &p.Name, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}

// GET BY ID
func (r *AccessPermissionRepository) FindByID(id int) (*AccessPermission, error) {
	query := "SELECT id, name, is_active, created_at, updated_at FROM access_permission WHERE id = $1"
	row := r.DB.QueryRow(query, id)

	var p AccessPermission
	err := row.Scan(&p.ID, &p.Name, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}