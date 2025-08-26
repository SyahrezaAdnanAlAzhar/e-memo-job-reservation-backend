package service

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type FileService struct {
	ticketRepo *repository.TicketRepository
	jobRepo    *repository.JobRepository
}

func NewFileService(ticketRepo *repository.TicketRepository, jobRepo *repository.JobRepository) *FileService {
	return &FileService{
		ticketRepo: ticketRepo,
		jobRepo:    jobRepo,
	}
}

func (s *FileService) GetAllFilesByTicketID(ctx context.Context, ticketID int) (*dto.AllFilesResponse, error) {
	supportFilePaths, supportFileTime, err := s.ticketRepo.GetSupportFilesByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err 
	}

	reportFilePaths, reportFileTime, err := s.jobRepo.GetReportFilesByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	response := &dto.AllFilesResponse{
		SupportFiles: formatFileResponses(supportFilePaths, supportFileTime),
		ReportFiles:  formatFileResponses(reportFilePaths, reportFileTime),
	}

	return response, nil
}

func formatFileResponses(paths []string, timestamp time.Time) []dto.FileResponse {
	if len(paths) == 0 || (len(paths) == 1 && paths[0] == "") {
		return []dto.FileResponse{}
	}

	responses := make([]dto.FileResponse, len(paths))
	for i, path := range paths {
		responses[i] = dto.FileResponse{
			FileName:   filepath.Base(path),
			FilePath:   path,
			FileType:   determineFileType(path),
			UploadedAt: timestamp,
		}
	}
	return responses
}

func determineFileType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp":
		return "image"
	case ".pdf":
		return "pdf"
	case ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt":
		return "document"
	case ".mp4", ".mov", ".avi", ".mkv":
		return "video"
	case ".zip", ".rar", ".7z", ".tar", ".gz":
		return "archive"
	default:
		return "unknown"
	}
}