package config

import (
	"strings"
	"testing"
)

func TestGetSafeDSN_MasksPassword(t *testing.T) {
	cfg := &Config{
		DatabaseType: "postgres",
		Postgres: DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "testuser",
			Password: "SuperSecretPassword123!",
			DBName:   "testdb",
			SSLMode:  "disable",
		},
	}

	safeDSN := cfg.GetSafeDSN()

	// Password should be masked
	if strings.Contains(safeDSN, "SuperSecretPassword123!") {
		t.Error("GetSafeDSN() contains plaintext password - security vulnerability!")
	}

	// Should contain masked password
	if !strings.Contains(safeDSN, "password=***") {
		t.Errorf("GetSafeDSN() should contain 'password=***', got: %s", safeDSN)
	}

	// Should still contain other information
	expectedParts := []string{
		"host=localhost",
		"port=5432",
		"user=testuser",
		"dbname=testdb",
		"sslmode=disable",
	}

	for _, part := range expectedParts {
		if !strings.Contains(safeDSN, part) {
			t.Errorf("GetSafeDSN() missing expected part '%s', got: %s", part, safeDSN)
		}
	}
}

func TestGetDSN_ContainsPassword(t *testing.T) {
	cfg := &Config{
		DatabaseType: "postgres",
		Postgres: DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "testuser",
			Password: "RealPassword",
			DBName:   "testdb",
			SSLMode:  "disable",
		},
	}

	dsn := cfg.GetDSN()

	// Real DSN should contain actual password (for database connection)
	if !strings.Contains(dsn, "password=RealPassword") {
		t.Errorf("GetDSN() should contain real password for database connection, got: %s", dsn)
	}
}

func TestGetSafeDSN_vs_GetDSN(t *testing.T) {
	cfg := &Config{
		DatabaseType: "postgres",
		Postgres: DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "admin",
			Password: "VerySecretPassword",
			DBName:   "production_db",
			SSLMode:  "require",
		},
	}

	dsn := cfg.GetDSN()
	safeDSN := cfg.GetSafeDSN()

	// They should be different
	if dsn == safeDSN {
		t.Error("GetDSN() and GetSafeDSN() should return different values")
	}

	// Real DSN has password
	if !strings.Contains(dsn, "VerySecretPassword") {
		t.Error("GetDSN() should contain real password")
	}

	// Safe DSN masks password
	if strings.Contains(safeDSN, "VerySecretPassword") {
		t.Error("GetSafeDSN() should NOT contain real password")
	}

	// Safe DSN has masked password
	if !strings.Contains(safeDSN, "password=***") {
		t.Error("GetSafeDSN() should have masked password")
	}
}

func TestGetSafeDSN_EmptyPassword(t *testing.T) {
	cfg := &Config{
		DatabaseType: "postgres",
		Postgres: DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "testuser",
			Password: "", // Empty password
			DBName:   "testdb",
			SSLMode:  "disable",
		},
	}

	safeDSN := cfg.GetSafeDSN()

	// Even empty password should be masked
	if !strings.Contains(safeDSN, "password=***") {
		t.Errorf("GetSafeDSN() should mask even empty password, got: %s", safeDSN)
	}
}

func TestGetSafeDSN_SQLite(t *testing.T) {
	cfg := &Config{
		DatabaseType: "sqlite",
		DatabasePath: "/path/to/database.db",
	}

	safeDSN := cfg.GetSafeDSN()

	// For SQLite, should just return the path
	if safeDSN != "/path/to/database.db" {
		t.Errorf("GetSafeDSN() for SQLite should return database path, got: %s", safeDSN)
	}
}
