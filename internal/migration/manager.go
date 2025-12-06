package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Manager handles migration operations
type Manager struct {
	db *sql.DB
}

// NewManager creates a new migration manager
func NewManager(db *sql.DB) *Manager {
	return &Manager{db: db}
}

// EnsureMigrationsTable creates the schema_migrations table if it doesn't exist
func (m *Manager) EnsureMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := m.db.Exec(query)
	return err
}

// GetAppliedMigrations returns a list of applied migration versions
func (m *Manager) GetAppliedMigrations() ([]string, error) {
	query := "SELECT version FROM schema_migrations ORDER BY version"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, nil
}

// GetPendingMigrations returns migrations that haven't been run yet
func (m *Manager) GetPendingMigrations() ([]string, error) {
	// Get all migration files
	files, err := filepath.Glob("db/migrations/*.go")
	if err != nil {
		return nil, err
	}

	// Get applied migrations
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}

	appliedMap := make(map[string]bool)
	for _, version := range applied {
		appliedMap[version] = true
	}

	// Find pending migrations
	var pending []string
	for _, file := range files {
		base := filepath.Base(file)
		// Extract timestamp from filename (first 14 characters)
		if len(base) >= 14 {
			version := base[:14]
			if !appliedMap[version] {
				pending = append(pending, version)
			}
		}
	}

	sort.Strings(pending)
	return pending, nil
}

// RecordMigration records a migration as applied
func (m *Manager) RecordMigration(version string) error {
	query := "INSERT INTO schema_migrations (version) VALUES ($1)"
	_, err := m.db.Exec(query, version)
	return err
}

// RemoveMigration removes a migration record
func (m *Manager) RemoveMigration(version string) error {
	query := "DELETE FROM schema_migrations WHERE version = $1"
	_, err := m.db.Exec(query, version)
	return err
}

// DumpSchema dumps the current database schema to db/schema.go
func (m *Manager) DumpSchema() error {
	// Get database driver name
	driverName := "unknown"

	// Try to detect SQLite by querying sqlite_master
	var testQuery string
	var tables []string

	// Try SQLite query first
	rows, err := m.db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name != 'schema_migrations' AND name NOT LIKE 'sqlite_%' ORDER BY name")
	if err == nil {
		// SQLite database
		driverName = "sqlite3"
		defer rows.Close()

		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err != nil {
				return err
			}
			tables = append(tables, tableName)
		}
	} else {
		// Try PostgreSQL/MySQL query
		testQuery = `
			SELECT table_name 
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name != 'schema_migrations'
			ORDER BY table_name
		`
		rows, err = m.db.Query(testQuery)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err != nil {
				return err
			}
			tables = append(tables, tableName)
		}
	}

	// Build schema content
	var schema strings.Builder
	schema.WriteString("# This file is auto-generated from the current state of the database.\n")
	schema.WriteString(fmt.Sprintf("# Generated at: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	for _, table := range tables {
		schema.WriteString(fmt.Sprintf("Table: %s\n", table))
		schema.WriteString(strings.Repeat("-", 60) + "\n")

		if driverName == "sqlite3" {
			// SQLite schema introspection
			columnRows, err := m.db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
			if err != nil {
				continue
			}

			for columnRows.Next() {
				var cid int
				var name, ctype string
				var notNull int
				var dfltValue sql.NullString
				var pk int

				if err := columnRows.Scan(&cid, &name, &ctype, &notNull, &dfltValue, &pk); err != nil {
					continue
				}

				nullable := ""
				if notNull == 1 {
					nullable = " NOT NULL"
				}

				defaultVal := ""
				if dfltValue.Valid {
					defaultVal = fmt.Sprintf(" DEFAULT %s", dfltValue.String)
				}

				primaryKey := ""
				if pk == 1 {
					primaryKey = " PRIMARY KEY"
				}

				schema.WriteString(fmt.Sprintf("  %s: %s%s%s%s\n", name, ctype, nullable, defaultVal, primaryKey))
			}
			columnRows.Close()
		} else {
			// PostgreSQL/MySQL schema introspection
			columnQuery := `
				SELECT column_name, data_type, is_nullable, column_default
				FROM information_schema.columns
				WHERE table_name = $1
				ORDER BY ordinal_position
			`

			columnRows, err := m.db.Query(columnQuery, table)
			if err != nil {
				continue
			}

			for columnRows.Next() {
				var columnName, dataType, isNullable string
				var columnDefault sql.NullString
				if err := columnRows.Scan(&columnName, &dataType, &isNullable, &columnDefault); err != nil {
					continue
				}

				nullable := ""
				if isNullable == "NO" {
					nullable = " NOT NULL"
				}

				defaultVal := ""
				if columnDefault.Valid {
					defaultVal = fmt.Sprintf(" DEFAULT %s", columnDefault.String)
				}

				schema.WriteString(fmt.Sprintf("  %s: %s%s%s\n", columnName, dataType, nullable, defaultVal))
			}
			columnRows.Close()
		}

		schema.WriteString("\n")
	}

	// Write to file
	schemaPath := filepath.Join("db", "schema.go")
	return ioutil.WriteFile(schemaPath, []byte(schema.String()), 0644)
}

// LoadSchema loads the schema from db/schema.go (for test database setup)
func (m *Manager) LoadSchema() error {
	schemaPath := filepath.Join("db", "schema.go")

	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		return fmt.Errorf("schema.go not found. Run migrations first")
	}

	// In a real implementation, this would parse schema.go and execute it
	// For now, we'll just return a message
	return fmt.Errorf("schema loading not yet implemented - run migrations instead")
}
