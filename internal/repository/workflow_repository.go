package repository

import (
	"context"
	"database/sql"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
)

type WorkflowRepository struct {
	DB *sql.DB
}

func NewWorkflowRepository(db *sql.DB) *WorkflowRepository {
	return &WorkflowRepository{DB: db}
}

// CREATE
func (r *WorkflowRepository) Create(ctx context.Context, tx *sql.Tx, name string) (*model.Workflow, error) {
	query := `
        INSERT INTO workflow (name, is_active) VALUES ($1, true)
        RETURNING id, name, is_active, created_at, updated_at`
	
	row := tx.QueryRowContext(ctx, query, name)

	var newWorkflow model.Workflow
	err := row.Scan(
		&newWorkflow.ID, &newWorkflow.Name, &newWorkflow.IsActive,
		&newWorkflow.CreatedAt, &newWorkflow.UpdatedAt,
	)
	return &newWorkflow, err
}

// HELPER
// GET INITIAL STATUS BY RULE
func (r *WorkflowRepository) GetInitialStatusByPosition(ctx context.Context, positionID int) (int, error) {
	var statusID int
	query := `
        SELECT ws.status_ticket_id
        FROM workflow_step ws
        JOIN position_to_workflow_mapping ptwm ON ws.workflow_id = ptwm.workflow_id
        WHERE ptwm.employee_position_id = $1 AND ws.step_sequence = 0
        LIMIT 1`
	
	err := r.DB.QueryRowContext(ctx, query, positionID).Scan(&statusID)
	if err != nil {
		return 0, err
	}
	return statusID, nil
}