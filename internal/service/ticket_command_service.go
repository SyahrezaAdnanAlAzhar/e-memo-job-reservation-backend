package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/websocket"
	"github.com/lib/pq"
)

type TicketCommandService struct {
	db                    *sql.DB
	ticketRepo            *repository.TicketRepository
	jobRepo               *repository.JobRepository
	workflowRepo          *repository.WorkflowRepository
	trackStatusTicketRepo *repository.TrackStatusTicketRepository
	employeeRepo          *repository.EmployeeRepository
	departmentRepo        *repository.DepartmentRepository
	actorRoleMappingRepo  *repository.ActorRoleMappingRepository
	actorRoleRepo         *repository.ActorRoleRepository
	statusTransitionRepo  *repository.StatusTransitionRepository
	hub                   *websocket.Hub
	queryService          *TicketQueryService
}

type TicketCommandServiceConfig struct {
	DB                    *sql.DB
	TicketRepo            *repository.TicketRepository
	JobRepo               *repository.JobRepository
	WorkflowRepo          *repository.WorkflowRepository
	TrackStatusTicketRepo *repository.TrackStatusTicketRepository
	EmployeeRepo          *repository.EmployeeRepository
	DepartmentRepo        *repository.DepartmentRepository
	ActorRoleMappingRepo  *repository.ActorRoleMappingRepository
	ActorRoleRepo         *repository.ActorRoleRepository
	StatusTransitionRepo  *repository.StatusTransitionRepository
	Hub                   *websocket.Hub
	QueryService          *TicketQueryService
}

func NewTicketCommandService(cfg *TicketCommandServiceConfig) *TicketCommandService {
	return &TicketCommandService{
		db:                    cfg.DB,
		ticketRepo:            cfg.TicketRepo,
		jobRepo:               cfg.JobRepo,
		workflowRepo:          cfg.WorkflowRepo,
		trackStatusTicketRepo: cfg.TrackStatusTicketRepo,
		employeeRepo:          cfg.EmployeeRepo,
		departmentRepo:        cfg.DepartmentRepo,
		actorRoleMappingRepo:  cfg.ActorRoleMappingRepo,
		actorRoleRepo:         cfg.ActorRoleRepo,
		statusTransitionRepo:  cfg.StatusTransitionRepo,
		hub:                   cfg.Hub,
		queryService:          cfg.QueryService,
	}
}

// CREATE TICKET
func (s *TicketCommandService) CreateTicket(ctx context.Context, req dto.CreateTicketRequest, requestor string, supportFiles []string) (*model.Ticket, error) {
	// VALIDATE DEPARTMENT
	canReceive, err := s.departmentRepo.IsReceiver(req.DepartmentTargetID)
	if err != nil {
		return nil, err // "department not found or is not active"
	}
	if !canReceive {
		return nil, errors.New("selected target department cannot receive jobs")
	}

	// GET EMPLOYEE DATA (TO GET THE POSITION)
	positionID, err := s.employeeRepo.GetEmployeePositionID(ctx, requestor)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("requestor not found")
		}
		return nil, err
	}

	// GET INITIAL STATUS FROM PREVIOUS GET EMPLOYEE DATA
	initialStatusID, err := s.workflowRepo.GetInitialStatusByPosition(ctx, positionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no workflow defined for this user's position")
		}
		return nil, err
	}

	deadline, err := repository.ParseDeadline(req.Deadline)
	if err != nil {
		return nil, errors.New("invalid deadline format, please use YYYY-MM-DD")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// GET LAST PRIORITY
	lastPriority, err := s.ticketRepo.GetLastPriority(ctx, tx, req.DepartmentTargetID)
	if err != nil {
		return nil, err
	}

	ticketData := model.Ticket{
		Requestor:           requestor,
		DepartmentTargetID:  req.DepartmentTargetID,
		PhysicalLocationID:  toNullInt64(req.PhysicalLocationID),
		SpecifiedLocationID: toNullInt64(req.SpecifiedLocationID),
		Description:         req.Description,
		TicketPriority:      lastPriority,
		Deadline:            deadline,
		SupportFile:         pq.StringArray(supportFiles),
	}

	// INSERT DATA TO TICKET TABLE
	createdTicket, err := s.ticketRepo.Create(ctx, tx, ticketData)
	if err != nil {
		return nil, err
	}

	// INSERT DATA TO JOB TABLE
	err = s.jobRepo.Create(ctx, tx, createdTicket.ID, createdTicket.TicketPriority)
	if err != nil {
		return nil, err
	}

	// INITIATE FIRST STATUS
	err = s.trackStatusTicketRepo.CreateInitialStatus(ctx, tx, createdTicket.ID, initialStatusID)
	if err != nil {
		return nil, err
	}

	// COMMIT DATA
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	ticketDetail, err := s.queryService.GetTicketByID(createdTicket.ID)
	if err != nil {
		log.Printf("CRITICAL: Failed to fetch new ticket details for broadcast. TicketID: %d, Error: %v", createdTicket.ID, err)
	} else {
		message, err := websocket.NewMessage("TICKET_CREATED", ticketDetail)
		if err != nil {
			log.Printf("CRITICAL: Failed to create websocket message for new ticket: %v", err)
		} else {
			s.hub.BroadcastMessage(message)
		}
	}

	return createdTicket, err
}

// UPDATE TICKET
func (s *TicketCommandService) UpdateTicket(ctx context.Context, ticketID int, req dto.UpdateTicketRequest, userNPK string) error {
	originalTicket, err := s.ticketRepo.FindByIDAsStruct(ctx, ticketID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("ticket not found")
		}
		return err
	}

	user, err := s.employeeRepo.FindByNPK(userNPK)
	if err != nil {
		return errors.New("user not found")
	}

	requestor, err := s.employeeRepo.FindByNPK(originalTicket.Requestor)
	if err != nil {
		return errors.New("original requestor not found")
	}

	isOriginalRequestor := user.NPK == originalTicket.Requestor
	isSameDeptApprover := user.DepartmentID == requestor.DepartmentID && (user.Position.Name == "Head of Department" || user.Position.Name == "Section")

	if !isOriginalRequestor && !isSameDeptApprover {
		return errors.New("user is not authorized to edit this ticket")
	}

	currentStatusID, currentStatusName, err := s.trackStatusTicketRepo.GetCurrentStatusByTicketID(ctx, ticketID)
	if err != nil {
		return errors.New("could not retrieve current ticket status")
	}

	if currentStatusName != "Ditolak" && currentStatusName != "Dibatalkan" {
		return errors.New("ticket can only be edited if status is 'Ditolak'")
	}

	userContexts := determineUserContexts(user, originalTicket, requestor, nil)

	actorRoles, err := s.actorRoleMappingRepo.GetRolesForUserContext(user.Position.ID, userContexts)
	if err != nil {
		return err
	}

	actorRoleIDs, err := s.actorRoleRepo.GetRoleIDsByNames(actorRoles)
	if err != nil {
		return err
	}

	isAuthorized, err := s.statusTransitionRepo.HasAvailableActionsForRoles(currentStatusID, actorRoleIDs)
	if err != nil {
		return err
	}

	if !isAuthorized {
		return errors.New("user is not authorized to edit this ticket")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rowsAffected, err := s.ticketRepo.Update(ctx, tx, ticketID, req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			return errors.New("invalid deadline format, please use YYYY-MM-DD")
		}
		return err
	}
	if rowsAffected == 0 {
		return errors.New("data conflict: ticket has been modified by another user, please refresh")
	}

	if req.Deadline != nil {
		if _, err := time.Parse("2006-01-02", *req.Deadline); err != nil {
			return errors.New("invalid deadline format, please use YYYY-MM-DD")
		}
	}

	updatedTicketDetail, err := s.queryService.GetTicketByID(ticketID)
	if err != nil {
		log.Printf("CRITICAL: Failed to fetch updated ticket for broadcast after update. TicketID: %d, Error: %v", ticketID, err)
	} else {
		message, err := websocket.NewMessage("TICKET_UPDATED", updatedTicketDetail)
		if err != nil {
			log.Printf("CRITICAL: Failed to create websocket message for updated ticket: %v", err)
		} else {
			s.hub.BroadcastMessage(message)
		}
	}

	return tx.Commit()
}

// HELPER
// CONVERTER
func toNullInt64(val *int) sql.NullInt64 {
	if val == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*val), Valid: true}
}
