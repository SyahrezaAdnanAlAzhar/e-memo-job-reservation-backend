package handler

import (
	"database/sql"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PhysicalLocationHandler struct {
	service *service.PhysicalLocationService
}

func NewPhysicalLocationHandler(service *service.PhysicalLocationService) *PhysicalLocationHandler {
	return &PhysicalLocationHandler{service: service}
}


// POST /physical-location
func (h *PhysicalLocationHandler) CreatePhysicalLocation(c *gin.Context) {
	var req repository.CreatePhysicalLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newLoc, err := h.service.CreatePhysicalLocation(req)
	if err != nil {
		if err.Error() == "physical location name already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create physical location"})
		return
	}
	c.JSON(http.StatusCreated, newLoc)
}



// GET /physical-location
func (h *PhysicalLocationHandler) GetAllPhysicalLocations(c *gin.Context) {
	filters := make(map[string]string)
	if isActive, exists := c.GetQuery("is_active"); exists {
		filters["is_active"] = isActive
	}

	locations, err := h.service.GetAllPhysicalLocations(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve physical locations"})
		return
	}

	if locations == nil {
		c.JSON(http.StatusOK, []repository.PhysicalLocation{})
		return
	}
	c.JSON(http.StatusOK, locations)
}


// GET /physical-location/:id
func (h *PhysicalLocationHandler) GetPhysicalLocationByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	location, err := h.service.GetPhysicalLocationByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Physical location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve physical location"})
		return
	}
	c.JSON(http.StatusOK, location)
}


// PUT /physical-location/:id
func (h *PhysicalLocationHandler) UpdatePhysicalLocation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req repository.UpdatePhysicalLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedLoc, err := h.service.UpdatePhysicalLocation(id, req)
	if err != nil {
		if err.Error() == "physical location name already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Physical location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update physical location"})
		return
	}
	c.JSON(http.StatusOK, updatedLoc)
}


// DELETE /physical-location/:id
func (h *PhysicalLocationHandler) DeletePhysicalLocation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeletePhysicalLocation(id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Physical location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete physical location"})
		return
	}
	c.Status(http.StatusNoContent)
}


// PATCH /physical-location/:id/status
func (h *PhysicalLocationHandler) UpdatePhysicalLocationActiveStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req repository.UpdatePhysicalLocationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdatePhysicalLocationActiveStatus(id, req); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Physical location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Physical location status updated successfully"})
}