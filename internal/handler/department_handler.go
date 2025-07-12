package handler

import (
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"net/http"
	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	service *service.DepartmentService
}

func NewDepartmentHandler(service *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

// GET /departments
func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
	departments, err := h.service.GetAllDepartments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve departments"})
		return
	}
	c.JSON(http.StatusOK, departments)
}