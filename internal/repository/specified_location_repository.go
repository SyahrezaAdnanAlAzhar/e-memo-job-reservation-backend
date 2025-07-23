package repository

import (
	"database/sql"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
)

type SpecifiedLocationRepository struct {
	DB *sql.DB
}

func NewSpecifiedLocationRepository(db *sql.DB) *SpecifiedLocationRepository {
	return &SpecifiedLocationRepository{DB: db}
}

// CREATE
func (r *SpecifiedLocationRepository) Create(req dto.CreateSpecifiedLocationRequest) (*model.SpecifiedLocation, error) {
	query := `
        INSERT INTO specified_location (physical_location_id, name, is_active) 
        VALUES ($1, $2, true)
        RETURNING id, physical_location_id, name, is_active, created_at, updated_at`

	row := r.DB.QueryRow(query, req.PhysicalLocationID, req.Name)

	var newLoc model.SpecifiedLocation
	err := row.Scan(
		&newLoc.ID, &newLoc.PhysicalLocationID, &newLoc.Name, &newLoc.IsActive,
		&newLoc.CreatedAt, &newLoc.UpdatedAt,
	)
	return &newLoc, err
}

// GET ALL
func (r *SpecifiedLocationRepository) FindAll() ([]model.SpecifiedLocation, error) {
	query := "SELECT id, physical_location_id, name, is_active, created_at, updated_at FROM specified_location ORDER BY physical_location_id, id ASC"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []model.SpecifiedLocation
	for rows.Next() {
		var loc model.SpecifiedLocation
		err := rows.Scan(
			&loc.ID, &loc.PhysicalLocationID, &loc.Name, &loc.IsActive,
			&loc.CreatedAt, &loc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

// GET ALL BY PHYSICAL LOCATION ID
func (r *SpecifiedLocationRepository) FindByPhysicalLocationID(physicalLocationID int) ([]model.SpecifiedLocation, error) {
	query := "SELECT id, physical_location_id, name, is_active, created_at, updated_at FROM specified_location WHERE physical_location_id = $1 ORDER BY id ASC"
	rows, err := r.DB.Query(query, physicalLocationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []model.SpecifiedLocation
	for rows.Next() {
		var loc model.SpecifiedLocation
		err := rows.Scan(
			&loc.ID, &loc.PhysicalLocationID, &loc.Name, &loc.IsActive,
			&loc.CreatedAt, &loc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

// GET BY ID
func (r *SpecifiedLocationRepository) FindByID(id int) (*model.SpecifiedLocation, error) {
	query := "SELECT id, physical_location_id, name, is_active, created_at, updated_at FROM specified_location WHERE id = $1"
	row := r.DB.QueryRow(query, id)
	var loc model.SpecifiedLocation
	err := row.Scan(
		&loc.ID, &loc.PhysicalLocationID, &loc.Name, &loc.IsActive,
		&loc.CreatedAt, &loc.UpdatedAt,
	)
	return &loc, err
}