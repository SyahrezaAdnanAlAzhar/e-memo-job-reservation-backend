package main

import (
	"log"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/auth"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/handler"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/router"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/websocket"
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

	hub := websocket.NewHub()
	go hub.Run()

	// DEPENDENCY INITIALIZATION (WIRING)
	// REPOSITORY
	actorRoleRepo := repository.NewActorRoleRepository(db)
	actorRoleMappingRepo := repository.NewActorRoleMappingRepository(db)
	appUserRepo := repository.NewAppUserRepository(db)
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
	jobQueryRepo := repository.NewJobQueryRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)
	trackStatusTicketRepo := repository.NewTrackStatusTicketRepository(db)
	positionPermissionRepo := repository.NewPositionPermissionRepository(db)
	employeePositionRepo := repository.NewEmployeePositionRepository(db)
	positionToWorkflowMappingRepo := repository.NewPositionToWorkflowMappingRepository(db)
	workflowStepRepo := repository.NewWorkflowStepRepository(db)
	specifiedLocationRepo := repository.NewSpecifiedLocationRepository(db)
	statusTransitionRepo := repository.NewStatusTransitionRepository(db)
	rejectedTicketRepo := repository.NewRejectedTicketRepository(db)
	ticketActionLogRepo := repository.NewTicketActionLogRepository(db)
	actionRepo := repository.NewActionRepository(db)

	// SERVICE
	authService := service.NewAuthService(authRepo, appUserRepo, positionPermissionRepo, employeeRepo)
	departmentService := service.NewDepartmentService(departmentRepo)
	employeeService := service.NewEmployeeService(employeeRepo)
	areaService := service.NewAreaService(areaRepo)
	physicalLocationService := service.NewPhysicalLocationService(physicalLocationRepo)
	accessPermissionService := service.NewAccessPermissionService(accessPermissionRepo)
	sectionStatusTicketService := service.NewSectionStatusTicketService(sectionStatusTicketRepo, statusTicketRepo, ticketRepo, db)
	statusTicketService := service.NewStatusTicketService(statusTicketRepo)
	jobService := service.NewJobService(jobRepo, jobQueryRepo, employeeRepo, positionPermissionRepo, db, hub)
	jobQueryService := service.NewJobQueryService(jobQueryRepo, employeeRepo, positionPermissionRepo)
	positionPermissionService := service.NewPositionPermissionService(positionPermissionRepo)
	workflowService := service.NewWorkflowService(workflowRepo, workflowStepRepo, db)
	specifiedLocationService := service.NewSpecifiedLocationService(specifiedLocationRepo)
	ticketActionService := service.NewTicketActionService(&service.TicketActionServiceConfig{
		TicketRepo:            ticketRepo,
		JobRepo:               jobRepo,
		EmployeeRepo:          employeeRepo,
		TrackStatusTicketRepo: trackStatusTicketRepo,
		StatusTransitionRepo:  statusTransitionRepo,
		ActorRoleMappingRepo:  actorRoleMappingRepo,
	})
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

	ticketQueryService := service.NewTicketQueryService(ticketRepo)

	ticketCommandService := service.NewTicketCommandService(&service.TicketCommandServiceConfig{
		DB:                    db,
		TicketRepo:            ticketRepo,
		JobRepo:               jobRepo,
		WorkflowRepo:          workflowRepo,
		TrackStatusTicketRepo: trackStatusTicketRepo,
		EmployeeRepo:          employeeRepo,
		DepartmentRepo:        departmentRepo,
		ActorRoleMappingRepo:  actorRoleMappingRepo,
		ActorRoleRepo:         actorRoleRepo,
		StatusTransitionRepo:  statusTransitionRepo,
		Hub:                   hub,
		QueryService:          ticketQueryService,
	})

	ticketWorkflowService := service.NewTicketWorkflowService(&service.TicketWorkflowServiceConfig{
		DB:                    db,
		TicketRepo:            ticketRepo,
		JobRepo:               jobRepo,
		EmployeeRepo:          employeeRepo,
		TrackStatusTicketRepo: trackStatusTicketRepo,
		StatusTicketRepo:      statusTicketRepo,
		StatusTransitionRepo:  statusTransitionRepo,
		ActorRoleRepo:         actorRoleRepo,
		ActorRoleMappingRepo:  actorRoleMappingRepo,
		TicketActionLogRepo:   ticketActionLogRepo,
	})

	ticketPriorityService := service.NewTicketPriorityService(db, hub, ticketRepo, employeeRepo)

	ticketHandler := handler.NewTicketHandler(&handler.TicketHandlerConfig{
		QueryService:    ticketQueryService,
		CommandService:  ticketCommandService,
		WorkflowService: ticketWorkflowService,
		PriorityService: ticketPriorityService,
		ActionService:   ticketActionService,
	})

	actionService := service.NewActionService(actionRepo)

	// HANDLER
	wsHandler := handler.NewWebSocketHandler(hub, authRepo)

	allHandlers := &router.AllHandlers{
		AuthHandler:                handler.NewAuthHandler(authService),
		DepartmentHandler:          handler.NewDepartmentHandler(departmentService),
		EmployeeHandler:            handler.NewEmployeeHandler(employeeService),
		AreaHandler:                handler.NewAreaHandler(areaService),
		PhysicalLocationHandler:    handler.NewPhysicalLocationHandler(physicalLocationService),
		AccessPermissionHandler:    handler.NewAccessPermissionHandler(accessPermissionService),
		SectionStatusTicketHandler: handler.NewSectionStatusTicketHandler(sectionStatusTicketService),
		StatusTicketHandler:        handler.NewStatusTicketHandler(statusTicketService),
		PositionPermissionHandler:  handler.NewPositionPermissionHandler(positionPermissionService),
		EmployeePositionHandler:    handler.NewEmployeePositionHandler(employeePositionService),
		WorkflowHandler:            handler.NewWorkflowHandler(workflowService),
		SpecifiedLocationHandler:   handler.NewSpecifiedLocationHandler(specifiedLocationService),
		RejectedTicketHandler:      handler.NewRejectedTicketHandler(rejectedTicketService),
		JobHandler:                 handler.NewJobHandler(jobService, jobQueryService),
		ActionHandler:              handler.NewActionHandler(actionService),
		TicketHandler:              ticketHandler,
	}

	allRepositories := &router.AllRepositories{
		PositionPermissionRepo: positionPermissionRepo,
	}

	// MIDDLEWARE
	authMiddleware := auth.NewAuthMiddleware(authRepo)

	// SET UP AND RUN SERVER
	appRouter := router.SetupRouter(allHandlers, allRepositories, authMiddleware, wsHandler)

	log.Println("Starting server on :8080...")
	appRouter.Run(":8080")
}
