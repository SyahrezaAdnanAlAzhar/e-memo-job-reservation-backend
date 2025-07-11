package main

import (
	"log"
	"net/http"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/handler"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/pkg/database"
)

func main() {
	db := database.Connect()
	defer db.Close()
	employeeRepo := repository.NewEmployeeRepository(db)
	employeeHandler := handler.NewEmployeeHandler(employeeRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/employees", employeeHandler.GetAllEmployees)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
