package repository

import (
	"context"
	"database/sql"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
)

type TicketActionLogRepository struct {
	DB *sql.DB
}

func NewTicketActionLogRepository(db *sql.DB) *TicketActionLogRepository {
	return &TicketActionLogRepository{DB: db}
}

// Create mencatat sebuah aksi ke dalam log di dalam sebuah transaksi.
func (r *TicketActionLogRepository) Create(ctx context.Context, tx *sql.Tx, logEntry model.TicketActionLog) error {
	query := `
        INSERT INTO ticket_action_log (
            ticket_id, action_id, performed_by_npk, details_text, 
            file_path, from_status_id, to_status_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := tx.ExecContext(ctx, query,
		logEntry.TicketID,
		logEntry.ActionID,
		logEntry.PerformedByNpk,
		logEntry.DetailsText,
		logEntry.FilePath,
		logEntry.FromStatusID,
		logEntry.ToStatusID,
	)
	return err
}
