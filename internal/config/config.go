package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig represents the JWT configuration
type JWTConfig struct {
	SecretKey  string
	Expiration time.Duration
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	// Set default values
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8000"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("PGHOST", "localhost"),
			Port:     getEnv("PGPORT", "5432"),
			Username: getEnv("PGUSER", "postgres"),
			Password: getEnv("PGPASSWORD", "postgres"),
			DBName:   getEnv("PGDATABASE", "todo_db"),
			SSLMode:  getEnv("PGSSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey:  getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
			Expiration: time.Duration(getEnvAsInt("JWT_EXPIRATION_HOURS", 24)) * time.Hour,
		},
	}

	// Check for DATABASE_URL (used by some PaaS providers)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		// Override individual database settings if DATABASE_URL is provided
		config.Database.Host = dbURL
		config.Database.Port = ""
		config.Database.Username = ""
		config.Database.Password = ""
		config.Database.DBName = ""
	}

	return config, nil
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	// If using the DATABASE_URL environment variable
	if c.Host != "" && c.Port == "" && c.Username == "" {
		return c.Host
	}

	// Build the DSN from individual parts
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.DBName, c.SSLMode)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
