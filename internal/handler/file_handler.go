package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"e-memo-job-reservation-api/internal/service"
	"e-memo-job-reservation-api/internal/util"

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

func (h *FileHandler) DownloadFile(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		util.ErrorResponse(c, http.StatusBadRequest, "File path is required", nil)
		return
	}

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./uploads"
	}

	cleanedPath := filepath.Clean(filePath)

	if !strings.HasPrefix(cleanedPath, storagePath) {
		util.ErrorResponse(c, http.StatusForbidden, "Access to the requested file path is forbidden", nil)
		return
	}

	c.FileAttachment(cleanedPath, filepath.Base(cleanedPath))
}

func (h *FileHandler) ViewFile(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		util.ErrorResponse(c, http.StatusBadRequest, "File path is required", nil)
		return
	}

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./uploads"
	}

	cleanedPath := filepath.Clean(filePath)
	if !strings.HasPrefix(cleanedPath, storagePath) {
		util.ErrorResponse(c, http.StatusForbidden, "Access to the requested file path is forbidden", nil)
		return
	}

	c.File(cleanedPath)
}
