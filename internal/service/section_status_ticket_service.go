package service

import (
	"errors"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type SectionStatusTicketService struct {
	repo *repository.SectionStatusTicketRepository
}

func NewSectionStatusTicketService(repo *repository.SectionStatusTicketRepository) *SectionStatusTicketService {
	return &SectionStatusTicketService{repo: repo}
}

// CREATE
func (s *SectionStatusTicketService) CreateSectionStatusTicket(req dto.CreateSectionStatusTicketRequest) (*model.SectionStatusTicket, error) {
	newSection, err := s.repo.Create(req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 = unique_violation
			return nil, errors.New("section name or sequence already exists")
		}
		return nil, err
	}
	return newSection, nil
}