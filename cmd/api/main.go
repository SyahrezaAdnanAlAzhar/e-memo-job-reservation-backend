package main

import (
	"log"

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


	// SETUP
	router := gin.Default()
	api := router.Group("/api/e-memo-job-reservation")
	{
		deptRoutes := api.Group("/departments")
		{
			deptRoutes.POST("", departmentHandler.CreateDepartment)
			deptRoutes.GET("", departmentHandler.GetAllDepartments)
			deptRoutes.GET("/:id", departmentHandler.GetDepartmentByID)
			deptRoutes.DELETE("/:id", departmentHandler.DeleteDepartment)
			deptRoutes.PUT("/:id", departmentHandler.UpdateDepartment)
		}
	}

	
	log.Println("Starting server on :8080...")
	router.Run(":8080")
}
