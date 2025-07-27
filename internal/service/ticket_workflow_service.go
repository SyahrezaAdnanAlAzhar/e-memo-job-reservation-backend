package service

import (
	"context"
	"database/sql"
	"errors"

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
func (s *TicketWorkflowService) ExecuteAction(ctx context.Context, ticketID int, userNPK string, req dto.ExecuteActionRequest, filePath string) error {
	// GET ALL DATA
	user, err := s.employeeRepo.FindByNPK(userNPK)
	if err != nil {
		return errors.New("user not found")
	}

	ticket, err := s.ticketRepo.FindByIDAsStruct(ctx, ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}

	requestor, err := s.employeeRepo.FindByNPK(ticket.Requestor)
	if err != nil {
		return errors.New("original requestor not found")
	}

	job, _ := s.jobRepo.GetPicByTicketID(ctx, ticketID)

	currentStatusID, _, err := s.trackStatusTicketRepo.GetCurrentStatusByTicketID(ctx, ticketID)
	if err != nil {
		return errors.New("could not get current ticket status")
	}

	// ACTOR RESOLUTION
	userContexts := s.determineUserContexts(user, ticket, requestor, job)
	actorRoles, err := s.actorRoleMappingRepo.GetRolesForUserContext(user.EmployeePositionID, userContexts)
	if err != nil {
		return err
	}

	if job != "" && user.NPK == job {
		actorRoles = append(actorRoles, "ASSIGNED_PIC")
	}

	// VALIDATE TRANSITION
	transition, err := s.statusTransitionRepo.FindValidTransition(currentStatusID, req.ActionName)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("action not allowed from the current status")
		}
		return err
	}

	// AUTHORIZATION
	requiredRole, err := s.actorRoleRepo.GetRoleNameByID(transition.ActorRoleID)
	if err != nil {
		return err
	}

	isAuthorized := false
	for _, role := range actorRoles {
		if role == requiredRole {
			isAuthorized = true
			break
		}
	}
	if !isAuthorized {
		return errors.New("user does not have the required role for this action")
	}

	// INPUT VALIDATION
	if transition.RequiresReason && req.Reason == "" {
		return errors.New("reason is required for this action")
	}
	if transition.RequiresFile && filePath == "" {
		return errors.New("file upload is required for this action")
	}

	// EXEC
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// SAVE FILE
	var filePaths pq.StringArray
	if filePath != "" {
		filePaths = pq.StringArray{filePath}
	}

	logEntry := model.TicketActionLog{
		FilePath: filePaths,
	}
	if err := s.ticketActionLogRepo.Create(ctx, tx, logEntry); err != nil {
		return err
	}

	// CHANGE STATUS TICKET
	if err := s.trackStatusTicketRepo.UpdateStatus(ctx, tx, ticketID, transition.ToStatusID); err != nil {
		return err
	}

	return tx.Commit()
}

// HELPER FUNCTION
func (s *TicketWorkflowService) determineUserContexts(user *model.Employee, ticket *model.Ticket, requestor *model.Employee, jobPIC string) []string {
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
	return contexts
}
