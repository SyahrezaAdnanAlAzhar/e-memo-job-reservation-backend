package handler

import (
	"database/sql"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SpecifiedLocationHandler struct {
	service *service.SpecifiedLocationService
}

func NewSpecifiedLocationHandler(service *service.SpecifiedLocationService) *SpecifiedLocationHandler {
	return &SpecifiedLocationHandler{service: service}
}

// POST /specified-location
func (h *SpecifiedLocationHandler) CreateSpecifiedLocation(c *gin.Context) {
	var req dto.CreateSpecifiedLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newLoc, err := h.service.CreateSpecifiedLocation(req)
	if err != nil {
		if err.Error() == "invalid physical_location_id" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "location name already exists in this physical location" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create specified location"})
		return
	}
	c.JSON(http.StatusCreated, newLoc)
}

// GET /specified-location
func (h *SpecifiedLocationHandler) GetAllSpecifiedLocations(c *gin.Context) {
	physicalLocationIDStr := c.Query("physical_location_id")
	if physicalLocationIDStr != "" {
		physicalLocationID, err := strconv.Atoi(physicalLocationIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid physical_location_id format"})
			return
		}
		locations, err := h.service.GetSpecifiedLocationsByPhysicalLocationID(physicalLocationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve specified locations"})
			return
		}
		c.JSON(http.StatusOK, locations)
		return
	}

	locations, err := h.service.GetAllSpecifiedLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve all specified locations"})
		return
	}
	c.JSON(http.StatusOK, locations)
}

// GET /specified-location/:id
func (h *SpecifiedLocationHandler) GetSpecifiedLocationByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	location, err := h.service.GetSpecifiedLocationByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Specified location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve specified location"})
		return
	}
	c.JSON(http.StatusOK, location)
}

// PUT /specified-location/:id
func (h *SpecifiedLocationHandler) UpdateSpecifiedLocation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req dto.UpdateSpecifiedLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedLoc, err := h.service.UpdateSpecifiedLocation(id, req)
	if err != nil {
		if err.Error() == "location name already exists in this physical location" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Specified location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update specified location"})
		return
	}
	c.JSON(http.StatusOK, updatedLoc)
}

// DELETE /specified-location/:id
func (h *SpecifiedLocationHandler) DeleteSpecifiedLocation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteSpecifiedLocation(id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Specified location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete specified location"})
		return
	}
	c.Status(http.StatusNoContent)
}

// PATCH /specified-location/:id/status
func (h *SpecifiedLocationHandler) UpdateSpecifiedLocationActiveStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req dto.UpdateSpecifiedLocationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateSpecifiedLocationActiveStatus(id, req); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Specified location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Specified location status updated successfully"})
}