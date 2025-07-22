package repository

import (
	"context"
	"database/sql"
)

type WorkflowRepository struct {
	DB *sql.DB
}

func NewWorkflowRepository(db *sql.DB) *WorkflowRepository {
	return &WorkflowRepository{DB: db}
}


// HELPER
// GET INITIAL STATUS BY RULE
func (r *WorkflowRepository) GetInitialStatusByPosition(ctx context.Context, positionID int) (int, error) {
	var statusID int
	query := `
        SELECT ws.status_ticket_id
        FROM workflow_step ws
        JOIN position_to_workflow_mapping ptwm ON ws.workflow_id = ptwm.workflow_id
        WHERE ptwm.position_id = $1 AND ws.step_sequence = 0
        LIMIT 1`
	
	err := r.DB.QueryRowContext(ctx, query, positionID).Scan(&statusID)
	if err != nil {
		return 0, err
	}
	return statusID, nil
}