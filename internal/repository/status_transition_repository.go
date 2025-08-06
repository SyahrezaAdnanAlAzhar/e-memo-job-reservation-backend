package repository

import (
	"database/sql"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
)

type StatusTransitionRepository struct {
	DB *sql.DB
}

func NewStatusTransitionRepository(db *sql.DB) *StatusTransitionRepository {
	return &StatusTransitionRepository{DB: db}
}

type TransitionDetail struct {
	RequiredActorRole string
	ActionDetail      dto.ActionResponse
}

func (r *StatusTransitionRepository) FindValidTransition(fromStatusID int, actionName string) (*model.StatusTransition, error) {
	query := `
        SELECT 
            st.id, st.from_status_id, st.to_status_id, st.action_id, 
            st.actor_role_id, st.require_reason, st.reason_label, st.require_file
        FROM status_transition st
        JOIN action a ON st.action_id = a.id
        WHERE st.from_status_id = $1 AND a.name = $2 AND st.is_active = true`

	row := r.DB.QueryRow(query, fromStatusID, actionName)

	var transition model.StatusTransition
	err := row.Scan(
		&transition.ID, &transition.FromStatusID, &transition.ToStatusID, &transition.ActionID,
		&transition.ActorRoleID, &transition.RequiresReason, &transition.ReasonLabel, &transition.RequiresFile,
	)
	if err != nil {
		return nil, err
	}
	return &transition, nil
}

func (r *StatusTransitionRepository) FindPossibleTransitionsWithDetails(fromStatusID int) ([]TransitionDetail, error) {
	query := `
        SELECT 
            ar.name as required_actor_role,
            a.name as action_name,
            a.hex_code,
            st.require_reason,
            st.reason_label,
            st.require_file
        FROM status_transition st
        JOIN action a ON st.action_id = a.id
        JOIN actor_role ar ON st.actor_role_id = ar.id
        WHERE st.from_status_id = $1 AND st.is_active = true AND a.is_active = true`

	rows, err := r.DB.Query(query, fromStatusID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transitions []TransitionDetail
	for rows.Next() {
		var t TransitionDetail
		err := rows.Scan(
			&t.RequiredActorRole,
			&t.ActionDetail.ActionName,
			&t.ActionDetail.HexCode,
			&t.ActionDetail.RequiresReason,
			&t.ActionDetail.ReasonLabel,
			&t.ActionDetail.RequiresFile,
		)
		if err != nil {
			return nil, err
		}
		transitions = append(transitions, t)
	}
	return transitions, nil
}
