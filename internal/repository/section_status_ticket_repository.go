package repository

import (
	"database/sql"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
)

type SectionStatusTicketRepository struct {
	DB *sql.DB
}

func NewSectionStatusTicketRepository(db *sql.DB) *SectionStatusTicketRepository {
	return &SectionStatusTicketRepository{DB: db}
}

// CREATE
func (r *SectionStatusTicketRepository) Create(req dto.CreateSectionStatusTicketRequest) (*model.SectionStatusTicket, error) {
	query := `
        INSERT INTO section_status_ticket (name, sequence, is_active) 
        VALUES ($1, $2, false)
        RETURNING id, name, sequence, is_active, created_at, updated_at`

	row := r.DB.QueryRow(query, req.Name, req.Sequence)

	var newSection model.SectionStatusTicket
	err := row.Scan(
		&newSection.ID, &newSection.Name, &newSection.Sequence, &newSection.IsActive,
		&newSection.CreatedAt, &newSection.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &newSection, nil
}