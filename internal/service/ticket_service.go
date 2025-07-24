package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type TicketService struct {
	ticketRepo            *repository.TicketRepository
	jobRepo               *repository.JobRepository
	workflowRepo          *repository.WorkflowRepository
	trackStatusTicketRepo *repository.TrackStatusTicketRepository
	employeeRepo          *repository.EmployeeRepository
	statusTicketRepo      *repository.StatusTicketRepository
	db                    *sql.DB
}

type TicketServiceConfig struct {
	TicketRepo            *repository.TicketRepository
	JobRepo               *repository.JobRepository
	WorkflowRepo          *repository.WorkflowRepository
	TrackStatusTicketRepo *repository.TrackStatusTicketRepository
	EmployeeRepo          *repository.EmployeeRepository
	StatusTicketRepo      *repository.StatusTicketRepository
	DB                    *sql.DB
}

func NewTicketService(cfg *TicketServiceConfig) *TicketService {
	return &TicketService{
		ticketRepo:            cfg.TicketRepo,
		jobRepo:               cfg.JobRepo,
		workflowRepo:          cfg.WorkflowRepo,
		trackStatusTicketRepo: cfg.TrackStatusTicketRepo,
		employeeRepo:          cfg.EmployeeRepo,
		statusTicketRepo:      cfg.StatusTicketRepo,
		db:                    cfg.DB,
	}
}

// CREATE TICKET
func (s *TicketService) CreateTicket(ctx context.Context, req dto.CreateTicketRequest, requestor string) (*model.Ticket, error) {
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

// GET ALL
func (s *TicketService) GetAllTickets(filters map[string]string) ([]map[string]interface{}, error) {
	allowedFilters := map[string]bool{
		"department_target_id": true,
		"current_status":       true,
		"requestor_npk":        true,
	}
	for key := range filters {
		if !allowedFilters[key] {
			return nil, errors.New("invalid filter key: " + key)
		}
	}
	return s.ticketRepo.FindAll(filters)
}

// GET BY ID
func (s *TicketService) GetTicketByID(id int) (map[string]interface{}, error) {
	return s.ticketRepo.FindByID(id)
}

// UPDATE TICKET
func (s *TicketService) UpdateTicket(ctx context.Context, ticketID int, req dto.UpdateTicketRequest, userNPK string) error {
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

// RE ORDER
func (s *TicketService) ReorderTickets(ctx context.Context, req dto.ReorderTicketsRequest) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, ticketID := range req.OrderedTicketIDs {
		newPriority := i + 1
		if err := s.ticketRepo.UpdatePriority(ctx, tx, ticketID, newPriority); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// CHANGE STATUS TO REJECT ("Ditolak")
func (s *TicketService) RejectTicket(ctx context.Context, ticketID int, req dto.RejectTicketRequest, userNPK string) error {
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

	rejectedStatus, err := s.statusTicketRepo.FindByName("Ditolak")
	if err != nil {
		return errors.New("critical configuration error: 'Ditolak' status not found")
	}

	if err := s.trackStatusTicketRepo.UpdateStatus(ctx, tx, ticketID, rejectedStatus.ID); err != nil {
		return err
	}

	log.Printf("Ticket %d rejected by %s. Reason: %s", ticketID, userNPK, req.Reason)

	return tx.Commit()
}

// CHANGE STATUS TO CANCEL ("Dibatalkan")
func (s *TicketService) CancelTicket(ctx context.Context, ticketID int, userNPK string) error {
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

// HELPER

// CONVERTER
func toNullInt64(val *int) sql.NullInt64 {
	if val == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*val), Valid: true}
}
