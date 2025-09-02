package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/lib/pq"
)

type TicketWorkflowService struct {
	db                    *sql.DB
	ticketRepo            *repository.TicketRepository
	jobRepo               *repository.JobRepository
	employeeRepo          *repository.EmployeeRepository
	trackStatusTicketRepo *repository.TrackStatusTicketRepository
	statusTicketRepo      *repository.StatusTicketRepository
	statusTransitionRepo  *repository.StatusTransitionRepository
	actorRoleRepo         *repository.ActorRoleRepository
	actorRoleMappingRepo  *repository.ActorRoleMappingRepository
	ticketActionLogRepo   *repository.TicketActionLogRepository
	actionService         *TicketActionService
}

type TicketWorkflowServiceConfig struct {
	DB                    *sql.DB
	TicketRepo            *repository.TicketRepository
	JobRepo               *repository.JobRepository
	EmployeeRepo          *repository.EmployeeRepository
	TrackStatusTicketRepo *repository.TrackStatusTicketRepository
	StatusTicketRepo      *repository.StatusTicketRepository
	StatusTransitionRepo  *repository.StatusTransitionRepository
	ActorRoleRepo         *repository.ActorRoleRepository
	ActorRoleMappingRepo  *repository.ActorRoleMappingRepository
	TicketActionLogRepo   *repository.TicketActionLogRepository
	ActionService         *TicketActionService
}

func NewTicketWorkflowService(cfg *TicketWorkflowServiceConfig) *TicketWorkflowService {
	return &TicketWorkflowService{
		db:                    cfg.DB,
		ticketRepo:            cfg.TicketRepo,
		jobRepo:               cfg.JobRepo,
		employeeRepo:          cfg.EmployeeRepo,
		trackStatusTicketRepo: cfg.TrackStatusTicketRepo,
		statusTicketRepo:      cfg.StatusTicketRepo,
		statusTransitionRepo:  cfg.StatusTransitionRepo,
		actorRoleRepo:         cfg.ActorRoleRepo,
		actorRoleMappingRepo:  cfg.ActorRoleMappingRepo,
		ticketActionLogRepo:   cfg.TicketActionLogRepo,
		actionService:         cfg.ActionService,
	}
}

// EXECUTE ACTION TO GET TO THE NEXT STATUS BASED ON STATE
func (s *TicketWorkflowService) ExecuteAction(ctx context.Context, ticketID int, userNPK string, req dto.ExecuteActionRequest, filesMetadata []model.FileMetadata) error {
	availableActions, err := s.actionService.GetAvailableActions(ctx, ticketID, userNPK)
	if err != nil {
		return err
	}

	var selectedAction *dto.AvailableTicketActionResponse
	for _, action := range availableActions {
		if action.ActionName == req.ActionName {
			act := action
			selectedAction = &act
			break
		}
	}

	if selectedAction == nil {
		return errors.New("user does not have the required role or action is not allowed from the current status")
	}

	if selectedAction.RequireReason && req.Reason == "" {
		errorMsg := "reason is required for this action"
		if selectedAction.ReasonLabel != nil {
			errorMsg = fmt.Sprintf("%s is required", *selectedAction.ReasonLabel)
		}
		return errors.New(errorMsg)
	}

	if selectedAction.RequireFile && len(filesMetadata) == 0 {
		return errors.New("file upload is required for this action")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if req.ActionName == "Selesaikan Job" {
		if len(filesMetadata) > 0 {
			err := s.jobRepo.UpdateJobCompletionDetails(ctx, tx, ticketID, filesMetadata, req.SpendingAmount)
			if err != nil {
				if err == sql.ErrNoRows {
					return errors.New("job associated with this ticket not found")
				}
				return errors.New("failed to update job completion details")
			}
		}
	}

	currentStatusID, _, _ := s.trackStatusTicketRepo.GetCurrentStatusByTicketID(ctx, ticketID)

	var filePathsForLog []string
	for _, meta := range filesMetadata {
		filePathsForLog = append(filePathsForLog, meta.FilePath)
	}

	logEntry := model.TicketActionLog{
		TicketID:       int64(ticketID),
		ActionID:       selectedAction.ActionID,
		PerformedByNpk: userNPK,
		DetailsText:    sql.NullString{String: req.Reason, Valid: req.Reason != ""},
		FilePath:       pq.StringArray(filePathsForLog),
		FromStatusID:   sql.NullInt32{Int32: int32(currentStatusID), Valid: true},
		ToStatusID:     selectedAction.ToStatusID,
	}
	if err := s.ticketActionLogRepo.Create(ctx, tx, logEntry); err != nil {
		return err
	}

	if err := s.trackStatusTicketRepo.UpdateStatus(ctx, tx, ticketID, selectedAction.ToStatusID); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *TicketWorkflowService) ValidateAndGetTransition(ctx context.Context, currentStatusID int, actionName string) (toStatusID int, allowedRoleIDs []int, err error) {
	toStatusID, allowedRoleIDs, err = s.statusTransitionRepo.FindValidTransition(currentStatusID, actionName)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil, errors.New("action not allowed from the current status")
		}
		return 0, nil, err
	}

	if len(allowedRoleIDs) == 0 {
		return 0, nil, errors.New("no roles are configured to perform this action")
	}

	return toStatusID, allowedRoleIDs, nil
}
