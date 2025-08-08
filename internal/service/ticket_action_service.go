package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type TicketActionService struct {
	ticketRepo            *repository.TicketRepository
	jobRepo               *repository.JobRepository
	employeeRepo          *repository.EmployeeRepository
	trackStatusTicketRepo *repository.TrackStatusTicketRepository
	statusTransitionRepo  *repository.StatusTransitionRepository
	actorRoleMappingRepo  *repository.ActorRoleMappingRepository
}

type TicketActionServiceConfig struct {
	TicketRepo            *repository.TicketRepository
	JobRepo               *repository.JobRepository
	EmployeeRepo          *repository.EmployeeRepository
	TrackStatusTicketRepo *repository.TrackStatusTicketRepository
	StatusTransitionRepo  *repository.StatusTransitionRepository
	ActorRoleMappingRepo  *repository.ActorRoleMappingRepository
}

func NewTicketActionService(cfg *TicketActionServiceConfig) *TicketActionService {
	return &TicketActionService{
		ticketRepo:            cfg.TicketRepo,
		jobRepo:               cfg.JobRepo,
		employeeRepo:          cfg.EmployeeRepo,
		trackStatusTicketRepo: cfg.TrackStatusTicketRepo,
		statusTransitionRepo:  cfg.StatusTransitionRepo,
		actorRoleMappingRepo:  cfg.ActorRoleMappingRepo,
	}
}

// GET AVAILABLE ACTIONS
func (s *TicketActionService) GetAvailableActions(ctx context.Context, ticketID int, userNPK string) ([]dto.ActionResponse, error) {
	user, err := s.employeeRepo.FindByNPK(userNPK)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user employee not found")
		}
		return nil, err
	}

	ticket, err := s.ticketRepo.FindByIDAsStruct(ctx, ticketID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}

	requestor, err := s.employeeRepo.FindByNPK(ticket.Requestor)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("requestor employee not found")
		}
		return nil, err
	}

	job, _ := s.jobRepo.FindByTicketID(ctx, ticketID)

	userContexts := determineUserContexts(user, ticket, requestor, job)

	userRoleIDs, err := s.actorRoleMappingRepo.GetRoleIDsForUserContext(user.Position.ID, userContexts)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("actor role mapping not found")
		}
		return nil, err
	}

	currentStatusID, _, err := s.trackStatusTicketRepo.GetCurrentStatusByTicketID(ctx, ticketID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("current status not found")
		}
		return nil, err
	}

	availableActions, err := s.statusTransitionRepo.FindAvailableTransitionsForRoles(currentStatusID, userRoleIDs)
	if err != nil {
		return nil, err
	}

	return availableActions, nil
}
