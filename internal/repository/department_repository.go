package repository

import (
	"database/sql"
	"time"
	"errors"
)

type Department struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	ReceiveJob bool      `json:"receive_job"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpdateDepartmentRequest struct {
	Name       string `json:"name" binding:"required"`
	ReceiveJob bool   `json:"receive_job"`
	IsActive   bool   `json:"is_active"`
}

type DepartmentRepository struct {
	DB *sql.DB
}

func NewDepartmentRepository(db *sql.DB) *DepartmentRepository {
	return &DepartmentRepository{DB: db}
}

// HELPER

// CHECK UNIQUE NAME
func (r *DepartmentRepository) IsNameTaken(name string, currentID int) (bool, error) {
	var existsID int
	query := "SELECT id FROM department WHERE name = $1 AND id != $2 LIMIT 1"

	err := r.DB.QueryRow(query, name, currentID).Scan(&existsID)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	
	return true, nil
} 


// MAIN

// GET ALL
func (r *DepartmentRepository) FindAll() ([]Department, error) {
	query := "SELECT id, name, receive_job, is_active, created_at, updated_at FROM department ORDER BY id ASC"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []Department
	for rows.Next() {
		var d Department
		err := rows.Scan(&d.ID, &d.Name, &d.ReceiveJob, &d.IsActive, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}
	return departments, nil
}


// SELECT BY ID
func (r *DepartmentRepository) FindByID(id int) (*Department, error) {
	query := "SELECT id, name, receive_job, is_active, created_at, updated_at FROM department WHERE id = $1"
	row := r.DB.QueryRow(query, id)

	var d Department
	err := row.Scan(&d.ID, &d.Name, &d.ReceiveJob, &d.IsActive, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}


// DELETE
func (r *DepartmentRepository) Delete(id int) error {
	query := "DELETE FROM department WHERE id = $1"
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("department not found or already deleted")
	}

	return nil
}


// UPDATE
func (r *DepartmentRepository) Update(id int, req UpdateDepartmentRequest) (*Department, error) {
	query := `
		UPDATE department 
        SET name = $1, receive_job = $2, is_active = $3, updated_at = NOW() 
        WHERE id = $4 
        RETURNING id, name, receive_job, is_active, created_at, updated_at`

	row := r.DB.QueryRow(query, req.Name, req.ReceiveJob, req.IsActive, id)

	var updatedDept Department
	err := row.Scan(
		&updatedDept.ID, &updatedDept.Name,
		&updatedDept.ReceiveJob, &updatedDept.IsActive,
		&updatedDept.CreatedAt, &updatedDept.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &updatedDept, nil
}