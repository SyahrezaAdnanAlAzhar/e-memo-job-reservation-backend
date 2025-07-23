package service

import (
	"errors"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type SpecifiedLocationService struct {
	repo *repository.SpecifiedLocationRepository
}

func NewSpecifiedLocationService(repo *repository.SpecifiedLocationRepository) *SpecifiedLocationService {
	return &SpecifiedLocationService{repo: repo}
}

// CREATE
func (s *SpecifiedLocationService) CreateSpecifiedLocation(req dto.CreateSpecifiedLocationRequest) (*model.SpecifiedLocation, error) {
	newLoc, err := s.repo.Create(req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" { // foreign_key_violation
				return nil, errors.New("invalid physical_location_id")
			}
			if pgErr.Code == "23505" { // unique_violation
				return nil, errors.New("location name already exists in this physical location")
			}
		}
		return nil, err
	}
	return newLoc, nil
}

// GET ALL
func (s *SpecifiedLocationService) GetAllSpecifiedLocations() ([]model.SpecifiedLocation, error) {
	return s.repo.FindAll()
}

// GET ALL BY PHYSICAL LOCATION ID
func (s *SpecifiedLocationService) GetSpecifiedLocationsByPhysicalLocationID(physicalLocationID int) ([]model.SpecifiedLocation, error) {
	return s.repo.FindByPhysicalLocationID(physicalLocationID)
}

// GET BY ID
func (s *SpecifiedLocationService) GetSpecifiedLocationByID(id int) (*model.SpecifiedLocation, error) {
	return s.repo.FindByID(id)
}