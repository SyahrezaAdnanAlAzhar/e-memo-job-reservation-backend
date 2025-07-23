package router

import (
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/auth"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/handler"

	"github.com/gin-gonic/gin"
)

type AllHandlers struct {
	AuthHandler                *handler.AuthHandler
	DepartmentHandler          *handler.DepartmentHandler
	AreaHandler                *handler.AreaHandler
	PhysicalLocationHandler    *handler.PhysicalLocationHandler
	AccessPermissionHandler    *handler.AccessPermissionHandler
	SectionStatusTicketHandler *handler.SectionStatusTicketHandler
	StatusTicketHandler        *handler.StatusTicketHandler
	TicketHandler              *handler.TicketHandler
	PositionPermissionHandler  *handler.PositionPermissionHandler
	EmployeePositionHandler    *handler.EmployeePositionHandler
	WorkflowHandler            *handler.WorkflowHandler
	SpecifiedLocationHandler   *handler.SpecifiedLocationHandler
}

func SetupRouter(h *AllHandlers, authMiddleware *auth.AuthMiddleware) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/e-memo-job-reservation")

	public := api.Group("")
	{
		public.POST("/login", h.AuthHandler.Login)
		public.POST("/refresh", h.AuthHandler.RefreshToken)
	}

	private := api.Group("")
	private.Use(authMiddleware.JWTMiddleware())
	{
		private.POST("/logout", h.AuthHandler.Logout)

		setupMasterDataRoutes(private, h)

		setupMainDataRoutes(private, h)
	}

	return router
}

func setupMasterDataRoutes(group *gin.RouterGroup, h *AllHandlers) {
	deptRoutes := group.Group("/department")
	{
		deptRoutes.POST("", h.DepartmentHandler.CreateDepartment)
		deptRoutes.GET("", h.DepartmentHandler.GetAllDepartments)
		deptRoutes.GET("/:id", h.DepartmentHandler.GetDepartmentByID)
		deptRoutes.DELETE("/:id", h.DepartmentHandler.DeleteDepartment)
		deptRoutes.PUT("/:id", h.DepartmentHandler.UpdateDepartment)
		deptRoutes.PATCH("/:id/status", h.DepartmentHandler.UpdateDepartmentActiveStatus)
	}

	physicalLocationRoutes := group.Group("/physical-location")
	{
		physicalLocationRoutes.POST("", h.PhysicalLocationHandler.CreatePhysicalLocation)
		physicalLocationRoutes.GET("", h.PhysicalLocationHandler.GetAllPhysicalLocations)
		physicalLocationRoutes.GET("/:id", h.PhysicalLocationHandler.GetPhysicalLocationByID)
		physicalLocationRoutes.PUT("/:id", h.PhysicalLocationHandler.UpdatePhysicalLocation)
		physicalLocationRoutes.DELETE("/:id", h.PhysicalLocationHandler.DeletePhysicalLocation)
		physicalLocationRoutes.PATCH("/:id/status", h.PhysicalLocationHandler.UpdatePhysicalLocationActiveStatus)
	}

	accessPermissionRoutes := group.Group("/access-permission")
	{
		accessPermissionRoutes.POST("", h.AccessPermissionHandler.CreateAccessPermission)
		accessPermissionRoutes.GET("", h.AccessPermissionHandler.GetAllAccessPermissions)
		accessPermissionRoutes.GET("/:id", h.AccessPermissionHandler.GetAccessPermissionByID)
		accessPermissionRoutes.PUT("/:id", h.AccessPermissionHandler.UpdateAccessPermission)
		accessPermissionRoutes.DELETE("/:id", h.AccessPermissionHandler.DeleteAccessPermission)
		accessPermissionRoutes.PATCH("/:id/status", h.AccessPermissionHandler.UpdateAccessPermissionActiveStatus)
	}

	areaRoutes := group.Group("/area")
	{
		areaRoutes.POST("", h.AreaHandler.CreateArea)
		areaRoutes.GET("", h.AreaHandler.GetAllAreas)
		areaRoutes.GET("/:id", h.AreaHandler.GetAreaByID)
		areaRoutes.PUT("/:id", h.AreaHandler.UpdateArea)
		areaRoutes.PATCH("/:id/status", h.AreaHandler.UpdateAreaActiveStatus)
	}

	statusTicketRoutes := group.Group("/status-ticket")
	{
		statusTicketRoutes.POST("", h.StatusTicketHandler.CreateStatusTicket)
		statusTicketRoutes.GET("", h.StatusTicketHandler.GetAllStatusTickets)
		statusTicketRoutes.GET("/:id", h.StatusTicketHandler.GetStatusTicketByID)
		statusTicketRoutes.DELETE("/:id", h.StatusTicketHandler.DeleteStatusTicket)
		statusTicketRoutes.PATCH("/:id/status", h.StatusTicketHandler.UpdateStatusTicketActiveStatus)
		statusTicketRoutes.PUT("/reorder", h.StatusTicketHandler.ReorderStatusTickets)
	}

	sectionRoutes := group.Group("/section-status-ticket")
	{
		sectionRoutes.POST("", h.SectionStatusTicketHandler.CreateSectionStatusTicket)
	}

	posPermRoutes := group.Group("/position-permissions")
	{
		posPermRoutes.POST("", h.PositionPermissionHandler.CreatePositionPermission)
		posPermRoutes.GET("", h.PositionPermissionHandler.GetAllPositionPermissions)
		posPermRoutes.PATCH("/positions/:posId/permissions/:permId/status", h.PositionPermissionHandler.UpdatePositionPermissionActiveStatus)
		posPermRoutes.DELETE("/positions/:posId/permissions/:permId", h.PositionPermissionHandler.DeletePositionPermission)
	}
	posRoutes := group.Group("/employee-position")
	{
		posRoutes.POST("", h.EmployeePositionHandler.CreateEmployeePosition)
		posRoutes.GET("", h.EmployeePositionHandler.GetAllEmployeePositions)
		posRoutes.GET("/:id", h.EmployeePositionHandler.GetEmployeePositionByID)
		posRoutes.PUT("/:id", h.EmployeePositionHandler.UpdateEmployeePosition)
		posRoutes.DELETE("/:id", h.EmployeePositionHandler.DeleteEmployeePosition)
		posRoutes.PATCH("/:id/status", h.EmployeePositionHandler.UpdateEmployeePositionActiveStatus)
	}
	workflowRoutes := group.Group("/workflow")
	{
		workflowRoutes.POST("", h.WorkflowHandler.CreateWorkflow)
		workflowRoutes.GET("", h.WorkflowHandler.GetAllWorkflows)
		workflowRoutes.GET("/:id", h.WorkflowHandler.GetWorkflowByID)
		workflowRoutes.PUT("/:id", h.WorkflowHandler.UpdateWorkflow)
		workflowRoutes.DELETE("/:id", h.WorkflowHandler.DeleteWorkflow)
		workflowRoutes.PATCH("/:id/status", h.WorkflowHandler.UpdateWorkflowActiveStatus)

		stepRoutes := group.Group("/workflow-step")
		{
			stepRoutes.POST("", h.WorkflowHandler.AddWorkflowStep)
			stepRoutes.GET("", h.WorkflowHandler.GetAllWorkflowSteps)
			stepRoutes.GET("/:id", h.WorkflowHandler.GetWorkflowStepByID)
			stepRoutes.DELETE("/:id", h.WorkflowHandler.DeleteWorkflowStep)
			stepRoutes.PATCH("/:id/status", h.WorkflowHandler.UpdateWorkflowStepActiveStatus)
		}
	}
	specLocRoutes := group.Group("/specified-location")
	{
		specLocRoutes.POST("", h.SpecifiedLocationHandler.CreateSpecifiedLocation)
		specLocRoutes.GET("", h.SpecifiedLocationHandler.GetAllSpecifiedLocations)
		specLocRoutes.GET("/:id", h.SpecifiedLocationHandler.GetSpecifiedLocationByID)
		specLocRoutes.PUT("/:id", h.SpecifiedLocationHandler.UpdateSpecifiedLocation)
		specLocRoutes.DELETE("/:id", h.SpecifiedLocationHandler.DeleteSpecifiedLocation)
		specLocRoutes.PATCH("/:id/status", h.SpecifiedLocationHandler.UpdateSpecifiedLocationActiveStatus)
	}
}

// MAIN TICKET
func setupMainDataRoutes(group *gin.RouterGroup, h *AllHandlers) {
	ticketRoutes := group.Group("/ticket")
	{
		ticketRoutes.POST("", h.TicketHandler.CreateTicket)
		ticketRoutes.GET("", h.TicketHandler.GetAllTickets)
		ticketRoutes.GET("/:id", h.TicketHandler.GetTicketByID)
		ticketRoutes.PUT("/:id", h.TicketHandler.UpdateTicket)
		ticketRoutes.PUT("/reorder", h.TicketHandler.ReorderTickets)
		ticketRoutes.POST("/:id/progress", h.TicketHandler.ProgressTicketStatus)
		ticketRoutes.PUT("/:id/change-status", h.TicketHandler.ChangeTicketStatus)
	}
}
