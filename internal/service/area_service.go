package service

import (	
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

type AreaService struct {
	repo *repository.AreaRepository
}

func NewAreaService(repo *repository.AreaRepository) *AreaService {
	return &AreaService{repo: repo}
}

// CREATE
func (s *AreaService) CreateArea(req repository.CreateAreaRequest) (*repository.Area, error) {
	if req.Name == "" {
		return nil, errors.New("area name is required")
	}
	if req.DepartmentID <= 0 {
		return nil, errors.New("valid department_id is required")
	}

	newArea, err := s.repo.Create(req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { 
			return nil, errors.New("area name already exists in this department")
		}
		return nil, err
	}

	return newArea, nil
}


// GET ALL
func (s *AreaService) GetAllAreas(filters map[string]string) ([]repository.Area, error) {
	return s.repo.FindAll(filters)
}