package service

import (
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type EmployeeService struct {
	repo *repository.EmployeeRepository
}

func NewEmployeeService(repo *repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

func (s *EmployeeService) GetAllEmployees(filters dto.EmployeeFilter) ([]model.Employee, error) {
	return s.repo.FindAll(filters)
}

func (s *EmployeeService) GetEmployeeOptions(filters dto.EmployeeOptionsFilter) ([]dto.EmployeeOptionResponse, error) {
	return s.repo.FindOptions(filters)
}