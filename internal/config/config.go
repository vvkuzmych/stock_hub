package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	ServerPort    string
	WebSocketPath string
	DatabasePath  string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8082"),
		WebSocketPath: getEnv("WS_PATH", "/ws"),
		DatabasePath:  getEnv("DATABASE_PATH", "stock_hub.db"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
