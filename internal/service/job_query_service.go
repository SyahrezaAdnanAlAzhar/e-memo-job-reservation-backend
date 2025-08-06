package service

import (
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type JobQueryService struct {
	repo *repository.JobQueryRepository
}

func NewJobQueryService(repo *repository.JobQueryRepository) *JobQueryService {
	return &JobQueryService{repo: repo}
}

// GET ALL
func (s *JobQueryService) GetAllJobs(filters dto.JobFilter) ([]dto.JobDetailResponse, error) {
	return s.repo.FindAll(filters)
}

// GET BY ID
func (s *JobQueryService) GetJobByID(id int) (*dto.JobDetailResponse, error) {
	return s.repo.FindByID(id)
}