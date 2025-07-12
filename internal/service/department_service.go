package service

import (
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

type DepartmentService struct {
	repo *repository.DepartmentRepository
}

func NewDepartmentService(repo *repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

// CREATE
func (s *DepartmentService) CreateDepartment(req repository.CreateDepartmentRequest) (*repository.Department, error) {
	if req.Name == "" {
		return nil, errors.New("department name is required")
	}

	newDept, err := s.repo.Create(req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errors.New("department name already exists")
		}
		return nil, err
	}

	return newDept, nil
}


// GET ALL
func (s *DepartmentService) GetAllDepartments(filters map[string]string) ([]repository.Department, error) {
	return s.repo.FindAll(filters)
}


// GET BY ID
func (s *DepartmentService) GetDepartmentByID(id int) (*repository.Department, error) {
	return s.repo.FindByID(id)
}


// DELETE
func (s *DepartmentService) DeleteDepartment(id int) error {
	return s.repo.Delete(id)
}


// UPDATE
func (s *DepartmentService) UpdateDepartment(id int, req repository.UpdateDepartmentRequest) (*repository.Department, error) {
	isTaken, err := s.repo.IsNameTaken(req.Name, id)
	if err != nil {
		
		return nil, err
	}
	if isTaken {
		return nil, errors.New("department name already exists")
	}

	return s.repo.Update(id, req)
}


// CHANGE ACTIVE STATUS
func (s *DepartmentService) UpdateDepartmentActiveStatus(id int, req repository.UpdateStatusRequest) error {
	return s.repo.UpdateActiveStatus(id, req.IsActive)
}