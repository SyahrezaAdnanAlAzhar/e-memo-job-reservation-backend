package repository

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type StatusTicket struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Sequence  int       `json:"sequence"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateStatusTicketRequest struct {
	Name     string `json:"name" binding:"required"`
	Sequence int    `json:"sequence"`
}

type UpdateStatusTicketStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type ReorderStatusTicketsRequest struct {
	DeleteSectionOrder   []int `json:"delete_section_order"`
	ApprovalSectionOrder []int `json:"approval_section_order"`
	ActualSectionOrder   []int `json:"actual_section_order"`
}

type StatusTicketRepository struct {
	DB *sql.DB
}

func NewStatusTicketRepository(db *sql.DB) *StatusTicketRepository {
	return &StatusTicketRepository{DB: db}
}

// MAIN

// CREATE
func (r *StatusTicketRepository) Create(req CreateStatusTicketRequest) (*StatusTicket, error) {
	query := `
        INSERT INTO status_ticket (name, sequence, is_active)
        VALUES ($1, $2, false)
        RETURNING id, name, sequence, is_active, created_at, updated_at`

	row := r.DB.QueryRow(query, req.Name, req.Sequence)

	var newStatus StatusTicket
	err := row.Scan(
		&newStatus.ID, &newStatus.Name, &newStatus.Sequence,
		&newStatus.IsActive, &newStatus.CreatedAt, &newStatus.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &newStatus, nil
}

// GET ALL
func (r *StatusTicketRepository) FindAll(filters map[string]string) ([]StatusTicket, error) {
	baseQuery := "SELECT id, name, sequence, is_active, created_at, updated_at FROM status_ticket"
	var conditions []string
	var args []interface{}
	argID := 1

	if val, ok := filters["is_active"]; ok {
		conditions = append(conditions, "is_active = $"+strconv.Itoa(argID))
		args = append(args, val)
		argID++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY sequence ASC"

	rows, err := r.DB.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []StatusTicket
	for rows.Next() {
		var s StatusTicket
		err := rows.Scan(&s.ID, &s.Name, &s.Sequence, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, nil
}

// GET BY ID
func (r *StatusTicketRepository) FindByID(id int) (*StatusTicket, error) {
	query := "SELECT id, name, sequence, is_active, created_at, updated_at FROM status_ticket WHERE id = $1"
	row := r.DB.QueryRow(query, id)

	var s StatusTicket
	err := row.Scan(&s.ID, &s.Name, &s.Sequence, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// DELETE
func (r *StatusTicketRepository) Delete(id int) error {
	query := "DELETE FROM status_ticket WHERE id = $1"
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// CHANGE ACTIVE STATUS
func (r *StatusTicketRepository) UpdateActiveStatus(id int, isActive bool) error {
	query := "UPDATE status_ticket SET is_active = $1, updated_at = NOW() WHERE id = $2"
	result, err := r.DB.Exec(query, isActive, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// REORDER
func (r *StatusTicketRepository) Reorder(ctx context.Context, tx *sql.Tx, id int, newSequence int) error {
	query := "UPDATE status_ticket SET sequence = $1, updated_at = NOW() WHERE id = $2"
	_, err := tx.ExecContext(ctx, query, newSequence, id)
	return err
}

// GET NEXT STATUS BASED ON SEQUENCE
func (r *StatusTicketRepository) GetNextStatusInSection(currentStatusID int) (*StatusTicket, error) {
	query := `
        WITH current_status AS (
            SELECT section_id, sequence
            FROM status_ticket
            WHERE id = $1
        )
        SELECT id, name, sequence, is_active
        FROM status_ticket
        WHERE section_id = (SELECT section_id FROM current_status)
          AND sequence > (SELECT sequence FROM current_status)
        ORDER BY sequence ASC
        LIMIT 1`

	var nextStatus StatusTicket
	err := r.DB.QueryRow(query, currentStatusID).Scan(&nextStatus.ID, &nextStatus.Name, &nextStatus.Sequence, &nextStatus.IsActive)
	if err != nil {
		return nil, err
	}
	return &nextStatus, nil
}
