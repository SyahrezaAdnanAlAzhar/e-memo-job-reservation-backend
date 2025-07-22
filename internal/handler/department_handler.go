package handler

import (
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"net/http"
	"database/sql"
	"strconv"
	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	service *service.DepartmentService
}

func NewDepartmentHandler(service *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

// POST /department
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newDept, err := h.service.CreateDepartment(req)
	if err != nil {
		if err.Error() == "department name already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()}) 
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create department"})
		return
	}

	c.JSON(http.StatusCreated, newDept)
}

// GET /department
func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
	filters := make(map[string]string)

	if isActive, exists := c.GetQuery("is_active"); exists {
		filters["is_active"] = isActive
	}
	if receiveJob, exists := c.GetQuery("receive_job"); exists {
		filters["receive_job"] = receiveJob
	}

	departments, err := h.service.GetAllDepartments(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve departments"})
		return
	}
	
	if departments == nil {
		c.JSON(http.StatusOK, []model.Department{})
        return
	}

	c.JSON(http.StatusOK, departments)
}

// GET /department/:id
func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	department, err := h.service.GetDepartmentByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve department"})
		return
	}

	c.JSON(http.StatusOK, department)
}

// DELETE /department/:id
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	err = h.service.DeleteDepartment(id)
	if err != nil {
		if err.Error() == "department not found or already deleted" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete department"})
		return
	}

	c.Status(http.StatusNoContent)
}

// PUT /department/:id
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedDept, err := h.service.UpdateDepartment(id, req)
	if err != nil {
		if err.Error() == "department name already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update department"})
		return
	}

	c.JSON(http.StatusOK, updatedDept)
}

// PATCH /department/:id/status
func (h *DepartmentHandler) UpdateDepartmentActiveStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateDepartmentActiveStatus(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update department status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Department status updated successfully"})
}