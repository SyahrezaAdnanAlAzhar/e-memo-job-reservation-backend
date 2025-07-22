package handler

import (
	"database/sql"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccessPermissionHandler struct {
	service *service.AccessPermissionService
}

func NewAccessPermissionHandler(service *service.AccessPermissionService) *AccessPermissionHandler {
	return &AccessPermissionHandler{service: service}
}

// POST /access-permissions
func (h *AccessPermissionHandler) CreateAccessPermission(c *gin.Context) {
	var req repository.CreateAccessPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPermission, err := h.service.CreateAccessPermission(req)
	if err != nil {
		if err.Error() == "permission name already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create permission"})
		return
	}
	c.JSON(http.StatusCreated, newPermission)
}

// GET /access-permissions
func (h *AccessPermissionHandler) GetAllAccessPermissions(c *gin.Context) {
	permissions, err := h.service.GetAllAccessPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve permissions"})
		return
	}
	if permissions == nil {
		c.JSON(http.StatusOK, []repository.AccessPermission{})
		return
	}
	c.JSON(http.StatusOK, permissions)
}

// GET /access-permissions/:id
func (h *AccessPermissionHandler) GetAccessPermissionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	permission, err := h.service.GetAccessPermissionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve permission"})
		return
	}
	c.JSON(http.StatusOK, permission)
}

// PUT /access-permissions/:id
func (h *AccessPermissionHandler) UpdateAccessPermission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req repository.UpdateAccessPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPermission, err := h.service.UpdateAccessPermission(id, req)
	if err != nil {
		if err.Error() == "permission name already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update permission"})
		return
	}
	c.JSON(http.StatusOK, updatedPermission)
}

// DELETE /access-permissions/:id
func (h *AccessPermissionHandler) DeleteAccessPermission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteAccessPermission(id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete permission"})
		return
	}
	c.Status(http.StatusNoContent)
}

// PATCH /access-permissions/:id/status
func (h *AccessPermissionHandler) UpdateAccessPermissionActiveStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req repository.UpdateAccessPermissionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateAccessPermissionActiveStatus(id, req); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Permission status updated successfully"})
}