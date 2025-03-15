package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"todo-api/internal/config"
	"todo-api/internal/domain/entity"
)

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(config config.DatabaseConfig) (*gorm.DB, error) {
	dsn := config.GetDSN()
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Todo{},
	)
}
