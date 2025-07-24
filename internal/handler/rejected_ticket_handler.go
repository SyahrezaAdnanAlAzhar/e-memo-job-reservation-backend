package handler

import (
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RejectedTicketHandler struct {
	service *service.RejectedTicketService
}

func NewRejectedTicketHandler(service *service.RejectedTicketService) *RejectedTicketHandler {
	return &RejectedTicketHandler{service: service}
}

// POST /rejected-tickets/
func (h *RejectedTicketHandler) CreateRejectedTicket(c *gin.Context) {
	var req dto.CreateRejectedTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userNPK := c.GetString("user_npk")

	newRejection, err := h.service.CreateRejectedTicket(c.Request.Context(), req, userNPK)
	if err != nil {
		if err.Error() == "ticket already has an active rejection that has not been seen" || err.Error() == "ticket is still in 'Ditolak' status from a previous rejection" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rejection", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newRejection)
}

// PUT /rejected-tickets/:id/feedback
func (h *RejectedTicketHandler) UpdateFeedback(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userNPK := c.GetString("user_npk")

	var req dto.UpdateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedRejection, err := h.service.UpdateFeedback(c.Request.Context(), id, req, userNPK)
	if err != nil {
		if err.Error() == "user is not authorized to update this feedback" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "rejection record not found" || err.Error() == "associated ticket not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update feedback", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedRejection)
}

// PATCH /rejected-ticket/:id/seen
func (h *RejectedTicketHandler) UpdateAlreadySeen(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userNPK := c.GetString("user_npk")

	var req dto.UpdateAlreadySeenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateAlreadySeen(c.Request.Context(), id, req, userNPK)
	if err != nil {
		if err.Error() == "user is not authorized to perform this action" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "rejection record not found" || err.Error() == "associated ticket not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update 'already_seen' status", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "'already_seen' status updated successfully"})
}

// DELETE /rejected-ticket/:id
func (h *RejectedTicketHandler) DeleteRejectedTicket(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userNPK := c.GetString("user_npk")

	err := h.service.DeleteRejectedTicket(c.Request.Context(), id, userNPK)
	if err != nil {
		if err.Error() == "user is not authorized to delete this rejection record" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "can only delete rejection record if ticket status is 'Ditolak'" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "rejection record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rejection record", "details": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}