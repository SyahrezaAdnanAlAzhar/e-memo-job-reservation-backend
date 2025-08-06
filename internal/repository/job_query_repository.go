package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
)

const baseJobQuery = `
    SELECT
        j.id as job_id,
        t.id as ticket_id,
        t.description,
        j.job_priority,
        t.ticket_priority,
        j.version,
        t.department_target_id as assigned_department_id,
        dept.name as assigned_department_name,
        current_st.name as current_status,
        current_st.hex_color as current_status_hex_code,
        pic_emp.name as pic_name,
        req_emp.name as requestor_name,
        (NOW()::date - t.created_at::date) as ticket_age_days,
        t.deadline,
        (t.deadline::date - NOW()::date) as days_remaining
    FROM job j
    JOIN ticket t ON j.ticket_id = t.id
    JOIN department dept ON t.department_target_id = dept.id -- Join ke department via ticket
    JOIN employee req_emp ON t.requestor = req_emp.npk
    LEFT JOIN employee pic_emp ON j.pic_job = pic_emp.npk
    LEFT JOIN (
        SELECT DISTINCT ON (ticket_id) ticket_id, status_ticket_id
        FROM track_status_ticket
        ORDER BY ticket_id, start_date DESC
    ) current_tst ON t.id = current_tst.ticket_id
    LEFT JOIN status_ticket current_st ON current_tst.status_ticket_id = current_st.id
`

type JobQueryRepository struct {
	DB *sql.DB
}

func NewJobQueryRepository(db *sql.DB) *JobQueryRepository {
	return &JobQueryRepository{DB: db}
}

// GET ALL
func (r *JobQueryRepository) FindAll(filters dto.JobFilter) ([]dto.JobDetailResponse, error) {
	query := baseJobQuery
	var conditions []string
	var args []interface{}
	argID := 1

	if filters.AssignedDepartmentID != 0 {
		conditions = append(conditions, fmt.Sprintf("t.department_target_id = $%d", argID))
		args = append(args, filters.AssignedDepartmentID)
		argID++
	}
	if filters.StatusID != 0 {
		conditions = append(conditions, fmt.Sprintf("current_st.id = $%d", argID))
		args = append(args, filters.StatusID)
		argID++
	}
	if filters.PicNPK != "" {
		conditions = append(conditions, fmt.Sprintf("j.pic_job = $%d", argID))
		args = append(args, filters.PicNPK)
		argID++
	}

	if filters.SearchQuery != "" {
		searchQuery := strings.ReplaceAll(strings.TrimSpace(filters.SearchQuery), " ", " & ")
		conditions = append(conditions, fmt.Sprintf("t.description_tsv @@ to_tsquery('simple', $%d)", argID))
		args = append(args, searchQuery)
		argID++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	orderByClause := " ORDER BY j.job_priority ASC"
	if filters.SortBy != "" {
		allowedSortColumns := map[string]string{
			"priority":  "j.job_priority",
			"deadline":  "t.deadline",
			"age":       "ticket_age_days",
			"status":    "current_st.name",
			"requestor": "req_emp.name",
			"pic":       "pic_emp.name",
		}

		var sortClauses []string
		sortParams := strings.Split(filters.SortBy, ",")

		for _, param := range sortParams {
			parts := strings.Split(strings.TrimSpace(param), "_")
			if len(parts) != 2 {
				continue
			}

			columnKey := parts[0]
			direction := strings.ToUpper(parts[1])

			if dbColumn, ok := allowedSortColumns[columnKey]; ok && (direction == "ASC" || direction == "DESC") {
				sortClauses = append(sortClauses, dbColumn+" "+direction)
			}
		}

		if len(sortClauses) > 0 {
			orderByClause = " ORDER BY " + strings.Join(sortClauses, ", ")
		}
	}
	query += orderByClause

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanJobDetails(rows)
}

// GET BY ID
func (r *JobQueryRepository) FindByID(id int) (*dto.JobDetailResponse, error) {
	query := baseJobQuery + " WHERE j.id = $1"

	rows, err := r.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs, err := scanJobDetails(rows)
	if err != nil {
		return nil, err
	}

	if len(jobs) == 0 {
		return nil, sql.ErrNoRows
	}

	return &jobs[0], nil
}

func scanJobDetails(rows *sql.Rows) ([]dto.JobDetailResponse, error) {
	var jobs []dto.JobDetailResponse
	for rows.Next() {
		var j dto.JobDetailResponse
		err := rows.Scan(
			&j.JobID, &j.TicketID, &j.Description, &j.JobPriority, &j.TicketPriority,
			&j.Version, &j.AssignedDepartmentID, &j.AssignedDepartmentName,
			&j.CurrentStatus, &j.CurrentStatusHexCode, &j.PicName, &j.RequestorName,
			&j.TicketAgeDays, &j.Deadline, &j.DaysRemaining,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}
