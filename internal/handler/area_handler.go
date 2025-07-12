package handler

import (
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"   
	"net/http"

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