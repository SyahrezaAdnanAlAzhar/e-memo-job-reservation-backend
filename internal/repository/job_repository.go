package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/lib/pq"
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

// FindByID
func (r *JobRepository) FindByID(id int) (*model.Job, error) {
	query := `
        SELECT id, ticket_id, pic_job_npk, assigned_department_id, job_priority, report_file, created_at, updated_at 
        FROM job WHERE id = $1`

	row := r.DB.QueryRow(query, id)

	var j model.Job
	err := row.Scan(
		&j.ID, &j.TicketID, &j.PicJob, &j.AssignedDepartmentID,
		&j.JobPriority, &j.ReportFile, &j.CreatedAt, &j.UpdatedAt,
	)
	return &j, err
}

// AssignPIC
func (r *JobRepository) AssignPIC(id int, picNpk string) error {
	query := "UPDATE job SET pic_job_npk = $1, updated_at = NOW() WHERE id = $2"
	result, err := r.DB.Exec(query, picNpk, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// UpdatePriority
func (r *JobRepository) UpdatePriority(ctx context.Context, tx *sql.Tx, jobID int, version int, newPriority int) (int64, error) {
	query := `
        UPDATE job 
        SET job_priority = $1, version = version + 1, updated_at = NOW()
        WHERE id = $2 AND version = $3`

	result, err := tx.ExecContext(ctx, query, newPriority, jobID, version)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// CheckJobsInDepartment
func (r *JobRepository) CheckJobsInDepartment(jobIDs []int, departmentID int) (int, error) {
	if len(jobIDs) == 0 {
		return 0, nil
	}

	query := "SELECT COUNT(id) FROM job WHERE id = ANY($1) AND assigned_department_id = $2"

	var count int
	err := r.DB.QueryRow(query, pq.Array(jobIDs), departmentID).Scan(&count)
	return count, err
}
