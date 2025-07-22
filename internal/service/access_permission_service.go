package service

import (
	"errors"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type AccessPermissionService struct {
	repo *repository.AccessPermissionRepository
}

func NewAccessPermissionService(repo *repository.AccessPermissionRepository) *AccessPermissionService {
	return &AccessPermissionService{repo: repo}
}

// CREATE
func (s *AccessPermissionService) CreateAccessPermission(req repository.CreateAccessPermissionRequest) (*repository.AccessPermission, error) {
	newPermission, err := s.repo.Create(req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { 
			return nil, errors.New("permission name already exists")
		}
		return nil, err
	}
	return newPermission, nil
}

// GET ALL
func (s *AccessPermissionService) GetAllAccessPermissions() ([]repository.AccessPermission, error) {
	return s.repo.FindAll()
}

// GET BY ID
func (s *AccessPermissionService) GetAccessPermissionByID(id int) (*repository.AccessPermission, error) {
	return s.repo.FindByID(id)
}
