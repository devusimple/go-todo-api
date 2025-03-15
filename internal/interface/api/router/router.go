package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"todo-api/internal/config"
	"todo-api/internal/infrastructure/repository/postgres"
	"todo-api/internal/interface/api/middleware"
	"todo-api/internal/interface/api/validator"
	"todo-api/internal/util/jwt"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config, logger *logrus.Logger) {
	// Set custom validator
	e.Validator = validator.NewCustomValidator()

	// Initialize JWT service
	jwtService := jwt.NewJWTService(cfg.JWT.SecretKey, cfg.JWT.Expiration)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	todoRepo := postgres.NewTodoRepository(db)

	// Set up routes
	SetupUserRoutes(e, userRepo, jwtService, authMiddleware, logger)
	SetupTodoRoutes(e, todoRepo, authMiddleware, logger)

	// Set up health check route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	// Set up API version info
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"name":    "Todo API",
			"version": "1.0.0",
		})
	})
}
