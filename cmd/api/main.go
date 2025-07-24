package main

import (
	"log"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/auth"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/handler"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/router"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/pkg/database"
	redisClient "github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/pkg/redis"

	"github.com/joho/godotenv"
)

func main() {
	// INITIAL SET UP
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables from OS")
	}

	db := database.Connect()
	defer db.Close()

	rdb := redisClient.Connect()
	defer rdb.Close()

	// DEPENDENCY INITIALIZATION (WIRING)
	// REPOSITORY
	authRepo := repository.NewAuthRepository(rdb)
	employeeRepo := repository.NewEmployeeRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)
	areaRepo := repository.NewAreaRepository(db)
	physicalLocationRepo := repository.NewPhysicalLocationRepository(db)
	accessPermissionRepo := repository.NewAccessPermissionRepository(db)
	sectionStatusTicketRepo := repository.NewSectionStatusTicketRepository(db)
	statusTicketRepo := repository.NewStatusTicketRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	jobRepo := repository.NewJobRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)
	trackStatusTicketRepo := repository.NewTrackStatusTicketRepository(db)
	positionPermissionRepo := repository.NewPositionPermissionRepository(db)
	employeePositionRepo := repository.NewEmployeePositionRepository(db)
	positionToWorkflowMappingRepo := repository.NewPositionToWorkflowMappingRepository(db)
	workflowStepRepo := repository.NewWorkflowStepRepository(db)
	specifiedLocationRepo := repository.NewSpecifiedLocationRepository(db)
	rejectedTicketRepo := repository.NewRejectedTicketRepository(db)

	// SERVICE
	authService := service.NewAuthService(authRepo, employeeRepo)
	departmentService := service.NewDepartmentService(departmentRepo)
	areaService := service.NewAreaService(areaRepo)
	physicalLocationService := service.NewPhysicalLocationService(physicalLocationRepo)
	accessPermissionService := service.NewAccessPermissionService(accessPermissionRepo)
	sectionStatusTicketService := service.NewSectionStatusTicketService(sectionStatusTicketRepo, statusTicketRepo, ticketRepo, db)
	statusTicketService := service.NewStatusTicketService(statusTicketRepo)
	positionPermissionService := service.NewPositionPermissionService(positionPermissionRepo)
	workflowService := service.NewWorkflowService(workflowRepo, workflowStepRepo, db)
	specifiedLocationService := service.NewSpecifiedLocationService(specifiedLocationRepo)
	employeePositionService := service.NewEmployeePositionService(
		employeePositionRepo,
		positionToWorkflowMappingRepo,
		ticketRepo,
		statusTicketRepo,
		db)
	rejectedTicketService := service.NewRejectedTicketService(
		rejectedTicketRepo,
		ticketRepo,
		trackStatusTicketRepo,
		statusTicketRepo,
		employeeRepo,
		db,
	)
	ticketService := service.NewTicketService(&service.TicketServiceConfig{
		TicketRepo:            ticketRepo,
		JobRepo:               jobRepo,
		WorkflowRepo:          workflowRepo,
		TrackStatusTicketRepo: trackStatusTicketRepo,
		EmployeeRepo:          employeeRepo,
		StatusTicketRepo:      statusTicketRepo,
		DB:                    db,
	})

	// HANDLER
	allHandlers := &router.AllHandlers{
		AuthHandler:                handler.NewAuthHandler(authService),
		DepartmentHandler:          handler.NewDepartmentHandler(departmentService),
		AreaHandler:                handler.NewAreaHandler(areaService),
		PhysicalLocationHandler:    handler.NewPhysicalLocationHandler(physicalLocationService),
		AccessPermissionHandler:    handler.NewAccessPermissionHandler(accessPermissionService),
		SectionStatusTicketHandler: handler.NewSectionStatusTicketHandler(sectionStatusTicketService),
		StatusTicketHandler:        handler.NewStatusTicketHandler(statusTicketService),
		TicketHandler:              handler.NewTicketHandler(ticketService),
		PositionPermissionHandler:  handler.NewPositionPermissionHandler(positionPermissionService),
		EmployeePositionHandler:    handler.NewEmployeePositionHandler(employeePositionService),
		WorkflowHandler:            handler.NewWorkflowHandler(workflowService),
		SpecifiedLocationHandler:   handler.NewSpecifiedLocationHandler(specifiedLocationService),
		RejectedTicketHandler:      handler.NewRejectedTicketHandler(rejectedTicketService),
	}

	// MIDDLEWARE
	authMiddleware := auth.NewAuthMiddleware(authRepo)

	// SET UP AND RUN SERVER
	appRouter := router.SetupRouter(allHandlers, authMiddleware)

	log.Println("Starting server on :8080...")
	appRouter.Run(":8080")
}
