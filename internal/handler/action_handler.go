package handler

import (
	"net/http"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/util"
	"github.com/gin-gonic/gin"
)

type ActionHandler struct {
	service *service.ActionService
}

func NewActionHandler(service *service.ActionService) *ActionHandler {
	return &ActionHandler{service: service}
}

// GET /actions
func (h *ActionHandler) GetAllActions(c *gin.Context) {
	actions, err := h.service.GetAllActions()
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve actions", err.Error())
		return
	}

	if actions == nil {
		util.SuccessResponse(c, http.StatusOK, []dto.ActionResponse{})
		return
	}

	util.SuccessResponse(c, http.StatusOK, actions)
}