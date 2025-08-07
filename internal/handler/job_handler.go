package handler

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/pkg/filehandler"
	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	commandService *service.JobService
	queryService   *service.JobQueryService
}

func NewJobHandler(commandService *service.JobService, queryService *service.JobQueryService) *JobHandler {
	return &JobHandler{commandService: commandService, queryService: queryService}
}

// GET /jobs
func (h *JobHandler) GetAllJobs(c *gin.Context) {
	var filters dto.JobFilter
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}

	jobs, err := h.queryService.GetAllJobs(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve jobs", "details": err.Error()})
		return
	}

	if jobs == nil {
		c.JSON(http.StatusOK, []dto.JobDetailResponse{})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// GET /jobs/:id
func (h *JobHandler) GetJobByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID format"})
		return
	}

	job, err := h.queryService.GetJobByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}

// PUT /jobs/:id/assign
func (h *JobHandler) AssignPIC(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID format"})
		return
	}
	userNPK := c.GetString("user_npk")

	var req dto.AssignPICRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.commandService.AssignPIC(c.Request.Context(), id, req, userNPK)
	if err != nil {
		switch err.Error() {
		case "job not found", "action performer not found", "new PIC employee data not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "user is not authorized to assign PIC for this job's department":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "new PIC must be from the same department as the job":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign PIC", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIC assigned successfully"})
}

// PUT /jobs/reorder
func (h *JobHandler) ReorderJobs(c *gin.Context) {
	userNPK := c.GetString("user_npk")

	var req dto.ReorderJobsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.commandService.ReorderJobs(c.Request.Context(), req, userNPK)
	if err != nil {
		switch err.Error() {
		case "user can only reorder jobs within their own department", "one or more job IDs do not belong to the specified department":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "action performer not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder jobs", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job priorities updated successfully"})
}

// GET /jobs/:id/available-actions
func (h *JobHandler) GetAvailableActions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID format"})
		return
	}
	userNPK := c.GetString("user_npk")

	actions, err := h.queryService.GetAvailableActions(c.Request.Context(), id, userNPK)
	if err != nil {
		if err.Error() == "job not found" || err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve available actions", "details": err.Error()})
		return
	}

	if actions == nil {
		c.JSON(http.StatusOK, []dto.AvailableActionResponse{})
		return
	}

	c.JSON(http.StatusOK, actions)
}

// POST /jobs/:id/report
func (h *JobHandler) UploadReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID format"})
		return
	}
	userNPK := c.GetString("user_npk")

	file, err := c.FormFile("report_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "report_file is required in the form data"})
		return
	}

	filePath, err := filehandler.SaveFile(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	err = h.commandService.UploadReport(c.Request.Context(), id, userNPK, filePath)
	if err != nil {
		os.Remove(filePath)

		if err.Error() == "user is not the assigned PIC for this job" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "job not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process report upload", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Report file uploaded successfully",
		"file_path": filePath,
	})
}
