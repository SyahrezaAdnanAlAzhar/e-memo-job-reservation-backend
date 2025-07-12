package service

import "github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"

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