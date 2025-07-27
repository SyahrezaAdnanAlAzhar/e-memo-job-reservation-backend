package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	queryService    *service.TicketQueryService
	commandService  *service.TicketCommandService
	workflowService *service.TicketWorkflowService
	priorityService *service.TicketPriorityService
}

type TicketHandlerConfig struct {
	QueryService    *service.TicketQueryService
	CommandService  *service.TicketCommandService
	WorkflowService *service.TicketWorkflowService
	PriorityService *service.TicketPriorityService
}

func NewTicketHandler(cfg *TicketHandlerConfig) *TicketHandler {
	return &TicketHandler{
		queryService:    cfg.QueryService,
		commandService:  cfg.CommandService,
		workflowService: cfg.WorkflowService,
		priorityService: cfg.PriorityService,
	}
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

	createdTicket, err := h.commandService.CreateTicket(c.Request.Context(), req, requestorNPK)
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

	tickets, err := h.queryService.GetAllTickets(filters)
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
	ticket, err := h.queryService.GetTicketByID(id)
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
	err := h.commandService.UpdateTicket(c.Request.Context(), id, req, userNPK)
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
	err := h.priorityService.ReorderTickets(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder tickets"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket priorities updated successfully"})
}
