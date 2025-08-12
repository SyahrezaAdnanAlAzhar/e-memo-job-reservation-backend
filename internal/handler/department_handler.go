package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/util"
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
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	newDept, err := h.service.CreateDepartment(req)
	if err != nil {
		if err.Error() == "department name already exists" {
			util.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to create department", nil)
		return
	}

	util.SuccessResponse(c, http.StatusCreated, newDept)
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
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve departments", nil)
		return
	}

	if departments == nil {
		util.SuccessResponse(c, http.StatusOK, []model.Department{})
		return
	}

	util.SuccessResponse(c, http.StatusOK, departments)
}

// GET /department/:id
func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid department ID", nil)
		return
	}

	department, err := h.service.GetDepartmentByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			util.ErrorResponse(c, http.StatusNotFound, "Department not found", nil)
			return
		}
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve department", nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, department)
}

// DELETE /department/:id
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid department ID", nil)
		return
	}

	err = h.service.DeleteDepartment(id)
	if err != nil {
		if err.Error() == "department not found or already deleted" {
			util.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete department", nil)
		return
	}

	c.Status(http.StatusNoContent)
}

// PUT /department/:id
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid department ID", nil)
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	updatedDept, err := h.service.UpdateDepartment(id, req)
	if err != nil {
		if err.Error() == "department name already exists" {
			util.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		if err == sql.ErrNoRows {
			util.ErrorResponse(c, http.StatusNotFound, "Department not found", nil)
			return
		}
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to update department", nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, updatedDept)
}

// PATCH /department/:id/status
func (h *DepartmentHandler) UpdateDepartmentActiveStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid department ID", nil)
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = h.service.UpdateDepartmentActiveStatus(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			util.ErrorResponse(c, http.StatusNotFound, "Department not found", nil)
			return
		}
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to update department status", nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, gin.H{"message": "Department status updated successfully"})
}
