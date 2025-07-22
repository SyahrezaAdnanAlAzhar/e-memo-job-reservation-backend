package handler

import (
	"database/sql"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StatusTicketHandler struct {
	service *service.StatusTicketService
}

func NewStatusTicketHandler(service *service.StatusTicketService) *StatusTicketHandler {
	return &StatusTicketHandler{service: service}
}


// POST /status-ticket
func (h *StatusTicketHandler) CreateStatusTicket(c *gin.Context) {
	var req dto.CreateStatusTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newStatus, err := h.service.CreateStatusTicket(req)
	if err != nil {
		if err.Error() == "status ticket name or sequence already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create status ticket"})
		return
	}

	c.JSON(http.StatusCreated, newStatus)
}

// GET /status-ticket
func (h *StatusTicketHandler) GetAllStatusTickets(c *gin.Context) {
	filters := make(map[string]string)
	if isActive, exists := c.GetQuery("is_active"); exists {
		filters["is_active"] = isActive
	}

	statuses, err := h.service.GetAllStatusTickets(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve status tickets"})
		return
	}

	if statuses == nil {
		c.JSON(http.StatusOK, []model.StatusTicket{})
		return
	}
	c.JSON(http.StatusOK, statuses)
}

// GET /status-ticket/:id
func (h *StatusTicketHandler) GetStatusTicketByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status ticket ID format"})
		return
	}

	status, err := h.service.GetStatusTicketByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Status ticket not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve status ticket"})
		return
	}

	c.JSON(http.StatusOK, status)
}

// DELETE /status-ticket/:id
func (h *StatusTicketHandler) DeleteStatusTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status ticket ID format"})
		return
	}
	if err := h.service.DeleteStatusTicket(id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Status ticket not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete status ticket"})
		return
	}
	c.Status(http.StatusNoContent)
}

// PATCH /status-ticket/:id/status
func (h *StatusTicketHandler) UpdateStatusTicketActiveStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status ticket ID format"})
		return
	}

	var req dto.UpdateStatusTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.service.UpdateStatusTicketActiveStatus(id, req); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Status ticket not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status ticket"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status ticket status updated successfully"})
}

// 
func (h *StatusTicketHandler) ReorderStatusTickets(c *gin.Context) {
	var req dto.ReorderStatusTicketsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.service.ReorderStatusTickets(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder status tickets", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status tickets reordered successfully"})
}