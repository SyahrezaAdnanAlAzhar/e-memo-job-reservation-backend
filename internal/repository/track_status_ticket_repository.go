package repository

import (
	"context"
	"database/sql"
)

type TrackStatusTicketRepository struct {
	DB *sql.DB
}

func NewTrackStatusTicketRepository(db *sql.DB) *TrackStatusTicketRepository {
	return &TrackStatusTicketRepository{DB: db}
}


// INITIATE FIRST STATUS
func (r *TrackStatusTicketRepository) CreateInitialStatus(ctx context.Context, tx *sql.Tx, ticketID int, statusID int) error {
	query := "INSERT INTO track_status_ticket (ticket_id, status_ticket_id, start_date) VALUES ($1, $2, NOW())"
	_, err := tx.ExecContext(ctx, query, ticketID, statusID)
	return err
}


// UPDATE STATUS
func (r *TrackStatusTicketRepository) UpdateStatus(ctx context.Context, tx *sql.Tx, ticketID int, newStatusID int) error {
	queryFinish := "UPDATE track_status_ticket SET finish_date = NOW() WHERE ticket_id = $1 AND finish_date IS NULL"
	_, err := tx.ExecContext(ctx, queryFinish, ticketID)
	if err != nil {
		return err
	}

	queryCreate := "INSERT INTO track_status_ticket (ticket_id, status_ticket_id, start_date) VALUES ($1, $2, NOW())"
	_, err = tx.ExecContext(ctx, queryCreate, ticketID, newStatusID)
	return err
}