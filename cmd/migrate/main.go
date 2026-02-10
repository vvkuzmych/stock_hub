package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"stock_hub/internal/config"
	"stock_hub/pkg/migrate"

	_ "github.com/lib/pq"
)

func main() {
	var (
		command = flag.String("command", "up", "Migration command: up or down")
		dbPath  = flag.String("db", "", "Database path (defaults to config)")
	)
	flag.Parse()

	// Load config
	cfg := config.Load()

	// Get DSN and driver
	dsn := cfg.GetDSN()
	driver := cfg.GetDriver()
	migrationsDir := cfg.GetMigrationsDir()

	// Override DSN if dbPath is provided
	if *dbPath != "" {
		dsn = *dbPath
	}

	log.Printf("Database type: %s", cfg.DatabaseType)
	log.Printf("Migrations directory: %s", migrationsDir)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	migrator := migrate.NewMigrator(db, migrationsDir)

	switch *command {
	case "up":
		if err := migrator.Up(); err != nil {
			log.Fatalf("Failed to run migrations up: %v", err)
		}
		fmt.Println("Migrations applied successfully")
	case "down":
		if err := migrator.Down(); err != nil {
			log.Fatalf("Failed to run migrations down: %v", err)
		}
		fmt.Println("Last migration rolled back successfully")
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s. Use 'up' or 'down'\n", *command)
		os.Exit(1)
	}
}
