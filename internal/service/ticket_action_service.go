package service

import (
	"context"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
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
		return nil, err
	}
	ticket, err := s.ticketRepo.FindByIDAsStruct(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	requestor, err := s.employeeRepo.FindByNPK(ticket.Requestor)
	if err != nil {
		return nil, err
	}
	jobPIC, _ := s.jobRepo.GetPicByTicketID(ctx, ticketID)

	userContexts := determineUserContexts(user, ticket, requestor, jobPIC)
	userRoles, err := s.actorRoleMappingRepo.GetRolesForUserContext(user.EmployeePositionID, userContexts)
	if err != nil {
		return nil, err
	}
	if jobPIC != "" && user.NPK == jobPIC {
		userRoles = append(userRoles, "ASSIGNED_PIC")
	}
	userRolesMap := make(map[string]bool)
	for _, role := range userRoles {
		userRolesMap[role] = true
	}

	currentStatusID, _, err := s.trackStatusTicketRepo.GetCurrentStatusByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	possibleTransitions, err := s.statusTransitionRepo.FindPossibleTransitionsWithDetails(currentStatusID)
	if err != nil {
		return nil, err
	}

	var availableActions []dto.ActionResponse
	for _, transition := range possibleTransitions {
		if _, ok := userRolesMap[transition.RequiredActorRole]; ok {
			availableActions = append(availableActions, transition.ActionDetail)
		}
	}

	return availableActions, nil
}

func determineUserContexts(user *model.Employee, ticket *model.Ticket, requestor *model.Employee, jobPIC string) []string {
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
