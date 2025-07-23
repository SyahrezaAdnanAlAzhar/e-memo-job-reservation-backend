package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type WorkflowService struct {
	workflowRepo *repository.WorkflowRepository
	stepRepo     *repository.WorkflowStepRepository
	db           *sql.DB
}

func NewWorkflowService(workflowRepo *repository.WorkflowRepository, stepRepo *repository.WorkflowStepRepository, db *sql.DB) *WorkflowService {
	return &WorkflowService{workflowRepo: workflowRepo, stepRepo: stepRepo, db: db}
}

// CREATE
func (s *WorkflowService) CreateWorkflowWithSteps(ctx context.Context, req dto.CreateWorkflowRequest) (*model.Workflow, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// CREATE WORKFLOW
	newWorkflow, err := s.workflowRepo.Create(ctx, tx, req.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errors.New("workflow name already exists")
		}
		return nil, err
	}

	// LOOP FOR EACH WORKFLOW STEP
	for i, statusID := range req.StatusTicketIDs {
		stepSequence := i 
		err := s.stepRepo.Create(ctx, tx, newWorkflow.ID, statusID, stepSequence)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23503" { // foreign_key_violation
					return nil, errors.New("one or more status_ticket_ids are invalid")
				}
				if pgErr.Code == "23505" { // unique_violation
					return nil, errors.New("cannot add the same status twice to a workflow")
				}
			}
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return newWorkflow, nil
}