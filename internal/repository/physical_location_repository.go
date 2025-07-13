package repository

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type PhysicalLocation struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePhysicalLocationRequest struct {
	Name string `json:"name" binding:"required"`
}

type PhysicalLocationRepository struct {
	DB *sql.DB
}

func NewPhysicalLocationRepository(db *sql.DB) *PhysicalLocationRepository {
	return &PhysicalLocationRepository{DB: db}
}


// CREATE
func (r *PhysicalLocationRepository) Create(req CreatePhysicalLocationRequest) (*PhysicalLocation, error) {
	query := `
        INSERT INTO physical_location (name, is_active) 
        VALUES ($1, false)
        RETURNING id, name, is_active, created_at, updated_at`

	row := r.DB.QueryRow(query, req.Name)

	var newLoc PhysicalLocation
	err := row.Scan(&newLoc.ID, &newLoc.Name, &newLoc.IsActive, &newLoc.CreatedAt, &newLoc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &newLoc, nil
}


// GET ALL
func (r *PhysicalLocationRepository) FindAll(filters map[string]string) ([]PhysicalLocation, error) {
	baseQuery := "SELECT id, name, is_active, created_at, updated_at FROM physical_location"
	var conditions []string
	var args []interface{}
	argID := 1

	if val, ok := filters["is_active"]; ok {
		conditions = append(conditions, "is_active = $"+strconv.Itoa(argID))
		args = append(args, val)
		argID++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}
	baseQuery += " ORDER BY id ASC"

	rows, err := r.DB.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []PhysicalLocation
	for rows.Next() {
		var loc PhysicalLocation
		err := rows.Scan(&loc.ID, &loc.Name, &loc.IsActive, &loc.CreatedAt, &loc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}


// GET BY ID
func (r *PhysicalLocationRepository) FindByID(id int) (*PhysicalLocation, error) {
	query := "SELECT id, name, is_active, created_at, updated_at FROM physical_location WHERE id = $1"
	row := r.DB.QueryRow(query, id)

	var loc PhysicalLocation
	err := row.Scan(&loc.ID, &loc.Name, &loc.IsActive, &loc.CreatedAt, &loc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &loc, nil
}