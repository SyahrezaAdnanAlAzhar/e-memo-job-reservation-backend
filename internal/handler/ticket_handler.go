package handler

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/pkg/filehandler"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	queryService    *service.TicketQueryService
	commandService  *service.TicketCommandService
	workflowService *service.TicketWorkflowService
	priorityService *service.TicketPriorityService
	actionService   *service.TicketActionService
}

type TicketHandlerConfig struct {
	QueryService    *service.TicketQueryService
	CommandService  *service.TicketCommandService
	WorkflowService *service.TicketWorkflowService
	PriorityService *service.TicketPriorityService
	ActionService   *service.TicketActionService
}

func NewTicketHandler(cfg *TicketHandlerConfig) *TicketHandler {
	return &TicketHandler{
		queryService:    cfg.QueryService,
		commandService:  cfg.CommandService,
		workflowService: cfg.WorkflowService,
		priorityService: cfg.PriorityService,
		actionService:   cfg.ActionService,
	}
}

// POST /tickets
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

	var filePaths []string
	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["support_files"]
		if len(files) > 0 {
			savedPaths, saveErr := filehandler.SaveFiles(c, files)
			if saveErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded files"})
				return
			}
			filePaths = savedPaths
		}
	}

	createdTicket, err := h.commandService.CreateTicket(c.Request.Context(), req, requestorNPK, filePaths)
	if err != nil {
		for _, p := range filePaths {
			os.Remove(p)
		}
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
	var filters dto.TicketFilter
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}

	tickets, err := h.queryService.GetAllTickets(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tickets", "details": err.Error()})
		return
	}

	if tickets == nil {
		c.JSON(http.StatusOK, []dto.TicketDetailResponse{})
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
		case "data conflict: ticket has been modified by another user, please refresh":
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
	userNPK := c.GetString("user_npk")
	var req dto.ReorderTicketsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.priorityService.ReorderTickets(c.Request.Context(), req, userNPK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder tickets"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket priorities updated successfully"})
}

// POST /tickets/:id/action
func (h *TicketHandler) ExecuteAction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID format"})
		return
	}
	userNPK := c.GetString("user_npk")

	var req dto.ExecuteActionRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	var filePaths []string
	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["files"]
		if len(files) > 0 {
			savedPaths, saveErr := filehandler.SaveFiles(c, files)
			if saveErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded files"})
				return
			}
			filePaths = savedPaths
		}
	}

	err = h.workflowService.ExecuteAction(c, id, userNPK, req, filePaths)

	if err != nil {
		for _, p := range filePaths {
			os.Remove(p)
		}
		switch err.Error() {
		case "ticket not found", "user not found", "original requestor not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "user does not have the required role for this action":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "action not allowed from the current status", "reason is required for this action", "file upload is required for this action":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute action", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Action '" + req.ActionName + "' executed successfully"})
}

// GET /tickets/:id/available-actions
func (h *TicketHandler) GetAvailableActions(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userNPK := c.GetString("user_npk")

	actions, err := h.actionService.GetAvailableActions(c.Request.Context(), id, userNPK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get available actions", "details": err.Error()})
		return
	}

	if actions == nil {
		c.JSON(http.StatusOK, []dto.ActionResponse{})
		return
	}
	c.JSON(http.StatusOK, actions)
}
