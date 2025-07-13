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