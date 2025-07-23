package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/model"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type EmployeePositionService struct {
	positionRepo *repository.EmployeePositionRepository
	mappingRepo  *repository.PositionToWorkflowMappingRepository
	db           *sql.DB
}

func NewEmployeePositionService(positionRepo *repository.EmployeePositionRepository, mappingRepo *repository.PositionToWorkflowMappingRepository, db *sql.DB) *EmployeePositionService {
	return &EmployeePositionService{positionRepo: positionRepo, mappingRepo: mappingRepo, db: db}
}

// CREATE
func (s *EmployeePositionService) CreateEmployeePosition(ctx context.Context, req dto.CreateEmployeePositionRequest) (*model.EmployeePosition, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	newPos, err := s.positionRepo.Create(ctx, tx, req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return nil, errors.New("position name already exists")
		}
		return nil, err
	}

	// Langkah 2: Buat mapping di tabel position_to_workflow_mapping
	err = s.mappingRepo.CreateWorkflowMapping(ctx, tx, newPos.ID, req.WorkflowID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign_key_violation
			return nil, errors.New("invalid workflow_id")
		}
		return nil, err
	}

	// Jika kedua langkah berhasil, commit transaksi
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return newPos, nil
}

// GET ALL
func (s *EmployeePositionService) GetAllEmployeePositions() ([]model.EmployeePosition, error) {
	return s.positionRepo.FindAll()
}

// GET BY ID
func (s *EmployeePositionService) GetEmployeePositionByID(id int) (*model.EmployeePosition, error) {
	return s.positionRepo.FindByID(id)
}
