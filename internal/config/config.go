package config

import (
	"fmt"
	"os"
)

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Config holds application configuration
type Config struct {
	ServerPort    string
	WebSocketPath string
	DatabaseType  string // "sqlite" or "postgres"
	DatabasePath  string // For SQLite
	Postgres      DBConfig
}

// Load loads configuration from environment variables
func Load() *Config {
	dbType := getEnv("DB_TYPE", "postgres") // Default to PostgreSQL

	cfg := &Config{
		ServerPort:    getEnv("SERVER_PORT", "8082"),
		WebSocketPath: getEnv("WS_PATH", "/ws"),
		DatabaseType:  dbType,
	}

	if dbType == "postgres" {
		cfg.Postgres = DBConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnv("POSTGRES_DB", "stock_hub"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		}
	} else {
		cfg.DatabasePath = getEnv("DATABASE_PATH", "stock_hub.db")
	}

	return cfg
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	if c.DatabaseType == "postgres" {
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			c.Postgres.Host,
			c.Postgres.Port,
			c.Postgres.User,
			c.Postgres.Password,
			c.Postgres.DBName,
			c.Postgres.SSLMode,
		)
	}
	return c.DatabasePath
}

// GetDriver returns the database driver name
func (c *Config) GetDriver() string {
	if c.DatabaseType == "postgres" {
		return "postgres"
	}
	return "sqlite"
}

// GetMigrationsDir returns the migrations directory path
func (c *Config) GetMigrationsDir() string {
	return "migrations"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
