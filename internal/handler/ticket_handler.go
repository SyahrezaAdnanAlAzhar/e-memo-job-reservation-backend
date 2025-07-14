package handler

import (
	"net/http"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	service *service.TicketService
}

func NewTicketHandler(service *service.TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

// POST /ticket
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var req repository.CreateTicketRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	requestorNPK := c.GetString("user_npk")
	if requestorNPK == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User NPK not found in token"})
		return
	}

	createdTicket, err := h.service.CreateTicket(c.Request.Context(), req, requestorNPK)
	if err != nil {
		switch err.Error() {
		case "requestor not found", "no workflow defined for this user's position":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, createdTicket)
}
