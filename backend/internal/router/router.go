package router

import (
	"cowork/internal/config"
	"cowork/internal/db"
	"cowork/internal/handler"
	"cowork/internal/middleware"
	"cowork/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, database *gorm.DB, cfg *config.Config) {
	// Apply global middleware on the engine.
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// Initialize JWT config for middleware
	middleware.SetJWTConfig(&cfg.JWT)

	// Initialize auth service and handler
	authSvc := &service.AuthService{
		DB:    database,
		Redis: db.Redis,
		JWT:   &cfg.JWT,
	}
	authHandler := &handler.AuthHandler{Svc: authSvc}

	// Serve uploaded files
	r.Static("/uploads", "./uploads")

	api := r.Group("/api")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes
	auth := api.Group("/auth")
	{
		// Public routes
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.GET("/captcha", authHandler.GetCaptcha)

		// Protected routes
		auth.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
	}

	// Projects
	projectSvc := &service.ProjectService{DB: database}
	projectHandler := &handler.ProjectHandler{Svc: projectSvc}

	projects := api.Group("/projects")
	projects.Use(middleware.AuthMiddleware())
	{
		projects.POST("", middleware.RequireRole("ProjectManager"), projectHandler.Create)
		projects.GET("", projectHandler.ListMyProjects)
		projects.GET("/:id", projectHandler.GetDetail)
		projects.PUT("/:id", projectHandler.Update)
		projects.POST("/:id/submit", projectHandler.SubmitForApproval)
		projects.POST("/:id/approve", middleware.RequireRole("ProjectManager"), projectHandler.Approve)
		projects.POST("/:id/start", projectHandler.StartDevelopment)
		projects.POST("/:id/complete", projectHandler.Complete)
		projects.POST("/:id/archive", projectHandler.Archive)
		projects.POST("/:id/members", projectHandler.AddMember)
		projects.DELETE("/:id/members/:uid", projectHandler.RemoveMember)
		projects.GET("/:id/members", projectHandler.ListMembers)
	}

	// Kanban
	kanbanHandler := &handler.KanbanHandler{DB: database}
	kanban := api.Group("")
	kanban.Use(middleware.AuthMiddleware())
	{
		kanban.GET("/projects/:id/kanban", kanbanHandler.GetKanban)
	}

	// Comments
	commentSvc := &service.CommentService{DB: database}
	commentHandler := &handler.CommentHandler{Svc: commentSvc}

	comments := api.Group("")
	comments.Use(middleware.AuthMiddleware())
	{
		comments.GET("/projects/:id/comments", commentHandler.ListByProject)
		comments.POST("/projects/:id/comments", commentHandler.CreateByProject)
		comments.GET("/tasks/:id/comments", commentHandler.ListByTask)
		comments.POST("/tasks/:id/comments", commentHandler.CreateByTask)
		comments.POST("/comments/:id/reply", commentHandler.CreateReply)
		comments.DELETE("/comments/:id", commentHandler.Delete)
	}

	// Attachments
	attachmentSvc := &service.AttachmentService{DB: database}
	attachmentHandler := &handler.AttachmentHandler{Svc: attachmentSvc}

	attachments := api.Group("")
	attachments.Use(middleware.AuthMiddleware())
	{
		attachments.POST("/attachments/upload", attachmentHandler.Upload)
		attachments.DELETE("/attachments/:id", attachmentHandler.Delete)
		attachments.GET("/tasks/:id/attachments", attachmentHandler.ListByTask)
		attachments.GET("/projects/:id/attachments", attachmentHandler.ListByProject)
	}

	// Notifications
	notificationSvc := &service.NotificationService{DB: database}
	notificationHandler := &handler.NotificationHandler{Svc: notificationSvc}

	notifications := api.Group("")
	notifications.Use(middleware.AuthMiddleware())
	{
		notifications.GET("/notifications", notificationHandler.List)
		notifications.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
		notifications.PUT("/notifications/read-all", notificationHandler.MarkAllRead)
		notifications.PUT("/notifications/:id/read", notificationHandler.MarkRead)
	}

	// Operation logs (ProjectManager only)
	logSvc := &service.LogService{DB: database}
	logHandler := &handler.LogHandler{Svc: logSvc}

	logs := api.Group("")
	logs.Use(middleware.AuthMiddleware())
	{
		logs.GET("/logs", middleware.RequireRole("ProjectManager"), logHandler.List)
	}

	// Tasks
	taskSvc := &service.TaskService{DB: database}
	taskHandler := &handler.TaskHandler{Svc: taskSvc}

	tasks := api.Group("")
	tasks.Use(middleware.AuthMiddleware())
	{
		tasks.POST("/projects/:id/tasks", taskHandler.Create)
		tasks.GET("/projects/:id/tasks", taskHandler.ListByProject)
		tasks.GET("/tasks/:id", taskHandler.GetDetail)
		tasks.PUT("/tasks/:id", taskHandler.Update)
		tasks.DELETE("/tasks/:id", taskHandler.Delete)
		tasks.PUT("/tasks/:id/assign", taskHandler.Assign)
		tasks.PUT("/tasks/:id/status", taskHandler.TransferStatus)
	}
}
