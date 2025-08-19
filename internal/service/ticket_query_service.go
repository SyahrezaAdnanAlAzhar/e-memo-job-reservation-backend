package service

import (
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type TicketQueryService struct {
	ticketRepo *repository.TicketRepository
}

func NewTicketQueryService(ticketRepo *repository.TicketRepository) *TicketQueryService {
	return &TicketQueryService{
		ticketRepo: ticketRepo,
	}
}

// GET ALL
func (s *TicketQueryService) GetAllTickets(filters dto.TicketFilter) ([]dto.TicketDetailResponse, error) {
	return s.ticketRepo.FindAll(filters) 
}

// GET BY ID
func (s *TicketQueryService) GetTicketByID(id int) (*dto.TicketDetailResponse, error) {
	return s.ticketRepo.FindByID(id)
}

func (s *TicketQueryService) GetTicketSummary(filters dto.TicketSummaryFilter) ([]dto.TicketSummaryResponse, error) {
	return s.ticketRepo.GetTicketSummary(filters)
}
