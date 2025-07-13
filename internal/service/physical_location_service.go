package service

import (
	"errors"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type PhysicalLocationService struct {
	repo *repository.PhysicalLocationRepository
}

func NewPhysicalLocationService(repo *repository.PhysicalLocationRepository) *PhysicalLocationService {
	return &PhysicalLocationService{repo: repo}
}


// CREATE
func (s *PhysicalLocationService) CreatePhysicalLocation(req repository.CreatePhysicalLocationRequest) (*repository.PhysicalLocation, error) {
	newLoc, err := s.repo.Create(req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { 
			return nil, errors.New("physical location name already exists")
		}
		return nil, err
	}
	return newLoc, nil
}


// GET ALL
func (s *PhysicalLocationService) GetAllPhysicalLocations(filters map[string]string) ([]repository.PhysicalLocation, error) {
	return s.repo.FindAll(filters)
}


// GET BY ID
func (s *PhysicalLocationService) GetPhysicalLocationByID(id int) (*repository.PhysicalLocation, error) {
	return s.repo.FindByID(id)
}