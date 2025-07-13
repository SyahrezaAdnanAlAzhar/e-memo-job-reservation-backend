package service

import (
	"errors"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository" 

	"github.com/jackc/pgx/v5/pgconn"
)

type StatusTicketService struct {
	repo *repository.StatusTicketRepository
}

func NewStatusTicketService(repo *repository.StatusTicketRepository) *StatusTicketService {
	return &StatusTicketService{repo: repo}
}

// CREATE
func (s *StatusTicketService) CreateStatusTicket(req repository.CreateStatusTicketRequest) (*repository.StatusTicket, error) {
	newStatus, err := s.repo.Create(req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errors.New("status ticket name or sequence already exists")
		}
		return nil, err
	}
	return newStatus, nil
}

// GET ALL
func (s *StatusTicketService) GetAllStatusTickets(filters map[string]string) ([]repository.StatusTicket, error) {
	return s.repo.FindAll(filters)
}

// GET BY ID
func (s *StatusTicketService) GetStatusTicketByID(id int) (*repository.StatusTicket, error) {
	return s.repo.FindByID(id)
}