package service

import (
	"context"
	"database/sql"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type TicketPriorityService struct {
	db         *sql.DB
	ticketRepo *repository.TicketRepository
}

// NewTicketPriorityService adalah constructor untuk TicketPriorityService.
func NewTicketPriorityService(db *sql.DB, ticketRepo *repository.TicketRepository) *TicketPriorityService {
	return &TicketPriorityService{
		db:         db,
		ticketRepo: ticketRepo,
	}
}

// RE ORDER
func (s *TicketPriorityService) ReorderTickets(ctx context.Context, req dto.ReorderTicketsRequest) error {
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
