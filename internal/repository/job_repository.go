package repository

import (
	"context"
	"database/sql"
	"errors"
)

type JobRepository struct {
	DB *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{DB: db}
}

// CREATE
func (r *JobRepository) Create(ctx context.Context, tx *sql.Tx, ticketID int, initialJobPriority int) error {
	query := "INSERT INTO job (ticket_id, job_priority) VALUES ($1, $2)"
	_, err := tx.ExecContext(ctx, query, ticketID, initialJobPriority)
	return err
}

// CHECK THE JOB ALREADY GET ASSIGN OR NOT
func (r *JobRepository) IsJobAssigned(ctx context.Context, ticketID int) (bool, error) {
	var isAssigned bool
	query := "SELECT (pic_job_npk IS NOT NULL) FROM job WHERE ticket_id = $1"

	err := r.DB.QueryRowContext(ctx, query, ticketID).Scan(&isAssigned)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return isAssigned, nil
}

// GET PIC
func (r *JobRepository) GetPicByTicketID(ctx context.Context, ticketID int) (string, error) {
	var picNpk sql.NullString
	query := "SELECT pic_job_npk FROM job WHERE ticket_id = $1"

	err := r.DB.QueryRowContext(ctx, query, ticketID).Scan(&picNpk)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("job not found for this ticket")
		}
		return "", err
	}

	if !picNpk.Valid {
		return "", nil
	}

	return picNpk.String, nil
}

// UPLOAD REPORT FILE
func (r *JobRepository) UpdateReportFile(ctx context.Context, tx *sql.Tx, ticketID int, reportFilePath string) error {
	query := "UPDATE job SET report_file = array_append(report_file, $1), updated_at = NOW() WHERE ticket_id = $2"
	_, err := tx.ExecContext(ctx, query, reportFilePath, ticketID)
	return err
}