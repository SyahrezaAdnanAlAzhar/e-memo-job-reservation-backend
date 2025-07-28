package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type JobService struct {
	jobRepo      *repository.JobRepository
	employeeRepo *repository.EmployeeRepository
	db           *sql.DB
}

func NewJobService(jobRepo *repository.JobRepository, employeeRepo *repository.EmployeeRepository, db *sql.DB) *JobService {
	return &JobService{
		jobRepo:      jobRepo,
		employeeRepo: employeeRepo,
		db:           db,
	}
}

// AssignPIC
func (s *JobService) AssignPIC(ctx context.Context, jobID int, req dto.AssignPICRequest, userNPK string) error {
	user, err := s.employeeRepo.FindByNPK(userNPK)
	if err != nil {
		return errors.New("action performer not found")
	}

	job, err := s.jobRepo.FindByID(jobID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("job not found")
		}
		return err
	}

	newPIC, err := s.employeeRepo.FindByNPK(req.PicJobNPK)
	if err != nil {
		return errors.New("new PIC employee data not found")
	}

	if user.DepartmentID != job.AssignedDepartmentID {
		return errors.New("user is not authorized to assign PIC for this job's department")
	}

	if newPIC.DepartmentID != job.AssignedDepartmentID {
		return errors.New("new PIC must be from the same department as the job")
	}

	return s.jobRepo.AssignPIC(jobID, req.PicJobNPK)
}

// ReorderJobs
func (s *JobService) ReorderJobs(ctx context.Context, req dto.ReorderJobsRequest, userNPK string) error {
	user, err := s.employeeRepo.FindByNPK(userNPK)
	if err != nil {
		return errors.New("action performer not found")
	}

	if user.DepartmentID != req.DepartmentTargetID {
		return errors.New("user can only reorder jobs within their own department")
	}

	validJobCount, err := s.jobRepo.CheckJobsInDepartment(req.OrderedJobIDs, req.DepartmentTargetID)
	if err != nil {
		return err
	}
	if validJobCount != len(req.OrderedJobIDs) {
		return errors.New("one or more job IDs do not belong to the specified department")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, jobID := range req.OrderedJobIDs {
		newPriority := i + 1 
		if err := s.jobRepo.UpdatePriority(ctx, tx, jobID, newPriority); err != nil {
			return err
		}
	}

	return tx.Commit()
}
