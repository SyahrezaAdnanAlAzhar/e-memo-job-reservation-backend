package repository

import (
	"database/sql"
	"strconv"
	"strings"
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


// GET ALL
func (r *AreaRepository) FindAll(filters map[string]string) ([]Area, error) {
	baseQuery := "SELECT id, department_id, name, is_active, created_at, updated_at FROM area"
	var conditions []string
	var args []interface{}
	argID := 1

	if val, ok := filters["is_active"]; ok {
		conditions = append(conditions, "is_active = $"+strconv.Itoa(argID))
		args = append(args, val)
		argID++
	}
	if val, ok := filters["department_id"]; ok {
		conditions = append(conditions, "department_id = $"+strconv.Itoa(argID))
		args = append(args, val)
		argID++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY department_id ASC, id ASC"

	rows, err := r.DB.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var areas []Area
	for rows.Next() {
		var a Area
		err := rows.Scan(&a.ID, &a.DepartmentID, &a.Name, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		areas = append(areas, a)
	}
	return areas, nil
}