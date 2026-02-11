package migrate

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// setupMockDB creates a new mock database for testing
// Returns the mock DB, mock controller, and a cleanup function
func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	t.Helper() // Mark this as a test helper

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, mock, cleanup
}

// setupMockDBWithMigrations creates a mock DB and temporary directory with test migrations
func setupMockDBWithMigrations(t *testing.T, migrations []TestMigration) (*sql.DB, sqlmock.Sqlmock, string, func()) {
	t.Helper()

	// Create mock DB
	db, mock, dbCleanup := setupMockDB(t)

	// Create temp directory
	tmpDir := t.TempDir()

	// Create migration files
	for _, m := range migrations {
		// Up migration
		upFile := filepath.Join(tmpDir, m.UpFilename)
		if err := os.WriteFile(upFile, []byte(m.UpSQL), 0644); err != nil {
			dbCleanup()
			t.Fatalf("Failed to create up migration %s: %v", upFile, err)
		}

		// Down migration
		downFile := filepath.Join(tmpDir, m.DownFilename)
		if err := os.WriteFile(downFile, []byte(m.DownSQL), 0644); err != nil {
			dbCleanup()
			t.Fatalf("Failed to create down migration %s: %v", downFile, err)
		}
	}

	cleanup := func() {
		dbCleanup()
		// tmpDir is cleaned up automatically by t.TempDir()
	}

	return db, mock, tmpDir, cleanup
}

// TestMigration represents a migration for testing
type TestMigration struct {
	Version      int
	Name         string
	UpFilename   string
	DownFilename string
	UpSQL        string
	DownSQL      string
}

// NewTestMigration creates a new test migration with standard naming
func NewTestMigration(version int, name, upSQL, downSQL string) TestMigration {
	return TestMigration{
		Version:      version,
		Name:         name,
		UpFilename:   formatMigrationFilename(version, name, "up"),
		DownFilename: formatMigrationFilename(version, name, "down"),
		UpSQL:        upSQL,
		DownSQL:      downSQL,
	}
}

// formatMigrationFilename formats a migration filename
func formatMigrationFilename(version int, name, direction string) string {
	return fmt.Sprintf("%03d_%s.%s.sql", version, name, direction)
}

// expectSchemaTable adds expectations for schema_migrations table operations
func expectSchemaTable(mock sqlmock.Sqlmock) {
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").
		WillReturnResult(sqlmock.NewResult(0, 0))
}

// expectGetAppliedMigrations adds expectations for getting applied migrations
func expectGetAppliedMigrations(mock sqlmock.Sqlmock, versions ...int) {
	expectSchemaTable(mock)

	rows := sqlmock.NewRows([]string{"version"})
	for _, v := range versions {
		rows.AddRow(v)
	}
	mock.ExpectQuery("SELECT version FROM schema_migrations").
		WillReturnRows(rows)
}

// expectMigrationUp adds expectations for a successful up migration
func expectMigrationUp(mock sqlmock.Sqlmock, version int, sql string) {
	mock.ExpectBegin()
	mock.ExpectExec(sql).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO schema_migrations").
		WithArgs(version).
		WillReturnResult(sqlmock.NewResult(int64(version), 1))
	mock.ExpectCommit()
}

// expectMigrationDown adds expectations for a successful down migration
func expectMigrationDown(mock sqlmock.Sqlmock, version int, sql string) {
	mock.ExpectBegin()
	mock.ExpectExec(sql).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM schema_migrations WHERE version").
		WithArgs(version).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
}

// assertMockExpectations verifies all mock expectations were met
func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled mock expectations: %v", err)
	}
}
