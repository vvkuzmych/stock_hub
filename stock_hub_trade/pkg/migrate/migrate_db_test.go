package migrate

import (
	"fmt"
	"testing"
)

// TestEnsureMigrationsTable tests that the migrations tracking table is created
func TestEnsureMigrationsTable(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	expectSchemaTable(mock)

	migrator := NewMigrator(db, "test_migrations")
	err := migrator.ensureMigrationsTable()
	if err != nil {
		t.Errorf("ensureMigrationsTable() failed: %v", err)
	}

	assertMockExpectations(t, mock)
}

// TestGetAppliedMigrations tests retrieving applied migrations
func TestGetAppliedMigrations(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	expectGetAppliedMigrations(mock, 1, 2, 3)

	migrator := NewMigrator(db, "test_migrations")
	applied, err := migrator.GetAppliedMigrations()
	if err != nil {
		t.Fatalf("GetAppliedMigrations() failed: %v", err)
	}

	expected := map[int]bool{1: true, 2: true, 3: true}
	if len(applied) != len(expected) {
		t.Errorf("Expected %d applied migrations, got %d", len(expected), len(applied))
	}

	for version := range expected {
		if !applied[version] {
			t.Errorf("Expected version %d to be applied", version)
		}
	}

	assertMockExpectations(t, mock)
}

// TestUpMigration tests running up migrations with mock DB
func TestUpMigration(t *testing.T) {
	migrations := []TestMigration{
		NewTestMigration(1, "create_test_table",
			"CREATE TABLE test_table (id INTEGER PRIMARY KEY);",
			"DROP TABLE test_table;"),
	}

	db, mock, tmpDir, cleanup := setupMockDBWithMigrations(t, migrations)
	defer cleanup()

	expectGetAppliedMigrations(mock) // No applied migrations
	expectMigrationUp(mock, 1, "CREATE TABLE test_table")

	migrator := NewMigrator(db, tmpDir)
	err := migrator.Up()
	if err != nil {
		t.Errorf("Up() failed: %v", err)
	}

	assertMockExpectations(t, mock)
}

// TestDownMigration tests rolling back migrations with mock DB
func TestDownMigration(t *testing.T) {
	migrations := []TestMigration{
		NewTestMigration(1, "create_test_table",
			"CREATE TABLE test_table (id INTEGER PRIMARY KEY);",
			"DROP TABLE test_table;"),
	}

	db, mock, tmpDir, cleanup := setupMockDBWithMigrations(t, migrations)
	defer cleanup()

	expectGetAppliedMigrations(mock, 1) // Migration 1 is applied
	expectMigrationDown(mock, 1, "DROP TABLE test_table")

	migrator := NewMigrator(db, tmpDir)
	err := migrator.Down()
	if err != nil {
		t.Errorf("Down() failed: %v", err)
	}

	assertMockExpectations(t, mock)
}

// TestUpMigrationAlreadyApplied tests that already-applied migrations are skipped
func TestUpMigrationAlreadyApplied(t *testing.T) {
	migrations := []TestMigration{
		NewTestMigration(1, "test", "CREATE TABLE test;", "DROP TABLE test;"),
	}

	db, mock, tmpDir, cleanup := setupMockDBWithMigrations(t, migrations)
	defer cleanup()

	expectGetAppliedMigrations(mock, 1) // Migration 1 already applied

	// No transaction should be started (migration already applied)

	migrator := NewMigrator(db, tmpDir)
	err := migrator.Up()
	if err != nil {
		t.Errorf("Up() failed: %v", err)
	}

	assertMockExpectations(t, mock)
}

// TestMigrationFailureRollback tests that failed migrations rollback properly
func TestMigrationFailureRollback(t *testing.T) {
	migrations := []TestMigration{
		NewTestMigration(1, "invalid", "THIS IS INVALID SQL;", "-- down"),
	}

	db, mock, tmpDir, cleanup := setupMockDBWithMigrations(t, migrations)
	defer cleanup()

	expectGetAppliedMigrations(mock) // No applied migrations

	mock.ExpectBegin()
	mock.ExpectExec("THIS IS INVALID SQL").
		WillReturnError(fmt.Errorf("syntax error"))
	mock.ExpectRollback()

	migrator := NewMigrator(db, tmpDir)
	err := migrator.Up()
	if err == nil {
		t.Error("Expected Up() to fail with invalid SQL, but it succeeded")
	}

	assertMockExpectations(t, mock)
}

// TestDownMigrationNotApplied tests rolling back when no migrations are applied
func TestDownMigrationNotApplied(t *testing.T) {
	migrations := []TestMigration{
		NewTestMigration(1, "test", "CREATE TABLE test;", "DROP TABLE test;"),
	}

	db, mock, tmpDir, cleanup := setupMockDBWithMigrations(t, migrations)
	defer cleanup()

	expectGetAppliedMigrations(mock) // No migrations applied

	migrator := NewMigrator(db, tmpDir)
	err := migrator.Down()
	if err == nil {
		t.Error("Expected Down() to fail when no migrations applied, but it succeeded")
	}

	assertMockExpectations(t, mock)
}

// TestMultipleMigrationsUp tests applying multiple migrations in order
func TestMultipleMigrationsUp(t *testing.T) {
	migrations := []TestMigration{
		NewTestMigration(1, "create_users", "CREATE TABLE users (id INTEGER PRIMARY KEY);", "-- down"),
		NewTestMigration(2, "create_orders", "CREATE TABLE orders (id INTEGER PRIMARY KEY);", "-- down"),
		NewTestMigration(3, "create_products", "CREATE TABLE products (id INTEGER PRIMARY KEY);", "-- down"),
	}

	db, mock, tmpDir, cleanup := setupMockDBWithMigrations(t, migrations)
	defer cleanup()

	expectGetAppliedMigrations(mock) // No migrations applied

	// Expect each migration in order
	expectMigrationUp(mock, 1, "CREATE TABLE users")
	expectMigrationUp(mock, 2, "CREATE TABLE orders")
	expectMigrationUp(mock, 3, "CREATE TABLE products")

	migrator := NewMigrator(db, tmpDir)
	err := migrator.Up()
	if err != nil {
		t.Errorf("Up() failed: %v", err)
	}

	assertMockExpectations(t, mock)
}
