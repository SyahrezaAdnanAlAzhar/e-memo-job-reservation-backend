package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type TicketPriorityService struct {
	db           *sql.DB
	ticketRepo   *repository.TicketRepository
	employeeRepo *repository.EmployeeRepository
}

func NewTicketPriorityService(db *sql.DB, ticketRepo *repository.TicketRepository, employeeRepo *repository.EmployeeRepository) *TicketPriorityService {
	return &TicketPriorityService{
		db:           db,
		ticketRepo:   ticketRepo,
		employeeRepo: employeeRepo,
	}
}

// RE ORDER
func (s *TicketPriorityService) ReorderTickets(ctx context.Context, req dto.ReorderTicketsRequest) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, item := range req.Items {
		newPriority := i + 1

		rowsAffected, err := s.ticketRepo.UpdatePriority(ctx, tx, item.TicketID, item.Version, newPriority)
		if err != nil {
			return err 
		}

		if rowsAffected == 0 {
			return errors.New("data conflict: one or more tickets have been modified by another user, please refresh")
		}
	}
	return tx.Commit()
}
