package config

import "os"

// EmailConfig holds email service configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// Config holds application configuration
type Config struct {
	ServerPort string
	Email      EmailConfig
	
	// Database connection (to fetch users/orders)
	DatabaseDSN string
	
	// Mock mode (for development)
	MockEmail bool
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8083"),
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnv("SMTP_PORT", "587"),
			SMTPUser:     getEnv("SMTP_USER", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("FROM_EMAIL", "noreply@stock-hub.com"),
			FromName:     getEnv("FROM_NAME", "Stock Hub"),
		},
		DatabaseDSN: getEnv("DATABASE_DSN", "host=localhost port=5432 user=postgres password=postgres dbname=stock_hub sslmode=disable"),
		MockEmail:   getEnv("MOCK_EMAIL", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
