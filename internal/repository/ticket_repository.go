package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
)

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

// SCAN TO MAP
func scanToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		rowData := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowData[col] = string(b)
			} else {
				rowData[col] = val
			}
		}
		results = append(results, rowData)
	}
	return results, nil
}


func toNullInt64(val *int) sql.NullInt64 {
	if val == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*val), Valid: true}
}


// MAIN

// CREATE TICKET
func (r *TicketRepository) Create(ctx context.Context, tx *sql.Tx, ticket model.Ticket) (*model.Ticket, error) {
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

	var newTicket model.Ticket = ticket
	err := row.Scan(&newTicket.ID, &newTicket.CreatedAt, &newTicket.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &newTicket, nil
}

// GET ALL
func (r *TicketRepository) FindAll(filters map[string]string) ([]map[string]interface{}, error) {
	baseQuery := "SELECT * FROM view_ticket_list"
	var conditions []string
	var args []interface{}
	argID := 1

	for key, val := range filters {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", key, argID))
		args = append(args, val)
		argID++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}
	baseQuery += " ORDER BY ticket_priority ASC"

	rows, err := r.DB.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanToMap(rows)
}

// GET BY ID
func (r *TicketRepository) FindByID(id int) (map[string]interface{}, error) {
	query := "SELECT * FROM view_ticket_list WHERE ticket_id = $1"
	rows, err := r.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := scanToMap(rows)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}
	return results[0], nil
}


// GET BY ID AS STRUCT
func (r *TicketRepository) FindByIDAsStruct(ctx context.Context, id int) (*model.Ticket, error) {
	query := "SELECT id, requestor_npk, department_target_id, physical_location_id, specified_location_id, description, ticket_priority FROM ticket WHERE id = $1"
	row := r.DB.QueryRowContext(ctx, query, id)

	var t model.Ticket
	err := row.Scan(&t.ID, &t.Requestor, &t.DepartmentTargetID, &t.PhysicalLocationID, &t.SpecifiedLocationID, &t.Description, &t.TicketPriority)
	if err != nil {
		return nil, err
	}
	return &t, nil
}


// UPDATE TICKET
func (r *TicketRepository) Update(ctx context.Context, tx *sql.Tx, id int, req dto.UpdateTicketRequest) error {
	query := `
        UPDATE ticket 
        SET department_target_id = $1, description = $2, physical_location_id = $3, specified_location_id = $4, updated_at = NOW()
        WHERE id = $5`

	_, err := tx.ExecContext(ctx, query,
		req.DepartmentTargetID,
		req.Description,
		toNullInt64(req.PhysicalLocationID),
		toNullInt64(req.SpecifiedLocationID),
		id)
	return err
}


// REORDER
func (r *TicketRepository) UpdatePriority(ctx context.Context, tx *sql.Tx, ticketID int, version int, newPriority int) (int64, error) {
	query := `
        UPDATE ticket 
        SET ticket_priority = $1, version = version + 1, updated_at = NOW()
        WHERE id = $2 AND version = $3`
	
	result, err := tx.ExecContext(ctx, query, newPriority, ticketID, version)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UPDATE TICKET TO FALLBACK STATUS
func (r *TicketRepository) MoveTicketsToFallbackStatus(ctx context.Context, tx *sql.Tx, sectionIDToDeactivate int, fallbackStatusID int) error {
	findTicketsQuery := `
        SELECT tst.ticket_id
        FROM track_status_ticket tst
        WHERE tst.finish_date IS NULL
        AND tst.status_ticket_id IN (
            SELECT id FROM status_ticket WHERE section_id = $1
        )`
	rows, err := tx.QueryContext(ctx, findTicketsQuery, sectionIDToDeactivate)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ticketIDsToMove []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ticketIDsToMove = append(ticketIDsToMove, id)
	}

	if len(ticketIDsToMove) == 0 {
		return nil
	}

	deleteQuery := `
        DELETE FROM track_status_ticket
        WHERE ticket_id = ANY($1)
        AND status_ticket_id IN (
            SELECT id FROM status_ticket WHERE section_id = $2
        )`
	_, err = tx.ExecContext(ctx, deleteQuery, ticketIDsToMove, sectionIDToDeactivate)
	if err != nil {
		return err
	}

	createQuery := "INSERT INTO track_status_ticket (ticket_id, status_ticket_id, start_date) VALUES ($1, $2, NOW())"
	stmt, err := tx.PrepareContext(ctx, createQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, ticketID := range ticketIDsToMove {
		_, err := stmt.ExecContext(ctx, ticketID, fallbackStatusID)
		if err != nil {
			return err 
		}
	}

	return nil
}