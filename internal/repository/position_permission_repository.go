package repository

import (
	"database/sql"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
)

type PositionPermissionRepository struct {
	DB *sql.DB
}

func NewPositionPermissionRepository(db *sql.DB) *PositionPermissionRepository {
	return &PositionPermissionRepository{DB: db}
}

// CREATE
func (r *PositionPermissionRepository) Create(req dto.CreatePositionPermissionRequest) (*model.PositionPermission, error) {
	query := `
        INSERT INTO position_permission (employee_position_id, access_permission_id, is_active) 
        VALUES ($1, $2, true)
        RETURNING employee_position_id, access_permission_id, is_active, created_at, updated_at`

	row := r.DB.QueryRow(query, req.EmployeePositionID, req.AccessPermissionID)

	var newPerm model.PositionPermission
	err := row.Scan(
		&newPerm.EmployeePositionID, &newPerm.AccessPermissionID, &newPerm.IsActive,
		&newPerm.CreatedAt, &newPerm.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &newPerm, nil
}

// GET ALL
func (r *PositionPermissionRepository) FindAll() ([]model.PositionPermission, error) {
	query := "SELECT employee_position_id, access_permission_id, is_active, created_at, updated_at FROM position_permission ORDER BY employee_position_id, access_permission_id ASC"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []model.PositionPermission
	for rows.Next() {
		var p model.PositionPermission
		err := rows.Scan(
			&p.EmployeePositionID, &p.AccessPermissionID, &p.IsActive,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}

// CHANGE STATUS
func (r *PositionPermissionRepository) UpdateActiveStatus(posID, permID int, isActive bool) error {
	query := `
        UPDATE position_permission 
        SET is_active = $1, updated_at = NOW() 
        WHERE employee_position_id = $2 AND access_permission_id = $3`
	result, err := r.DB.Exec(query, isActive, posID, permID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DELETE
func (r *PositionPermissionRepository) Delete(posID, permID int) error {
	query := "DELETE FROM position_permission WHERE employee_position_id = $1 AND access_permission_id = $2"
	result, err := r.DB.Exec(query, posID, permID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}