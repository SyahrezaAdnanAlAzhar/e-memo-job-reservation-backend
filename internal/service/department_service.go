package service

import (
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"errors"
)

type DepartmentService struct {
	repo *repository.DepartmentRepository
}

func NewDepartmentService(repo *repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) GetAllDepartments() ([]repository.Department, error) {
	return s.repo.FindAll()
}

func (s *DepartmentService) GetDepartmentByID(id int) (*repository.Department, error) {
	return s.repo.FindByID(id)
}

func (s *DepartmentService) DeleteDepartment(id int) error {
	return s.repo.Delete(id)
}

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