package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/gin-gonic/gin"
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
	}
}

// EXECUTE ACTION TO GET TO THE NEXT STATUS BASED ON STATE
func (s *TicketWorkflowService) ExecuteAction(c *gin.Context, ticketID int, userNPK string, req dto.ExecuteActionRequest, filePaths []string) error {
	log.Printf("--- START ExecuteAction for Ticket %d by User %s ---", ticketID, userNPK)
	ctx := c.Request.Context()

	user, err := s.employeeRepo.FindByNPK(userNPK)
	if err != nil {
		return errors.New("user not found")
	}
	log.Printf("Action Performer: %+v", user)

	ticket, err := s.ticketRepo.FindByIDAsStruct(ctx, ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}
	log.Printf("Target Ticket: %+v", ticket)

	requestor, err := s.employeeRepo.FindByNPK(ticket.Requestor)
	if err != nil {
		return errors.New("original requestor not found")
	}
	log.Printf("Original Requestor: %+v", requestor)

	job, _ := s.jobRepo.FindByTicketID(ctx, ticketID)

	currentStatusID, currentStatusName, err := s.trackStatusTicketRepo.GetCurrentStatusByTicketID(ctx, ticketID)
	if err != nil {
		return errors.New("could not get current ticket status")
	}
	log.Printf("Current Status: ID=%d, Name=%s", currentStatusID, currentStatusName)

	log.Printf("DEBUG: Executing action '%s' for ticket %d. Current status is '%s' (ID: %d)", req.ActionName, ticketID, currentStatusName, currentStatusID)

	toStatusID, allowedRoleIDs, err := s.statusTransitionRepo.FindValidTransition(currentStatusID, req.ActionName)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("action not allowed from the current status")
		}
		return err
	}

	userContexts := s.determineUserContexts(user, ticket, requestor, job)
	userRoleIDs, err := s.actorRoleMappingRepo.GetRoleIDsForUserContext(user.Position.ID, userContexts)
	log.Printf("User Contexts: %v", userContexts)

	if err != nil {
		return err
	}

	log.Printf("Resolved Actor Roles: %v", userRoleIDs)

	isAuthorized := false
	for _, userRoleID := range userRoleIDs {
		for _, allowedRoleID := range allowedRoleIDs {
			if userRoleID == allowedRoleID {
				isAuthorized = true
				break
			}
		}
		if isAuthorized {
			break
		}
	}
	if !isAuthorized {
		return errors.New("user does not have the required role for this action")
	}

	log.Println("Authorization successful.")

	transitionDetails, err := s.statusTransitionRepo.GetTransitionDetails(currentStatusID, req.ActionName)
	if err != nil {
		return errors.New("could not retrieve transition details")
	}

	if transitionDetails.RequiresReason && req.Reason == "" {
		errorMsg := "reason is required for this action"
		if transitionDetails.ReasonLabel.Valid {
			errorMsg = fmt.Sprintf("%s is required", transitionDetails.ReasonLabel.String)
		}
		return errors.New(errorMsg)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	logEntry := model.TicketActionLog{
		TicketID:       int64(ticketID),
		ActionID:       transitionDetails.ActionID,
		PerformedByNpk: userNPK,
		DetailsText:    sql.NullString{String: req.Reason, Valid: req.Reason != ""},
		FilePath:       pq.StringArray(filePaths),
		FromStatusID:   sql.NullInt32{Int32: int32(currentStatusID), Valid: true},
		ToStatusID:     transitionDetails.ToStatusID,
	}

	if err := s.ticketActionLogRepo.Create(ctx, tx, logEntry); err != nil {
		return err
	}

	if err := s.trackStatusTicketRepo.UpdateStatus(ctx, tx, ticketID, toStatusID); err != nil {
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

// HELPER FUNCTION
func (s *TicketWorkflowService) determineUserContexts(user *model.Employee, ticket *model.Ticket, requestor *model.Employee, job *model.Job) []string {
	var contexts []string
	if user.NPK == ticket.Requestor {
		contexts = append(contexts, "SELF")
	}
	if user.DepartmentID == requestor.DepartmentID {
		contexts = append(contexts, "REQUESTOR_DEPT")
	}
	if user.DepartmentID == ticket.DepartmentTargetID {
		contexts = append(contexts, "TARGET_DEPT")
	}
	if job != nil && job.PicJob.Valid && user.NPK == job.PicJob.String {
		contexts = append(contexts, "ASSIGNED")
	}
	return contexts
}
