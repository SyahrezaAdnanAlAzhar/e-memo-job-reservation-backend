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
	// MASTER DATA INDEPENDENT
	physicalLocationRepo := repository.NewPhysicalLocationRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)

	// MASTER DATA DEPENDENT
	areaRepo := repository.NewAreaRepository(db)
	statusTicketRepo := repository.NewStatusTicketRepository(db)

	// MAIN DATA
	ticketRepo := repository.NewTicketRepository(db)
	jobRepo := repository.NewJobRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)


	
	// SERVICE
	// MASTER DATA INDEPENDENT
	departmentService := service.NewDepartmentService(departmentRepo)
	physicalLocationService := service.NewPhysicalLocationService(physicalLocationRepo)

	// MASTER DATA DEPENDENT
	areaService := service.NewAreaService(areaRepo)
	statusTicketService := service.NewStatusTicketService(statusTicketRepo) 

	// MAIN DATA
	ticketService := service.NewTicketService(&service.TicketServiceConfig{
		TicketRepo:   ticketRepo,
		JobRepo:      jobRepo,
		WorkflowRepo: workflowRepo,
		DB:           db, 
	})



	// HANDLER
	// employeeHandler := handler.NewEmployeeHandler(employeeRepo)
	// MASTER DATA INDEPENDENT
	departmentHandler := handler.NewDepartmentHandler(departmentService)
	physicalLocationHandler := handler.NewPhysicalLocationHandler(physicalLocationService)

	// MASTER DATA DEPENDENT
	areaHandler := handler.NewAreaHandler(areaService)
	statusTicketHandler := handler.NewStatusTicketHandler(statusTicketService)

	// MAIN DATA
	ticketHandler := handler.NewTicketHandler(ticketService)

	// SETUP
	router := gin.Default()
	api := router.Group("/api/e-memo-job-reservation")
	{
		// MASTER DATA INDEPENDENT
		deptRoutes := api.Group("/department")
		{
			deptRoutes.POST("", departmentHandler.CreateDepartment)
			deptRoutes.GET("", departmentHandler.GetAllDepartments)
			deptRoutes.GET("/:id", departmentHandler.GetDepartmentByID)
			deptRoutes.DELETE("/:id", departmentHandler.DeleteDepartment)
			deptRoutes.PUT("/:id", departmentHandler.UpdateDepartment)
			deptRoutes.PATCH("/:id/status", departmentHandler.UpdateDepartmentActiveStatus)
		}
		physicalLocationRoutes := api.Group("/physical-location")
		{
			physicalLocationRoutes.POST("", physicalLocationHandler.CreatePhysicalLocation)
			physicalLocationRoutes.GET("", physicalLocationHandler.GetAllPhysicalLocations)
			physicalLocationRoutes.GET("/:id", physicalLocationHandler.GetPhysicalLocationByID)
			physicalLocationRoutes.PUT("/:id", physicalLocationHandler.UpdatePhysicalLocation)
			physicalLocationRoutes.DELETE("/:id", physicalLocationHandler.DeletePhysicalLocation)
			physicalLocationRoutes.PATCH("/:id/status", physicalLocationHandler.UpdatePhysicalLocationActiveStatus)
		}

		// MASTER DATA DEPENDENT
		areaRoutes := api.Group("/area")
		{
			areaRoutes.POST("", areaHandler.CreateArea)
			areaRoutes.GET("", areaHandler.GetAllAreas)
			areaRoutes.GET("/:id", areaHandler.GetAreaByID)
			areaRoutes.PUT("/:id", areaHandler.UpdateArea)
			areaRoutes.PATCH("/:id/status", areaHandler.UpdateAreaActiveStatus)
		}

		statusTicketRoutes := api.Group("/status-ticket")
		{
			statusTicketRoutes.POST("", statusTicketHandler.CreateStatusTicket)
			statusTicketRoutes.GET("", statusTicketHandler.GetAllStatusTickets)
			statusTicketRoutes.GET("/:id", statusTicketHandler.GetStatusTicketByID)
			statusTicketRoutes.DELETE("/:id", statusTicketHandler.DeleteStatusTicket)
			statusTicketRoutes.PATCH("/:id/status", statusTicketHandler.UpdateStatusTicketActiveStatus)
			statusTicketRoutes.PUT("/reorder", statusTicketHandler.ReorderStatusTickets)
		}

		// MAIN DATA
		ticketRoutes := api.Group("/ticket")
		{
			ticketRoutes.POST("", ticketHandler.CreateTicket)
		}
	}

	
	log.Println("Starting server on :8080...")
	router.Run(":8080")
}
