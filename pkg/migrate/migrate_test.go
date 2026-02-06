package migrate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrationNameParsing(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()

	testCases := []struct {
		filename     string
		expectedName string
	}{
		{
			filename:     "001_create_messages_table.up.sql",
			expectedName: "create_messages_table",
		},
		{
			filename:     "002_create_users_table.down.sql",
			expectedName: "create_users_table",
		},
		{
			filename:     "003_add_user_email.up.sql",
			expectedName: "add_user_email",
		},
		{
			filename:     "010_create_orders.up.sql",
			expectedName: "create_orders",
		},
		{
			filename:     "100_update_schema_with_long_name.down.sql",
			expectedName: "update_schema_with_long_name",
		},
	}

	// Create test migration files
	for _, tc := range testCases {
		upFile := filepath.Join(tmpDir, tc.filename)
		if err := os.WriteFile(upFile, []byte("-- test migration"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", upFile, err)
		}
	}

	// Create migrator and load migrations
	migrator := NewMigrator(nil, tmpDir)
	migrations, err := migrator.LoadMigrations()
	if err != nil {
		t.Fatalf("Failed to load migrations: %v", err)
	}

	// Verify migration names
	if len(migrations) != len(testCases) {
		t.Fatalf("Expected %d migrations, got %d", len(testCases), len(migrations))
	}

	for i, migration := range migrations {
		expected := testCases[i].expectedName
		if migration.Name != expected {
			t.Errorf("Migration %d: expected name %q, got %q (from file %s)",
				i, expected, migration.Name, testCases[i].filename)
		}
	}
}

func TestMigrationNameParsingEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()

	edgeCases := []struct {
		filename     string
		expectedName string
		description  string
	}{
		{
			filename:     "001_simple.up.sql",
			expectedName: "simple",
			description:  "Single word migration name",
		},
		{
			filename:     "002_add_column_to_users_table_v2.up.sql",
			expectedName: "add_column_to_users_table_v2",
			description:  "Multiple underscores in name",
		},
		{
			filename:     "999_final_migration.down.sql",
			expectedName: "final_migration",
			description:  "Three-digit version number",
		},
	}

	// Create test files
	for _, tc := range edgeCases {
		file := filepath.Join(tmpDir, tc.filename)
		if err := os.WriteFile(file, []byte("-- "+tc.description), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Load and verify
	migrator := NewMigrator(nil, tmpDir)
	migrations, err := migrator.LoadMigrations()
	if err != nil {
		t.Fatalf("Failed to load migrations: %v", err)
	}

	for i, migration := range migrations {
		expected := edgeCases[i].expectedName
		if migration.Name != expected {
			t.Errorf("%s: expected name %q, got %q",
				edgeCases[i].description, expected, migration.Name)
		}
	}
}

func TestMigrationPairing(t *testing.T) {
	tmpDir := t.TempDir()

	// Create matching up and down migrations
	files := []string{
		"001_create_table.up.sql",
		"001_create_table.down.sql",
		"002_add_column.up.sql",
		"002_add_column.down.sql",
	}

	for _, file := range files {
		path := filepath.Join(tmpDir, file)
		if err := os.WriteFile(path, []byte("-- migration"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	migrator := NewMigrator(nil, tmpDir)
	migrations, err := migrator.LoadMigrations()
	if err != nil {
		t.Fatalf("Failed to load migrations: %v", err)
	}

	// Should have 2 migrations (paired up+down)
	if len(migrations) != 2 {
		t.Fatalf("Expected 2 migrations, got %d", len(migrations))
	}

	// Verify both have up and down
	for _, migration := range migrations {
		if migration.Up == "" {
			t.Errorf("Migration %d (%s) missing up migration", migration.Version, migration.Name)
		}
		if migration.Down == "" {
			t.Errorf("Migration %d (%s) missing down migration", migration.Version, migration.Name)
		}
	}

	// Verify names
	expectedNames := []string{"create_table", "add_column"}
	for i, migration := range migrations {
		if migration.Name != expectedNames[i] {
			t.Errorf("Migration %d: expected name %q, got %q",
				i, expectedNames[i], migration.Name)
		}
	}
}
