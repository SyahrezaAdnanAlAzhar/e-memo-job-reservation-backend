package service

import (
	"errors"

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
func (s *TicketQueryService) GetAllTickets(filters map[string]string) ([]map[string]interface{}, error) {
	allowedFilters := map[string]bool{
		"department_target_id": true,
		"current_status":       true,
		"requestor_npk":        true,
	}
	for key := range filters {
		if !allowedFilters[key] {
			return nil, errors.New("invalid filter key: " + key)
		}
	}
	return s.ticketRepo.FindAll(filters)
}

// GET BY ID
func (s *TicketQueryService) GetTicketByID(id int) (map[string]interface{}, error) {
	return s.ticketRepo.FindByID(id)
}
