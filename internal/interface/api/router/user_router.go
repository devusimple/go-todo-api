package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"todo-api/internal/domain/repository"
	"todo-api/internal/domain/usecase"
	"todo-api/internal/interface/api/handler"
	"todo-api/internal/interface/api/middleware"
	"todo-api/internal/util/jwt"
)

// SetupUserRoutes sets up routes related to user operations
func SetupUserRoutes(
	e *echo.Echo,
	userRepo repository.UserRepository,
	jwtService jwt.JWTService,
	authMiddleware *middleware.AuthMiddleware,
	logger *logrus.Logger,
) {
	// Initialize user use case
	userUseCase := usecase.NewUserUseCase(userRepo, jwtService)

	// Initialize user handler
	userHandler := handler.NewUserHandler(userUseCase, logger)

	// Define public user routes
	authGroup := e.Group("/api/auth")
	authGroup.POST("/register", userHandler.Register)
	authGroup.POST("/login", userHandler.Login)

	// Define protected user routes
	userGroup := e.Group("/api/users")
	userGroup.Use(authMiddleware.Authenticate)
	userGroup.GET("/me", userHandler.GetProfile)
	userGroup.PUT("/me", userHandler.UpdateProfile)
	userGroup.PUT("/me/password", userHandler.UpdatePassword)
}
