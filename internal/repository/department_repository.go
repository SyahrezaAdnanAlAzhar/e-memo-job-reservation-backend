package repository

import (
	"database/sql"
	"time"
)

type Department struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	ReceiveJob bool      `json:"receive_job"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DepartmentRepository struct {
	DB *sql.DB
}

func NewDepartmentRepository(db *sql.DB) *DepartmentRepository {
	return &DepartmentRepository{DB: db}
}


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