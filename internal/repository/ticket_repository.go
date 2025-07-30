package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/lib/pq"
)

type TicketRepository struct {
	DB *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{DB: db}
}

const baseTicketQuery = `
    SELECT
        t.id,
        t.description,
        t.department_target_id,
        dt.name as department_target_name,
        t.ticket_priority,
        t.version,
        j.id as job_id,
        j.job_priority,
        pl.name as location_name,
        sl.name as specified_location_name,
        t.created_at,
        (NOW()::date - t.created_at::date) as ticket_age_days,
        t.deadline,
        (t.deadline::date - NOW()::date) as days_remaining,
        req_emp.name as requestor_name,
        req_dept.name as requestor_department,
        pic_emp.name as pic_name,
        pic_area.name as pic_area_name,
        current_st.name as current_status,
        current_st.hex_color as current_status_hex_code,
        current_sst.name as current_section_name
    FROM ticket t
    LEFT JOIN job j ON t.id = j.ticket_id
    LEFT JOIN department dt ON t.department_target_id = dt.id
    LEFT JOIN physical_location pl ON t.physical_location_id = pl.id
    LEFT JOIN specified_location sl ON t.specified_location_id = sl.id
    JOIN employee req_emp ON t.requestor = req_emp.npk
    LEFT JOIN department req_dept ON req_emp.department_id = req_dept.id
    LEFT JOIN employee pic_emp ON j.pic_job = pic_emp.npk
    LEFT JOIN area pic_area ON pic_emp.area_id = pic_area.id
    LEFT JOIN (
        SELECT DISTINCT ON (ticket_id) ticket_id, status_ticket_id
        FROM track_status_ticket
        ORDER BY ticket_id, start_date DESC
    ) current_tst ON t.id = current_tst.ticket_id
    LEFT JOIN status_ticket current_st ON current_tst.status_ticket_id = current_st.id
    LEFT JOIN section_status_ticket current_sst ON current_st.section_id = current_sst.id
`

// MAIN

// CREATE TICKET
func (r *TicketRepository) Create(ctx context.Context, tx *sql.Tx, ticket model.Ticket) (*model.Ticket, error) {
	query := `
        INSERT INTO ticket (requestor, department_target_id, physical_location_id, specified_location_id, description, ticket_priority, deadline)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at`

	row := tx.QueryRowContext(ctx, query,
		ticket.Requestor,
		ticket.DepartmentTargetID,
		ticket.PhysicalLocationID,
		ticket.SpecifiedLocationID,
		ticket.Description,
		ticket.TicketPriority,
		ticket.Deadline,
	)

	var newTicket model.Ticket = ticket
	err := row.Scan(&newTicket.ID, &newTicket.CreatedAt, &newTicket.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &newTicket, nil
}

// GET ALL
func (r *TicketRepository) FindAll() ([]dto.TicketDetailResponse, error) {
	// Untuk saat ini, kita hanya menjalankan query dasar tanpa filter.
	// Nanti kita akan tambahkan logika WHERE dan ORDER BY di sini.
	query := baseTicketQuery + " ORDER BY t.ticket_priority ASC"

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTicketDetails(rows)
}

// GET BY ID
func (r *TicketRepository) FindByID(id int) (*dto.TicketDetailResponse, error) {
	query := baseTicketQuery + " WHERE t.id = $1"
	rows, err := r.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets, err := scanTicketDetails(rows)
	if err != nil {
		return nil, err
	}
	if len(tickets) == 0 {
		return nil, sql.ErrNoRows
	}
	return &tickets[0], nil
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
        SET department_target_id = $1, description = $2, physical_location_id = $3, specified_location_id = $4, deadline = $5, updated_at = NOW()
        WHERE id = $5`

	deadline, _ := ParseDeadline(req.Deadline)

	_, err := tx.ExecContext(ctx, query,
		req.DepartmentTargetID,
		req.Description,
		toNullInt64(req.PhysicalLocationID),
		toNullInt64(req.SpecifiedLocationID),
		deadline,
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

// FORCE REORDER
func (r *TicketRepository) ForceUpdatePriority(ctx context.Context, tx *sql.Tx, ticketID int, newPriority int) error {
	query := `
        UPDATE ticket 
        SET ticket_priority = $1, version = version + 1, updated_at = NOW()
        WHERE id = $2`

	_, err := tx.ExecContext(ctx, query, newPriority, ticketID)
	return err
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

func (r *TicketRepository) CheckTicketsFromDepartment(ticketIDs []int, requestorDepartmentID int) (int, error) {
	if len(ticketIDs) == 0 {
		return 0, nil
	}

	query := `
        SELECT COUNT(t.id) 
        FROM ticket t
        JOIN employee e ON t.requestor_npk = e.npk
        WHERE t.id = ANY($1) AND e.department_id = $2`

	var count int
	err := r.DB.QueryRow(query, pq.Array(ticketIDs), requestorDepartmentID).Scan(&count)
	return count, err
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

// PARSE DEADLINE
func ParseDeadline(deadlineStr *string) (sql.NullTime, error) {
	if deadlineStr == nil {
		return sql.NullTime{Valid: false}, nil
	}
	// Format "2006-01-02"
	t, err := time.Parse("2006-01-02", *deadlineStr)
	if err != nil {
		return sql.NullTime{Valid: false}, err
	}
	return sql.NullTime{Time: t, Valid: true}, nil
}

// QUERY MAPPING
func scanTicketDetails(rows *sql.Rows) ([]dto.TicketDetailResponse, error) {
	var tickets []dto.TicketDetailResponse
	for rows.Next() {
		var t dto.TicketDetailResponse
		err := rows.Scan(
			&t.TicketID,
			&t.Description,
			&t.DepartmentTargetID,
			&t.DepartmentTargetName,
			&t.TicketPriority,
			&t.Version,
			&t.JobID,
			&t.JobPriority,
			&t.LocationName,
			&t.SpecifiedLocationName,
			&t.CreatedAt,
			&t.TicketAgeDays,
			&t.Deadline,
			&t.DaysRemaining,
			&t.RequestorName,
			&t.RequestorDepartment,
			&t.PicName,
			&t.PicAreaName,
			&t.CurrentStatus,
			&t.CurrentStatusHexCode,
			&t.CurrentSectionName,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}
