package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AreaHandler struct {
	service *service.AreaService
}

func NewAreaHandler(service *service.AreaService) *AreaHandler {
	return &AreaHandler{service: service}
}

// POST /area
func (h *AreaHandler) CreateArea(c *gin.Context) {
	var req repository.CreateAreaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newArea, err := h.service.CreateArea(req)
	if err != nil {
		if err.Error() == "area name already exists in this department" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create area"})
		return
	}

	c.JSON(http.StatusCreated, newArea)
}

// GET /area
func (h *AreaHandler) GetAllAreas(c *gin.Context) {
	filters := make(map[string]string)

	if isActive, exists := c.GetQuery("is_active"); exists {
		filters["is_active"] = isActive
	}

	if deptID, exists := c.GetQuery("department_id"); exists {
		filters["department_id"] = deptID
	}

	areas, err := h.service.GetAllAreas(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve areas"})
		return
	}

	if areas == nil {
		c.JSON(http.StatusOK, []repository.Area{})
		return
	}

	c.JSON(http.StatusOK, areas)
}

// GET /area/:id
func (h *AreaHandler) GetAreaByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area ID format"})
		return
	}

	area, err := h.service.GetAreaByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Area not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve area"})
		return
	}

	c.JSON(http.StatusOK, area)
}

// DELETE /area/:id
func (h *AreaHandler) DeleteArea(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area ID format"})
		return
	}

	err = h.service.DeleteArea(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Area not found or already deleted"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete area"})
		return
	}

	c.Status(http.StatusNoContent)
}

// PUT /area/:id
func (h *AreaHandler) UpdateArea(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area ID format"})
		return
	}

	var req repository.UpdateAreaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedArea, err := h.service.UpdateArea(id, req)
	if err != nil {
		if err.Error() == "area name already exists in this department" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Area not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update area"})
		return
	}

	c.JSON(http.StatusOK, updatedArea)
}
