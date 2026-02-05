package config

import (
	"fmt"
)

// PostgresConfig holds PostgreSQL database configuration
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConfigPostgres holds application configuration for PostgreSQL
type ConfigPostgres struct {
	ServerPort    string
	WebSocketPath string
	DatabaseType  string // "sqlite" or "postgres"
	DatabasePath  string // For SQLite
	Postgres      PostgresConfig
}

// LoadPostgres loads configuration from environment variables
func LoadPostgres() *ConfigPostgres {
	dbType := getEnv("DB_TYPE", "sqlite") // Default to SQLite for backwards compatibility

	cfg := &ConfigPostgres{
		ServerPort:    getEnv("SERVER_PORT", "8082"),
		WebSocketPath: getEnv("WS_PATH", "/ws"),
		DatabaseType:  dbType,
	}

	if dbType == "postgres" {
		cfg.Postgres = PostgresConfig{
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
func (c *ConfigPostgres) GetDSN() string {
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
func (c *ConfigPostgres) GetDriver() string {
	if c.DatabaseType == "postgres" {
		return "postgres"
	}
	return "sqlite"
}

// GetMigrationsDir returns the migrations directory path
func (c *ConfigPostgres) GetMigrationsDir() string {
	if c.DatabaseType == "postgres" {
		return "migrations_postgres"
	}
	return "migrations"
}
