package handler

import (
	"net/http"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/util"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	service *service.EmployeeService
}

func NewEmployeeHandler(service *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) GetAllEmployees(c *gin.Context) {
	var filters dto.EmployeeFilter
	if err := c.ShouldBindQuery(&filters); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	employees, err := h.service.GetAllEmployees(filters)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve employees", err.Error())
		return
	}

	if employees == nil {
		util.SuccessResponse(c, http.StatusOK, []gin.H{})
		return
	}
	util.SuccessResponse(c, http.StatusOK, employees)
}
