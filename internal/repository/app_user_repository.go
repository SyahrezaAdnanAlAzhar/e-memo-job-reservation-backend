package repository

import (
	"database/sql"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
)

type AppUserRepository struct {
	DB *sql.DB
}

func NewAppUserRepository(db *sql.DB) *AppUserRepository {
	return &AppUserRepository{DB: db}
}

// FindByUsername
func (r *AppUserRepository) FindByUsername(username string) (*model.AppUser, error) {
	query := `
        SELECT 
            id, username, password_hash, user_type, 
            employee_npk, employee_position_id, is_active, 
            created_at, updated_at
        FROM app_user 
        WHERE username = $1 AND is_active = true`

	row := r.DB.QueryRow(query, username)

	var user model.AppUser
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.UserType,
		&user.EmployeeNPK,
		&user.EmployeePositionID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByID
func (r *AppUserRepository) FindByID(id int) (*model.AppUser, error) {
	query := `
        SELECT 
            id, username, password_hash, user_type, 
            employee_npk, employee_position_id, is_active, 
            created_at, updated_at
        FROM app_user 
        WHERE id = $1`

	row := r.DB.QueryRow(query, id)

	var user model.AppUser
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.UserType,
		&user.EmployeeNPK,
		&user.EmployeePositionID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// FindByNPK
func (r *AppUserRepository) FindByNPK(npk string) (*model.AppUser, error) {
	query := `
        SELECT 
            id, username, password_hash, user_type, 
            employee_npk, employee_position_id, is_active, 
            created_at, updated_at
        FROM app_user 
        WHERE employee_npk = $1`

	row := r.DB.QueryRow(query, npk)

	var user model.AppUser
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.UserType,
		&user.EmployeeNPK,
		&user.EmployeePositionID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, err
}
