package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/util"
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
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	jobs, err := h.queryService.GetAllJobs(filters)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve jobs", err.Error())
		return
	}

	if jobs == nil {
		util.SuccessResponse(c, http.StatusOK, []dto.JobDetailResponse{})
		return
	}

	util.SuccessResponse(c, http.StatusOK, jobs)
}

// GET /jobs/:id
func (h *JobHandler) GetJobByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid job ID format", nil)
		return
	}

	job, err := h.queryService.GetJobByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			util.ErrorResponse(c, http.StatusNotFound, "Job not found", nil)
			return
		}
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve job", err.Error())
		return
	}

	util.SuccessResponse(c, http.StatusOK, job)
}

// PUT /jobs/:id/assign
func (h *JobHandler) AssignPIC(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid job ID format", nil)
		return
	}
	userNPK := c.GetString("user_npk")

	var req dto.AssignPICRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = h.commandService.AssignPIC(c.Request.Context(), id, req, userNPK)
	if err != nil {
		switch err.Error() {
		case "job not found", "action performer not found", "new PIC employee data not found":
			util.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		case "user is not authorized to assign PIC for this job's department":
			util.ErrorResponse(c, http.StatusForbidden, err.Error(), nil)
		case "new PIC must be from the same department as the job":
			util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		default:
			util.ErrorResponse(c, http.StatusInternalServerError, "Failed to assign PIC", err.Error())
		}
		return
	}

	util.SuccessResponse(c, http.StatusOK, gin.H{"message": "PIC assigned successfully"})
}

// PUT /jobs/reorder
func (h *JobHandler) ReorderJobs(c *gin.Context) {
	userNPK := c.GetString("user_npk")

	var req dto.ReorderJobsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.commandService.ReorderJobs(c.Request.Context(), req, userNPK)
	if err != nil {
		switch err.Error() {
		case "user can only reorder jobs within their own department", "one or more job IDs do not belong to the specified department":
			util.ErrorResponse(c, http.StatusForbidden, err.Error(), nil)
		case "action performer not found":
			util.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			util.ErrorResponse(c, http.StatusInternalServerError, "Failed to reorder jobs", err.Error())
		}
		return
	}

	util.SuccessResponse(c, http.StatusOK, gin.H{"message": "Job priorities updated successfully"})
}

// GET /jobs/:id/available-actions
func (h *JobHandler) GetAvailableActions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid job ID format", nil)
		return
	}
	userNPK := c.GetString("user_npk")

	actions, err := h.queryService.GetAvailableActions(c.Request.Context(), id, userNPK)
	if err != nil {
		if err.Error() == "job not found" || err.Error() == "user not found" {
			util.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve available actions", err.Error())
		return
	}

	if actions == nil {
		util.SuccessResponse(c, http.StatusOK, []dto.AvailableActionResponse{})
		return
	}

	util.SuccessResponse(c, http.StatusOK, actions)
}

// POST /jobs/:id/files
func (h *JobHandler) AddReportFiles(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid job ID format", nil)
		return
	}
	userNPK := c.GetString("user_npk")

	form, err := c.MultipartForm()
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid form data", err.Error())
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		util.ErrorResponse(c, http.StatusBadRequest, "At least one file must be uploaded", nil)
		return
	}

	savedPaths, err := h.commandService.AddReportFiles(c.Request.Context(), c, id, userNPK, files)
	if err != nil {
		if err.Error() == "job not found" {
			util.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		} else if err.Error() == "user is not the assigned PIC for this job" {
			util.ErrorResponse(c, http.StatusForbidden, err.Error(), nil)
		} else {
			util.ErrorResponse(c, http.StatusInternalServerError, "Failed to add report files", err.Error())
		}
		return
	}
	util.SuccessResponse(c, http.StatusOK, gin.H{"message": "Report files uploaded successfully", "file_paths": savedPaths})
}

// DELETE /jobs/:id/files
func (h *JobHandler) RemoveReportFiles(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid job ID format", nil)
		return
	}
	userNPK := c.GetString("user_npk")

	var req dto.DeleteJobFilesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err = h.commandService.RemoveReportFiles(c.Request.Context(), id, userNPK, req)
	if err != nil {
		if err.Error() == "job not found" {
			util.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		} else if err.Error() == "user is not the assigned PIC for this job" {
			util.ErrorResponse(c, http.StatusForbidden, err.Error(), nil)
		} else {
			util.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove report files", err.Error())
		}
		return
	}
	util.SuccessResponse(c, http.StatusOK, gin.H{"message": "Selected report files removed successfully"})
}
