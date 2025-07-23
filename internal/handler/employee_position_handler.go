package handler

import (
	"database/sql"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmployeePositionHandler struct {
	service *service.EmployeePositionService
}

func NewEmployeePositionHandler(service *service.EmployeePositionService) *EmployeePositionHandler {
	return &EmployeePositionHandler{service: service}
}

//POST /employee-position
func (h *EmployeePositionHandler) CreateEmployeePosition(c *gin.Context) {
	var req dto.CreateEmployeePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPos, err := h.service.CreateEmployeePosition(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "position name already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "invalid workflow_id" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create position"})
		return
	}
	c.JSON(http.StatusCreated, newPos)
}

// GET /api/v1/employee-position
func (h *EmployeePositionHandler) GetAllEmployeePositions(c *gin.Context) {
	positions, err := h.service.GetAllEmployeePositions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve positions"})
		return
	}
	if positions == nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}
	c.JSON(http.StatusOK, positions)
}

// GET /employee-position/:id
func (h *EmployeePositionHandler) GetEmployeePositionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	position, err := h.service.GetEmployeePositionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve position"})
		return
	}
	c.JSON(http.StatusOK, position)
}

// PUT /employee-position/:id
func (h *EmployeePositionHandler) UpdateEmployeePosition(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateEmployeePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPos, err := h.service.UpdateEmployeePosition(id, req)
	if err != nil {
		if err.Error() == "position name already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update position"})
		return
	}
	c.JSON(http.StatusOK, updatedPos)
}

// DELETE /employee-position/:id
func (h *EmployeePositionHandler) DeleteEmployeePosition(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteEmployeePosition(c.Request.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete position", "details": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// PATCH /employee-position/:id/status
func (h *EmployeePositionHandler) UpdateEmployeePositionActiveStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateEmployeePositionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateEmployeePositionActiveStatus(id, req); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Position status updated successfully"})
}