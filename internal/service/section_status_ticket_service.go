package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type SectionStatusTicketService struct {
	repo             *repository.SectionStatusTicketRepository
	statusTicketRepo *repository.StatusTicketRepository
	ticketRepo       *repository.TicketRepository
	db               *sql.DB
}

func NewSectionStatusTicketService(
	repo *repository.SectionStatusTicketRepository,
	statusTicketRepo *repository.StatusTicketRepository,
	ticketRepo *repository.TicketRepository, db *sql.DB) *SectionStatusTicketService {
	return &SectionStatusTicketService{repo: repo, statusTicketRepo: statusTicketRepo, ticketRepo: ticketRepo, db: db}
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

// GET ALL
func (s *SectionStatusTicketService) GetAllSectionStatusTickets() ([]model.SectionStatusTicket, error) {
	return s.repo.FindAll()
}

// GET BY ID
func (s *SectionStatusTicketService) GetSectionStatusTicketByID(id int) (*model.SectionStatusTicket, error) {
	return s.repo.FindByID(id)
}

// CHANGE STATUS
func (s *SectionStatusTicketService) UpdateSectionStatusTicketActiveStatus(ctx context.Context, id int, req dto.UpdateSectionStatusTicketStatusRequest) error {
	section, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if !req.IsActive { 
		activeCount, err := s.repo.CountActiveSections()
		if err != nil { return err }
		if activeCount <= 2 {
			return errors.New("cannot deactivate, must have at least two active sections")
		}
		if section.Sequence == 1 {
			return errors.New("cannot deactivate the first section")
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil { return err }
	defer tx.Rollback()

	if err := s.repo.UpdateActiveStatus(ctx, tx, id, req.IsActive); err != nil {
		return err
	}
	if err := s.statusTicketRepo.UpdateActiveStatusBySectionID(ctx, tx, id, req.IsActive); err != nil {
		return err
	}

	if !req.IsActive {
		fallbackStatusID, err := s.statusTicketRepo.GetDynamicFallbackStatusID(ctx, tx, section.Sequence)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("could not find a valid active fallback status")
			}
			return err
		}
		
		if err := s.ticketRepo.MoveTicketsToFallbackStatus(ctx, tx, id, fallbackStatusID); err != nil {
			return err
		}
	}

	return tx.Commit()
}