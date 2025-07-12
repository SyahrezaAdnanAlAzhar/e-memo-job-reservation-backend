package repository

import (
	"database/sql"
	"time"
	"errors"
	"strconv"
    "strings"
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

type CreateDepartmentRequest struct {
	Name       string `json:"name" binding:"required"`
	ReceiveJob bool   `json:"receive_job"`
}

type UpdateStatusRequest struct {
	IsActive bool `json:"is_active"`
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

// CREATE
func (r *DepartmentRepository) Create(req CreateDepartmentRequest) (*Department, error) {
	query := `
        INSERT INTO department (name, receive_job, is_active)
        VALUES ($1, $2, false)
        RETURNING id, name, receive_job, is_active, created_at, updated_at`

	row := r.DB.QueryRow(query, req.Name, req.ReceiveJob)
	
	var newDept Department
	err := row.Scan(
		&newDept.ID, &newDept.Name, &newDept.ReceiveJob, 
		&newDept.IsActive, &newDept.CreatedAt, &newDept.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &newDept, nil
}


// GET ALL
func (r *DepartmentRepository) FindAll(filters map[string]string) ([]Department, error) {
	query := "SELECT id, name, receive_job, is_active, created_at, updated_at FROM department"
	
	var conditions []string
	var args []interface{}
	argID := 1

	if val, ok := filters["is_active"]; ok {
		conditions = append(conditions, "is_active = $"+strconv.Itoa(argID))
		args = append(args, val)
		argID++
	}
	if val, ok := filters["receive_job"]; ok {
		conditions = append(conditions, "receive_job = $"+strconv.Itoa(argID))
		args = append(args, val)
		argID++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY id ASC"
	
	rows, err := r.DB.Query(query, args...)
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


// CHANGE ACTIVE STATUS
func (r *DepartmentRepository) UpdateActiveStatus(id int, isActive bool) error {
	query := "UPDATE department SET is_active = $1, updated_at = NOW() WHERE id = $2"
	result, err := r.DB.Exec(query, isActive, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows 
	}
	return nil
}