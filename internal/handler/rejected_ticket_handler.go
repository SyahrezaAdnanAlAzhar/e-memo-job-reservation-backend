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

// CREATE
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

// UPDATE FEEDBACK
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