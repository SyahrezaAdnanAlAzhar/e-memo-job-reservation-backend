package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"

	"github.com/gin-gonic/gin"
)

type WorkflowHandler struct {
	service *service.WorkflowService
}

func NewWorkflowHandler(service *service.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{service: service}
}

// POST /workflow
func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	var req dto.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newWorkflow, err := h.service.CreateWorkflowWithSteps(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "workflow name already exists" || err.Error() == "cannot add the same status twice to a workflow" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "one or more status_ticket_ids are invalid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workflow"})
		return
	}
	c.JSON(http.StatusCreated, newWorkflow)
}

//POST /workflow/step
func (h *WorkflowHandler) AddWorkflowStep(c *gin.Context) {
	var req dto.AddWorkflowStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.service.AddWorkflowStep(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "status ticket is already in this workflow" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add workflow step", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Workflow step added successfully"})
}

// GET ALL
func (h *WorkflowHandler) GetAllWorkflows(c *gin.Context) {
	workflows, err := h.service.GetAllWorkflows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflows"})
		return
	}
	c.JSON(http.StatusOK, workflows)
}

func (h *WorkflowHandler) GetAllWorkflowSteps(c *gin.Context) {
	workflowIDStr := c.Query("workflow_id")
	if workflowIDStr != "" {
		workflowID, err := strconv.Atoi(workflowIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow_id format"})
			return
		}
		steps, err := h.service.GetWorkflowStepsByWorkflowID(workflowID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflow steps"})
			return
		}
		c.JSON(http.StatusOK, steps)
		return
	}

	steps, err := h.service.GetAllWorkflowSteps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve all workflow steps"})
		return
	}
	c.JSON(http.StatusOK, steps)
}

// GET BY ID
func (h *WorkflowHandler) GetWorkflowByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	workflow, err := h.service.GetWorkflowByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflow"})
		return
	}
	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) GetWorkflowStepByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	step, err := h.service.GetWorkflowStepByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workflow step not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflow step"})
		return
	}
	c.JSON(http.StatusOK, step)
}