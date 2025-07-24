package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type TicketWorkflowService struct {
	db                    *sql.DB
	ticketRepo            *repository.TicketRepository
	employeeRepo          *repository.EmployeeRepository
	trackStatusTicketRepo *repository.TrackStatusTicketRepository
	statusTicketRepo      *repository.StatusTicketRepository
	rejectedTicketService *RejectedTicketService
	workflowRepo          *repository.WorkflowRepository
}

type TicketWorkflowServiceConfig struct {
	DB                    *sql.DB
	TicketRepo            *repository.TicketRepository
	EmployeeRepo          *repository.EmployeeRepository
	TrackStatusTicketRepo *repository.TrackStatusTicketRepository
	StatusTicketRepo      *repository.StatusTicketRepository
	RejectedTicketService *RejectedTicketService
	WorkflowRepo          *repository.WorkflowRepository
}

func NewTicketWorkflowService(cfg *TicketWorkflowServiceConfig) *TicketWorkflowService {
	return &TicketWorkflowService{
		db:                    cfg.DB,
		ticketRepo:            cfg.TicketRepo,
		employeeRepo:          cfg.EmployeeRepo,
		trackStatusTicketRepo: cfg.TrackStatusTicketRepo,
		statusTicketRepo:      cfg.StatusTicketRepo,
		rejectedTicketService: cfg.RejectedTicketService,
		workflowRepo:          cfg.WorkflowRepo,
	}
}

// CHANGE STATUS TO REJECT ("Ditolak")
func (s *TicketWorkflowService) RejectTicket(ctx context.Context, ticketID int, req dto.RejectTicketRequest, userNPK string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ticket, err := s.ticketRepo.FindByIDAsStruct(ctx, ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}

	user, err := s.employeeRepo.FindByNPK(userNPK)
	if err != nil {
		return errors.New("user not found")
	}

	isAllowed := user.DepartmentID == ticket.DepartmentTargetID && (user.Position.Name == "Head of Department" || user.Position.Name == "Section")
	if !isAllowed {
		return errors.New("user not authorized to reject this ticket")
	}

	_, err = s.statusTicketRepo.FindByName("Ditolak")
	if err != nil {
		return errors.New("critical configuration error: 'Ditolak' status not found")
	}

	rejectionReq := dto.CreateRejectedTicketRequest{
		TicketID: int64(ticketID),
		Feedback: req.Reason,
	}

	_, err = s.rejectedTicketService.CreateRejectedTicket(ctx, rejectionReq, userNPK)

	return err
}

// CHANGE STATUS TO CANCEL ("Dibatalkan")
func (s *TicketWorkflowService) CancelTicket(ctx context.Context, ticketID int, userNPK string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ticket, err := s.ticketRepo.FindByIDAsStruct(ctx, ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}

	user, err := s.employeeRepo.FindByNPK(userNPK)
	if err != nil {
		return errors.New("user not found")
	}

	requestor, err := s.employeeRepo.FindByNPK(ticket.Requestor)
	if err != nil {
		return errors.New("original requestor not found")
	}

	isOriginalRequestor := user.NPK == ticket.Requestor
	isSameDeptApprover := user.DepartmentID == requestor.DepartmentID && (user.Position.Name == "Head of Department" || user.Position.Name == "Section")

	if !isOriginalRequestor && !isSameDeptApprover {
		return errors.New("user not authorized to cancel this ticket")
	}

	cancelledStatus, err := s.statusTicketRepo.FindByName("Dibatalkan")
	if err != nil {
		return errors.New("critical configuration error: 'Dibatalkan' status not found")
	}

	if err := s.trackStatusTicketRepo.UpdateStatus(ctx, tx, ticketID, cancelledStatus.ID); err != nil {
		return err
	}

	return tx.Commit()
}

// CHANGE STATUS TO APPROVAL SECTION ("Approval Section")
func (s *TicketWorkflowService) ApproveSection(ctx context.Context, ticketID int, userNPK string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ticket, err := s.ticketRepo.FindByIDAsStruct(ctx, ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}

	user, err := s.employeeRepo.FindByNPK(userNPK)
	if err != nil {
		return errors.New("user not found")
	}

	requestor, err := s.employeeRepo.FindByNPK(ticket.Requestor)
	if err != nil {
		return errors.New("original requestor not found")
	}

	isAllowed := user.DepartmentID == requestor.DepartmentID && (user.Position.Name == "Head of Department" || user.Position.Name == "Section")
	if !isAllowed {
		return errors.New("user is not authorized to perform this approval")
	}

	currentStatusID, currentStatusName, err := s.trackStatusTicketRepo.GetCurrentStatusByTicketID(ctx, ticketID)
	if err != nil {
		return err
	}
	if currentStatusName != "Approval Section" {
		return errors.New("ticket is not in 'Approval Section' status")
	}

	nextStatusID, isFinalStep, err := s.workflowRepo.GetNextWorkflowStep(ctx, currentStatusID)
	if err != nil {
		return err
	}
	if isFinalStep {
		return errors.New("workflow configuration error: no next step found")
	}

	if err := s.trackStatusTicketRepo.UpdateStatus(ctx, tx, ticketID, nextStatusID); err != nil {
		return err
	}

	return tx.Commit()
}

// CHANGE STATUS TO APPROVAL DEPARTMENT ("Approval Department")
func (s *TicketWorkflowService) ApproveDepartment(ctx context.Context, ticketID int, userNPK string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ticket, err := s.ticketRepo.FindByIDAsStruct(ctx, ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}

	user, err := s.employeeRepo.FindByNPK(userNPK)
	if err != nil {
		return errors.New("user not found")
	}

	requestor, err := s.employeeRepo.FindByNPK(ticket.Requestor)
	if err != nil {
		return errors.New("original requestor not found")
	}

	isAllowed := user.DepartmentID == requestor.DepartmentID && user.Position.Name == "Head of Department"
	if !isAllowed {
		return errors.New("user is not authorized to perform this approval")
	}

	currentStatusID, currentStatusName, err := s.trackStatusTicketRepo.GetCurrentStatusByTicketID(ctx, ticketID)
	if err != nil {
		return err
	}
	if currentStatusName != "Approval Department" {
		return errors.New("ticket is not in 'Approval Department' status")
	}

	nextStatusID, isFinalStep, err := s.workflowRepo.GetNextWorkflowStep(ctx, currentStatusID)
	if err != nil {
		return err
	}
	if isFinalStep {
		return errors.New("workflow configuration error: no next step found")
	}

	if err := s.trackStatusTicketRepo.UpdateStatus(ctx, tx, ticketID, nextStatusID); err != nil {
		return err
	}

	return tx.Commit()
}
