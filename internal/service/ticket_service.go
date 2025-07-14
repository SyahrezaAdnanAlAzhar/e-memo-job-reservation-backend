package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type TicketService struct {
	ticketRepo   *repository.TicketRepository
	jobRepo      *repository.JobRepository
	workflowRepo *repository.WorkflowRepository
	db           *sql.DB
}

type TicketServiceConfig struct {
	TicketRepo   *repository.TicketRepository
	JobRepo      *repository.JobRepository
	WorkflowRepo *repository.WorkflowRepository
	DB           *sql.DB
}

func NewTicketService(cfg *TicketServiceConfig) *TicketService {
	return &TicketService{
		ticketRepo:   cfg.TicketRepo,
		jobRepo:      cfg.JobRepo,
		workflowRepo: cfg.WorkflowRepo,
		db:           cfg.DB,
	}
}

// CREATE TICKET
func (s *TicketService) CreateTicket(ctx context.Context, req repository.CreateTicketRequest, requestor string) (*repository.Ticket, error) {
	// GET EMPLOYEE DATA (TO GET THE POSITION)
	positionID, err := s.workflowRepo.GetEmployeeData(ctx, requestor)
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

	ticketData := repository.Ticket{
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
	err = s.ticketRepo.CreateInitialStatus(ctx, tx, createdTicket.ID, initialStatusID)
	if err != nil {
		return nil, err
	}

	// COMMIT DATA
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return createdTicket, nil
}

// HELPER CONVERTER
func toNullInt64(val *int) sql.NullInt64 {
	if val == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*val), Valid: true}
}
