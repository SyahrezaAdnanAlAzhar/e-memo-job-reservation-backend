package handler

import (
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SectionStatusTicketHandler struct {
	service *service.SectionStatusTicketService
}

func NewSectionStatusTicketHandler(service *service.SectionStatusTicketService) *SectionStatusTicketHandler {
	return &SectionStatusTicketHandler{service: service}
}

// POST /section-status-ticket
func (h *SectionStatusTicketHandler) CreateSectionStatusTicket(c *gin.Context) {
	var req dto.CreateSectionStatusTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newSection, err := h.service.CreateSectionStatusTicket(req)
	if err != nil {
		if err.Error() == "section name or sequence already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create section status ticket"})
		return
	}
	c.JSON(http.StatusCreated, newSection)
}