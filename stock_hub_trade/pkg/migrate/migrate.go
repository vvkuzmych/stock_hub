package migrate

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Up      string
	Down    string
	Name    string
}

// Migrator handles database migrations
type Migrator struct {
	db            *sql.DB
	migrationsDir string
	placeholder   string // "?" for SQLite, "$1" for PostgreSQL
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}
	
	// PostgreSQL uses $1 placeholder
	placeholder := "$1"
	
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
		placeholder:   placeholder,
	}
}

// LoadMigrations loads all migration files from the migrations directory
func (m *Migrator) LoadMigrations() ([]Migration, error) {
	migrations := make(map[int]Migration)

	entries, err := os.ReadDir(m.migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		// Parse filename: 001_create_messages_table.up.sql or 001_create_messages_table.down.sql
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		// Extract name and direction
		// Example: 001_create_messages_table.up.sql
		ext := filepath.Ext(filename)                    // .sql
		nameWithDir := strings.TrimSuffix(filename, ext) // 001_create_messages_table.up
		dirParts := strings.Split(nameWithDir, ".")      // [001_create_messages_table, up]
		if len(dirParts) < 2 {
			continue
		}

		direction := dirParts[len(dirParts)-1] // "up" or "down"
		
		// Extract migration name from base (without version prefix)
		// dirParts[0] = "001_create_messages_table"
		baseName := dirParts[0]
		nameParts := strings.SplitN(baseName, "_", 2) // Split only on first underscore
		var name string
		if len(nameParts) > 1 {
			name = nameParts[1] // "create_messages_table"
		} else {
			name = baseName // Fallback if no underscore found
		}

		// Read migration content
		path := filepath.Join(m.migrationsDir, filename)
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", path, err)
		}

		migration, exists := migrations[version]
		if !exists {
			migration = Migration{
				Version: version,
				Name:    name,
			}
		}

		if direction == "up" {
			migration.Up = string(content)
		} else if direction == "down" {
			migration.Down = string(content)
		}

		migrations[version] = migration
	}

	// Convert map to sorted slice
	result := make([]Migration, 0, len(migrations))
	for _, migration := range migrations {
		result = append(result, migration)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Version < result[j].Version
	})

	return result, nil
}

// ensureMigrationsTable creates the migrations tracking table if it doesn't exist
func (m *Migrator) ensureMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := m.db.Exec(query)
	return err
}

// GetAppliedMigrations returns a set of applied migration versions
func (m *Migrator) GetAppliedMigrations() (map[int]bool, error) {
	if err := m.ensureMigrationsTable(); err != nil {
		return nil, err
	}

	rows, err := m.db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// Up runs all pending up migrations
func (m *Migrator) Up() error {
	migrations, err := m.LoadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, migration := range migrations {
		if applied[migration.Version] {
			continue
		}

		if migration.Up == "" {
			return fmt.Errorf("migration %d (%s) has no up migration", migration.Version, migration.Name)
		}

		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(migration.Up); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d (%s): %w", migration.Version, migration.Name, err)
		}

		query := fmt.Sprintf("INSERT INTO schema_migrations (version) VALUES (%s)", m.placeholder)
		if _, err := tx.Exec(query, migration.Version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}
	}

	return nil
}

// Down runs the last down migration
func (m *Migrator) Down() error {
	migrations, err := m.LoadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	if len(migrations) == 0 {
		return fmt.Errorf("no migrations found")
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Find the highest applied migration
	var lastMigration *Migration
	for i := len(migrations) - 1; i >= 0; i-- {
		if applied[migrations[i].Version] {
			lastMigration = &migrations[i]
			break
		}
	}

	if lastMigration == nil {
		return fmt.Errorf("no migrations to rollback")
	}

	if lastMigration.Down == "" {
		return fmt.Errorf("migration %d (%s) has no down migration", lastMigration.Version, lastMigration.Name)
	}

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if _, err := tx.Exec(lastMigration.Down); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute down migration %d (%s): %w", lastMigration.Version, lastMigration.Name, err)
	}

	query := fmt.Sprintf("DELETE FROM schema_migrations WHERE version = %s", m.placeholder)
	if _, err := tx.Exec(query, lastMigration.Version); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %d: %w", lastMigration.Version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit down migration %d: %w", lastMigration.Version, err)
	}

	return nil
}
