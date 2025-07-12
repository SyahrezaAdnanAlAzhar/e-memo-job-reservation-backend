package main

import (
	"log"
	"net/http"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/handler"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.Connect()
	defer db.Close()

	// REPOSITORY
	// employeeRepo := repository.NewEmployeeRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)


	// SERVICE
	departmentService := service.NewDepartmentService(departmentRepo)


	// HANDLER
	// employeeHandler := handler.NewEmployeeHandler(employeeRepo)
	departmentHandler := handler.NewDepartmentHandler(departmentService)


	router := gin.Default()
	apiV1 := router.Group("/api/e-memo-job-reservation")
	{
		apiV1.GET("/departments", departmentHandler.GetAllDepartments)
	}

	
	log.Println("Starting server on :8080...")
	router.Run(":8080")
}
