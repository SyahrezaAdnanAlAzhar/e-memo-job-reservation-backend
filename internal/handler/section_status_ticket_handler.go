package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"

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

// GET /section-status-ticket
func (h *SectionStatusTicketHandler) GetAllSectionStatusTickets(c *gin.Context) {
	sections, err := h.service.GetAllSectionStatusTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sections"})
		return
	}
	c.JSON(http.StatusOK, sections)
}

// GET /section-status-ticket/:id
func (h *SectionStatusTicketHandler) GetSectionStatusTicketByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	section, err := h.service.GetSectionStatusTicketByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve section"})
		return
	}
	c.JSON(http.StatusOK, section)
}

// PATCH /section-status-ticket/:id/status
func (h *SectionStatusTicketHandler) UpdateSectionStatusTicketActiveStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateSectionStatusTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateSectionStatusTicketActiveStatus(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "cannot deactivate, must have at least two active sections" || err.Error() == "cannot deactivate the first section" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update section status", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Section status and related data updated successfully"})
}