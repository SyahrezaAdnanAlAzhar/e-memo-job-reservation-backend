package handler

import (
	"net/http"
	"strconv"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/util"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	service *service.FileService
}

func NewFileHandler(service *service.FileService) *FileHandler {
	return &FileHandler{service: service}
}

// GET /tickets/:id/files
func (h *FileHandler) GetAllFilesByTicketID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid ticket ID format", nil)
		return
	}

	files, err := h.service.GetAllFilesByTicketID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "ticket not found" {
			util.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve files", err.Error())
		return
	}

	util.SuccessResponse(c, http.StatusOK, files)
}
