package repository

import (
	"context"
	"database/sql"
	"time"
)

type Ticket struct {
	ID                  int           `json:"id"`
	Requestor           string        `json:"requestor"`
	DepartmentTargetID  int           `json:"department_target_id"`
	PhysicalLocationID  sql.NullInt64 `json:"physical_location_id"`
	SpecifiedLocationID sql.NullInt64 `json:"specified_location_id"`
	Description         string        `json:"description"`
	TicketPriority      int           `json:"ticket_priority"`
	SupportFile         []string      `json:"support_file"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
}

type CreateTicketRequest struct {
	DepartmentTargetID  int    `json:"department_target_id" binding:"required,gt=0"`
	PhysicalLocationID  *int   `json:"physical_location_id"`
	SpecifiedLocationID *int   `json:"specified_location_id"`
	Description         string `json:"description" binding:"required"`
}

type TicketRepository struct {
	DB *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{DB: db}
}

// HELPER

// GET LAST PRIORITY
func (r *TicketRepository) GetLastPriority(ctx context.Context, tx *sql.Tx, departmentTargetID int) (int, error) {
	var lastPriority sql.NullInt64
	query := "SELECT MAX(ticket_priority) FROM ticket WHERE department_target_id = $1"
	err := tx.QueryRowContext(ctx, query, departmentTargetID).Scan(&lastPriority)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if !lastPriority.Valid { 
		return 1, nil 
	}
	return int(lastPriority.Int64) + 1, nil
}

// MAIN

// CREATE TICKET
func (r *TicketRepository) Create(ctx context.Context, tx *sql.Tx, ticket Ticket) (*Ticket, error) {
	query := `
        INSERT INTO ticket (requestor, department_target_id, physical_location_id, specified_location_id, description, ticket_priority)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at`

	row := tx.QueryRowContext(ctx, query,
		ticket.Requestor,
		ticket.DepartmentTargetID,
		ticket.PhysicalLocationID,
		ticket.SpecifiedLocationID,
		ticket.Description,
		ticket.TicketPriority,
	)

	var newTicket Ticket = ticket 
	err := row.Scan(&newTicket.ID, &newTicket.CreatedAt, &newTicket.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &newTicket, nil
}

// INITIATE FIRST STATUS
func (r *TicketRepository) CreateInitialStatus(ctx context.Context, tx *sql.Tx, ticketID int, statusID int) error {
	query := "INSERT INTO track_status_ticket (ticket_id, status_ticket_id, start_date) VALUES ($1, $2, NOW())"
	_, err := tx.ExecContext(ctx, query, ticketID, statusID)
	return err
}
