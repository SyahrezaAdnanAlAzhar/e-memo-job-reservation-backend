package repository

import (
	"context"
	"database/sql"
)

type WorkflowStepRepository struct {
	DB *sql.DB
}

func NewWorkflowStepRepository(db *sql.DB) *WorkflowStepRepository {
	return &WorkflowStepRepository{DB: db}
}

// CREATE
func (r *WorkflowStepRepository) Create(ctx context.Context, tx *sql.Tx, workflowID, statusTicketID, stepSequence int) error {
	query := "INSERT INTO workflow_step (workflow_id, status_ticket_id, step_sequence, is_active) VALUES ($1, $2, $3, true)"
	_, err := tx.ExecContext(ctx, query, workflowID, statusTicketID, stepSequence)
	return err
}