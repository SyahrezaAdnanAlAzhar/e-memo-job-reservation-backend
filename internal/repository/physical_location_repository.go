package repository

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
	"errors"
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

type UpdatePhysicalLocationRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdatePhysicalLocationStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type PhysicalLocationRepository struct {
	DB *sql.DB
}

func NewPhysicalLocationRepository(db *sql.DB) *PhysicalLocationRepository {
	return &PhysicalLocationRepository{DB: db}
}

// HELPER

// CHECK UNIQUE NAME
func (r *PhysicalLocationRepository) IsNameTaken(name string, currentID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM physical_location WHERE name = $1 AND id != $2)"
	err := r.DB.QueryRow(query, name, currentID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}



// MAIN

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


// UPDATE
func (r *PhysicalLocationRepository) Update(id int, req UpdatePhysicalLocationRequest) (*PhysicalLocation, error) {
	query := `
        UPDATE physical_location 
        SET name = $1, updated_at = NOW()
        WHERE id = $2
        RETURNING id, name, is_active, created_at, updated_at`

	row := r.DB.QueryRow(query, req.Name, id)

	var updatedLoc PhysicalLocation
	err := row.Scan(&updatedLoc.ID, &updatedLoc.Name, &updatedLoc.IsActive, &updatedLoc.CreatedAt, &updatedLoc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &updatedLoc, nil
}

// DELETE
func (r *PhysicalLocationRepository) Delete(id int) error {
	query := "DELETE FROM physical_location WHERE id = $1"
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("department not found or already deleted")
	}
	return nil
}

// CHANGE STATUS
func (r *PhysicalLocationRepository) UpdateActiveStatus(id int, isActive bool) error {
	query := "UPDATE physical_location SET is_active = $1, updated_at = NOW() WHERE id = $2"
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