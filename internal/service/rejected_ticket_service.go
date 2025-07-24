package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type RejectedTicketService struct {
	repo                  *repository.RejectedTicketRepository
	ticketRepo            *repository.TicketRepository
	trackStatusTicketRepo *repository.TrackStatusTicketRepository
	statusTicketRepo      *repository.StatusTicketRepository
	employeeRepo          *repository.EmployeeRepository
	db                    *sql.DB
}

// CREATE
func (s *RejectedTicketService) CreateRejectedTicket(ctx context.Context, req dto.CreateRejectedTicketRequest, userNPK string) (*model.RejectedTicket, error) {
	latestRejection, err := s.repo.FindLatestByTicketID(ctx, req.TicketID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if latestRejection != nil {
		if !latestRejection.AlreadySeen {
			return nil, errors.New("ticket already has an active rejection that has not been seen")
		}
		_, currentStatusName, err := s.trackStatusTicketRepo.GetCurrentStatusByTicketID(ctx, int(req.TicketID))
		if err != nil {
			return nil, err
		}
		if currentStatusName == "Ditolak" {
			return nil, errors.New("ticket is still in 'Ditolak' status from a previous rejection")
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	newRejection, err := s.repo.Create(ctx, tx, req, userNPK)
	if err != nil {
		return nil, err
	}

	rejectedStatus, err := s.statusTicketRepo.FindByName("Ditolak")
	if err != nil {
		return nil, errors.New("critical configuration error: 'Ditolak' status not found")
	}
	if err := s.trackStatusTicketRepo.UpdateStatus(ctx, tx, int(req.TicketID), rejectedStatus.ID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return newRejection, nil
}

// UPDATE FEEDBACK
func (s *RejectedTicketService) UpdateFeedback(ctx context.Context, rejectionID int64, req dto.UpdateFeedbackRequest, userNPK string) (*model.RejectedTicket, error) {
	user, err := s.employeeRepo.FindByNPK(userNPK)
	if err != nil {
		return nil, errors.New("user not found")
	}

	rejection, err := s.repo.FindByID(rejectionID)
	if err != nil {
		return nil, errors.New("rejection record not found")
	}

	ticket, err := s.ticketRepo.FindByIDAsStruct(ctx, int(rejection.TicketID))
	if err != nil {
		return nil, errors.New("associated ticket not found")
	}

	isAllowed := user.DepartmentID == ticket.DepartmentTargetID && (user.Position.Name == "Head of Department" || user.Position.Name == "Section")
	if !isAllowed {
		return nil, errors.New("user is not authorized to update this feedback")
	}

	return s.repo.UpdateFeedback(rejectionID, req.Feedback)
}
