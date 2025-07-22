package handler

import (
	"database/sql"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PositionPermissionHandler struct {
	service *service.PositionPermissionService
}

func NewPositionPermissionHandler(service *service.PositionPermissionService) *PositionPermissionHandler {
	return &PositionPermissionHandler{service: service}
}


// POST /position-permissions
func (h *PositionPermissionHandler) CreatePositionPermission(c *gin.Context) {
	var req dto.CreatePositionPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPerm, err := h.service.CreatePositionPermission(req)
	if err != nil {
		if err.Error() == "invalid employee_position_id or access_permission_id" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "this permission is already assigned to the position" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create position permission"})
		return
	}
	c.JSON(http.StatusCreated, newPerm)
}

// GET /position-permissions
func (h *PositionPermissionHandler) GetAllPositionPermissions(c *gin.Context) {
	permissions, err := h.service.GetAllPositionPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve position permissions"})
		return
	}
	if permissions == nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}
	c.JSON(http.StatusOK, permissions)
}

// PATCH /position-permissions/positions/:posId/permissions/:permId/status
func (h *PositionPermissionHandler) UpdatePositionPermissionActiveStatus(c *gin.Context) {
	posID, _ := strconv.Atoi(c.Param("posId"))
	permID, _ := strconv.Atoi(c.Param("permId"))

	var req dto.UpdatePositionPermissionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdatePositionPermissionActiveStatus(posID, permID, req); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Position permission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Position permission status updated successfully"})
}

// DELETE /api/v1/position-permissions/positions/:posId/permissions/:permId
func (h *PositionPermissionHandler) DeletePositionPermission(c *gin.Context) {
	posID, _ := strconv.Atoi(c.Param("posId"))
	permID, _ := strconv.Atoi(c.Param("permId"))

	if err := h.service.DeletePositionPermission(posID, permID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Position permission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete position permission"})
		return
	}
	c.Status(http.StatusNoContent)
}