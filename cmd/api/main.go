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
	areaRepo := repository.NewAreaRepository(db)


	// SERVICE
	departmentService := service.NewDepartmentService(departmentRepo)
	areaService := service.NewAreaService(areaRepo)


	// HANDLER
	// employeeHandler := handler.NewEmployeeHandler(employeeRepo)
	departmentHandler := handler.NewDepartmentHandler(departmentService)
	areaHandler := handler.NewAreaHandler(areaService)


	// SETUP
	router := gin.Default()
	api := router.Group("/api/e-memo-job-reservation")
	{
		deptRoutes := api.Group("/department")
		{
			deptRoutes.POST("", departmentHandler.CreateDepartment)
			deptRoutes.GET("", departmentHandler.GetAllDepartments)
			deptRoutes.GET("/:id", departmentHandler.GetDepartmentByID)
			deptRoutes.DELETE("/:id", departmentHandler.DeleteDepartment)
			deptRoutes.PUT("/:id", departmentHandler.UpdateDepartment)
			deptRoutes.PATCH("/:id/status", departmentHandler.UpdateDepartmentActiveStatus)
		}

		areaRoutes := api.Group("/area")
		{
			areaRoutes.POST("", areaHandler.CreateArea)
			areaRoutes.GET("", areaHandler.GetAllAreas)
			areaRoutes.GET("/:id", areaHandler.GetAreaByID)
			areaRoutes.PUT("/:id", areaHandler.UpdateArea)
			areaRoutes.PATCH("/:id/status", areaHandler.UpdateAreaActiveStatus)
		}
	}

	
	log.Println("Starting server on :8080...")
	router.Run(":8080")
}
