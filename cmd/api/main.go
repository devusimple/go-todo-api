package main

import (
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"todo-api/internal/config"
	"todo-api/internal/infrastructure/repository/postgres"
	"todo-api/internal/interface/api/router"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := postgres.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate database schemas
	if err := postgres.AutoMigrate(db); err != nil {
		logger.Fatalf("Failed to migrate database schemas: %v", err)
	}

	// Initialize Echo framework
	e := echo.New()

	// Set up middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Set up routes
	router.SetupRoutes(e, db, cfg, logger)

	// Start server
	serverAddr := fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port)
	log.Fatal(e.Start(serverAddr))
}
