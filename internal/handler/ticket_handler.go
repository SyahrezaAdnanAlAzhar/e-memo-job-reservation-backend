package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"

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
	var req dto.CreateTicketRequest

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

// GET ALL
func (h *TicketHandler) GetAllTickets(c *gin.Context) {
	filters := make(map[string]string)
	for key, value := range c.Request.URL.Query() {
		if len(value) > 0 {
			filters[key] = value[0]
		}
	}

	tickets, err := h.service.GetAllTickets(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tickets == nil {
		c.JSON(http.StatusOK, []map[string]interface{}{})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

// GET BY ID
func (h *TicketHandler) GetTicketByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ticket, err := h.service.GetTicketByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ticket"})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

// PUT UPDATE
func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userNPK := c.GetString("user_npk")
	err := h.service.UpdateTicket(c.Request.Context(), id, req, userNPK)
	if err != nil {
		switch err.Error() {
		case "ticket not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "user is not authorized to edit this ticket":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "ticket cannot be edited in its current state":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket updated and resubmitted for approval"})
}

// PUT REORDER
func (h *TicketHandler) ReorderTickets(c *gin.Context) {
	var req dto.ReorderTicketsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.ReorderTickets(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder tickets"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket priorities updated successfully"})
}

// POST /ticket/:id/reject
func (h *TicketHandler) RejectTicket(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userNPK := c.GetString("user_npk")

	var req dto.RejectTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reason for rejection is required"})
		return
	}

	err := h.service.RejectTicket(c.Request.Context(), id, req, userNPK)
	if err != nil {
		if err.Error() == "user not authorized to reject this ticket" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "ticket not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject ticket", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket has been rejected"})
}

// POST /ticket/:id/cancel
func (h *TicketHandler) CancelTicket(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userNPK := c.GetString("user_npk")

	err := h.service.CancelTicket(c.Request.Context(), id, userNPK)
	if err != nil {
		if err.Error() == "user not authorized to cancel this ticket" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "ticket not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel ticket", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket has been cancelled"})
}