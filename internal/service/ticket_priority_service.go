package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/websocket"
	"github.com/gin-gonic/gin"
)

type TicketPriorityService struct {
	db           *sql.DB
	hub          *websocket.Hub
	ticketRepo   *repository.TicketRepository
	employeeRepo *repository.EmployeeRepository
}

func NewTicketPriorityService(db *sql.DB, hub *websocket.Hub, ticketRepo *repository.TicketRepository, employeeRepo *repository.EmployeeRepository) *TicketPriorityService {
	return &TicketPriorityService{
		db:           db,
		hub:          hub,
		ticketRepo:   ticketRepo,
		employeeRepo: employeeRepo,
	}
}

// RE ORDER
func (s *TicketPriorityService) ReorderTickets(ctx context.Context, req dto.ReorderTicketsRequest, userNPK string) error {
	// --- Validasi Kontekstual (Gerbang Kedua) ---
	// Logika ini memastikan seorang Head hanya bisa reorder tiket yang dibuat oleh departemennya.
	user, err := s.employeeRepo.FindByNPK(userNPK)
	if err != nil {
		return errors.New("action performer not found")
	}

	// Ekstrak ticket IDs dari DTO
	ticketIDs := make([]int, len(req.Items))
	for i, item := range req.Items {
		ticketIDs[i] = item.TicketID
	}

	validTicketCount, err := s.ticketRepo.CheckTicketsFromDepartment(ticketIDs, user.DepartmentID)
	if err != nil {
		return err
	}
	if validTicketCount != len(req.Items) {
		return errors.New("user can only reorder tickets requested by their own department")
	}

	// --- Eksekusi dalam Transaksi ---
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, item := range req.Items {
		newPriority := i + 1
		rowsAffected, err := s.ticketRepo.UpdatePriority(ctx, tx, item.TicketID, item.Version, newPriority)
		if err != nil {
			return err // Rollback akan terjadi
		}
		if rowsAffected == 0 {
			// Konflik versi, batalkan seluruh operasi
			return errors.New("data conflict: one or more tickets have been modified by another user, please refresh")
		}
	}

	// [FIX] Commit HANYA dilakukan di sini, setelah loop selesai.
	if err := tx.Commit(); err != nil {
		return err
	}

	// --- [BARU] Siarkan Pembaruan via WebSocket ---
	// Setelah commit berhasil, beritahu semua klien.
	payload := gin.H{
		"department_target_id": req.DepartmentTargetID,
		"message":              "Ticket priorities have been updated.",
	}
	message, err := websocket.NewMessage("TICKET_PRIORITY_UPDATED", payload)
	if err != nil {
		log.Printf("CRITICAL: Failed to create websocket message for ticket reorder: %v", err)
	} else {
		s.hub.BroadcastMessage(message)
	}

	return nil
}
