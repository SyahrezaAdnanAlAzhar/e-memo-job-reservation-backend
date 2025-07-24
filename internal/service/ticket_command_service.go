package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type TicketCommandService struct {
	db                    *sql.DB
	ticketRepo            *repository.TicketRepository
	jobRepo               *repository.JobRepository
	workflowRepo          *repository.WorkflowRepository
	trackStatusTicketRepo *repository.TrackStatusTicketRepository
	employeeRepo          *repository.EmployeeRepository
}

type TicketCommandServiceConfig struct {
	DB                    *sql.DB
	TicketRepo            *repository.TicketRepository
	JobRepo               *repository.JobRepository
	WorkflowRepo          *repository.WorkflowRepository
	TrackStatusTicketRepo *repository.TrackStatusTicketRepository
	EmployeeRepo          *repository.EmployeeRepository
}

func NewTicketCommandService(cfg *TicketCommandServiceConfig) *TicketCommandService {
	return &TicketCommandService{
		db:                    cfg.DB,
		ticketRepo:            cfg.TicketRepo,
		jobRepo:               cfg.JobRepo,
		workflowRepo:          cfg.WorkflowRepo,
		trackStatusTicketRepo: cfg.TrackStatusTicketRepo,
		employeeRepo:          cfg.EmployeeRepo,
	}
}

// CREATE TICKET
func (s *TicketCommandService) CreateTicket(ctx context.Context, req dto.CreateTicketRequest, requestor string) (*model.Ticket, error) {
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

	return createdTicket, nil
}

// UPDATE TICKET
func (s *TicketCommandService) UpdateTicket(ctx context.Context, ticketID int, req dto.UpdateTicketRequest, userNPK string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// GET TICKET
	originalTicket, err := s.ticketRepo.FindByIDAsStruct(ctx, ticketID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("ticket not found")
		}
		return err
	}

	// VALIDATE: ONLY THIS TICKET REQUESTOR THAT CAN EDIT
	if originalTicket.Requestor != userNPK {
		return errors.New("user is not authorized to edit this ticket")
	}

	// GET CURRENT STATUS
	_, currentStatusName, err := s.trackStatusTicketRepo.GetCurrentStatusByTicketID(ctx, ticketID)
	if err != nil {
		return errors.New("could not retrieve current ticket status")
	}

	// CHECK THE TICKET ABLE TO EDIT OR NOT
	canEdit := false
	switch currentStatusName {
	case "Ditolak":
		canEdit = true
	case "Menunggu Job":
		isAssigned, err := s.jobRepo.IsJobAssigned(ctx, ticketID)
		if err != nil {
			return errors.New("could not verify job assignment status")
		}
		if !isAssigned {
			canEdit = true
		}
	}

	if !canEdit {
		return errors.New("ticket cannot be edited in its current state")
	}

	// EXECUTE UPDATE
	if err := s.ticketRepo.Update(ctx, tx, ticketID, req); err != nil {
		return err
	}

	// GET USER POSITION
	positionID, err := s.employeeRepo.GetEmployeePositionID(ctx, userNPK)
	if err != nil {
		return err
	}

	// GET INITAL STATUS
	initialStatusID, err := s.workflowRepo.GetInitialStatusByPosition(ctx, positionID)
	if err != nil {
		return err
	}

	// CHANGE TICKET STATUS
	if err := s.trackStatusTicketRepo.UpdateStatus(ctx, tx, ticketID, initialStatusID); err != nil {
		return err
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
