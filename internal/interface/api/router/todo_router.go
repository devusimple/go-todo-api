package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"todo-api/internal/domain/repository"
	"todo-api/internal/domain/usecase"
	"todo-api/internal/interface/api/handler"
	"todo-api/internal/interface/api/middleware"
)

// SetupTodoRoutes sets up routes related to todo operations
func SetupTodoRoutes(
	e *echo.Echo,
	todoRepo repository.TodoRepository,
	authMiddleware *middleware.AuthMiddleware,
	logger *logrus.Logger,
) {
	// Initialize todo use case
	todoUseCase := usecase.NewTodoUseCase(todoRepo)

	// Initialize todo handler
	todoHandler := handler.NewTodoHandler(todoUseCase, logger)

	// Define todo routes
	todoGroup := e.Group("/api/todos")
	
	// Add authentication middleware to all todo routes
	todoGroup.Use(authMiddleware.Authenticate)

	// Routes
	todoGroup.POST("", todoHandler.CreateTodo)
	todoGroup.GET("", todoHandler.GetTodos)
	todoGroup.GET("/:id", todoHandler.GetTodo)
	todoGroup.PUT("/:id", todoHandler.UpdateTodo)
	todoGroup.DELETE("/:id", todoHandler.DeleteTodo)
	todoGroup.PATCH("/:id/complete", todoHandler.CompleteTodo)
}
