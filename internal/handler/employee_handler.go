package handler

import (
	"encoding/json"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository" 
	"log"
	"net/http"
)

type EmployeeHandler struct {
	Repo *repository.EmployeeRepository
}

func NewEmployeeHandler(repo *repository.EmployeeRepository) *EmployeeHandler {
	return &EmployeeHandler{Repo: repo}
}

func (h *EmployeeHandler) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.Repo.GetAllEmployees()
	if err != nil {
		log.Printf("Error getting employees: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}