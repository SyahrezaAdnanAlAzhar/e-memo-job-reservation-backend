package repository

import (
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