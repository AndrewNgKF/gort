package commands

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AndrewNgKF/gort/internal/migration"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations",
	Long:  "Run database migrations, rollbacks, and schema dumps",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the database",
	Long:  "Create the database from config/database.yml",
	Run: func(cmd *cobra.Command, args []string) {
		if err := createDatabase(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error creating database: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Database created successfully!")
	},
}

var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop the database",
	Long:  "Drop the database (WARNING: This will delete all data!)",
	Run: func(cmd *cobra.Command, args []string) {
		if err := dropDatabase(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error dropping database: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Database dropped successfully!")
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Create database and run migrations",
	Long:  "Create the database, run all migrations, and load seed data",
	Run: func(cmd *cobra.Command, args []string) {
		// Create database
		fmt.Print("Creating database...")
		if err := createDatabase(); err != nil {
			fmt.Printf(" ❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(" ✅")

		// Run migrations
		fmt.Println("\nRunning migrations...")
		migrateCmd.Run(cmd, args)

		fmt.Println("\n✅ Database setup complete!")
	},
}

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Aliases: []string{"migrate"},
	Short:   "Run pending migrations",
	Long:    "Run all pending database migrations and update schema.go",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := connectDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		mgr := migration.NewManager(db)

		// Ensure migrations table exists
		if err := mgr.EnsureMigrationsTable(); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating migrations table: %v\n", err)
			os.Exit(1)
		}

		// Get pending migrations
		pending, err := mgr.GetPendingMigrations()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting pending migrations: %v\n", err)
			os.Exit(1)
		}

		if len(pending) == 0 {
			fmt.Println("No pending migrations")
			return
		}

		fmt.Printf("Running %d migration(s)...\n", len(pending))

		// Run each pending migration
		for _, version := range pending {
			fmt.Printf("  Running migration %s...", version)

			// Find and run the migration file
			migrationFile := filepath.Join("db", "migrations", version+"_*.go")
			matches, err := filepath.Glob(migrationFile)
			if err != nil || len(matches) == 0 {
				fmt.Println(" FAILED (file not found)")
				continue
			}

			// Execute the Up function
			if err := runMigrationUp(db, version); err != nil {
				fmt.Printf(" FAILED: %v\n", err)
				os.Exit(1)
			}

			// Record migration
			if err := mgr.RecordMigration(version); err != nil {
				fmt.Printf(" FAILED to record: %v\n", err)
				os.Exit(1)
			}

			fmt.Println(" OK")
		}

		// Dump schema
		fmt.Print("Updating schema.rb...")
		if err := mgr.DumpSchema(); err != nil {
			fmt.Printf(" FAILED: %v\n", err)
		} else {
			fmt.Println(" OK")
		}
	},
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback the last migration",
	Long:  "Rollback the most recently applied migration and update schema.go",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := connectDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		mgr := migration.NewManager(db)

		// Get applied migrations
		applied, err := mgr.GetAppliedMigrations()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting applied migrations: %v\n", err)
			os.Exit(1)
		}

		if len(applied) == 0 {
			fmt.Println("No migrations to rollback")
			return
		}

		// Get the last migration
		lastVersion := applied[len(applied)-1]
		fmt.Printf("Rolling back migration %s...", lastVersion)

		// Execute the Down function
		if err := runMigrationDown(db, lastVersion); err != nil {
			fmt.Printf(" FAILED: %v\n", err)
			os.Exit(1)
		}

		// Remove migration record
		if err := mgr.RemoveMigration(lastVersion); err != nil {
			fmt.Printf(" FAILED to remove record: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(" OK")

		// Dump schema
		fmt.Print("Updating schema.go...")
		if err := mgr.DumpSchema(); err != nil {
			fmt.Printf(" FAILED: %v\n", err)
		} else {
			fmt.Println(" OK")
		}
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the database",
	Long:  "Drop all tables and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := connectDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		mgr := migration.NewManager(db)

		fmt.Print("Dropping all tables...")
		// Get all applied migrations in reverse order
		applied, err := mgr.GetAppliedMigrations()
		if err != nil {
			fmt.Printf(" FAILED: %v\n", err)
			os.Exit(1)
		}

		// Rollback all migrations
		for i := len(applied) - 1; i >= 0; i-- {
			version := applied[i]
			if err := runMigrationDown(db, version); err != nil {
				fmt.Printf(" FAILED on %s: %v\n", version, err)
				os.Exit(1)
			}
			if err := mgr.RemoveMigration(version); err != nil {
				fmt.Printf(" FAILED to remove %s: %v\n", version, err)
				os.Exit(1)
			}
		}
		fmt.Println(" OK")

		// Re-run all migrations
		fmt.Println("Re-running all migrations...")
		migrateCmd.Run(cmd, args)
	},
}

var schemaCmd = &cobra.Command{
	Use:   "schema:dump",
	Short: "Dump the database schema",
	Long:  "Export the current database schema to db/schema.go",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := connectDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		mgr := migration.NewManager(db)

		fmt.Print("Dumping schema...")
		if err := mgr.DumpSchema(); err != nil {
			fmt.Printf(" FAILED: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(" OK")
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(createCmd)
	dbCmd.AddCommand(dropCmd)
	dbCmd.AddCommand(setupCmd)
	dbCmd.AddCommand(migrateCmd)
	dbCmd.AddCommand(rollbackCmd)
	dbCmd.AddCommand(resetCmd)
	dbCmd.AddCommand(schemaCmd)

	// Add Rails-style aliases with colons
	rootCmd.AddCommand(&cobra.Command{
		Use:    "db:create",
		Hidden: true,
		Run:    createCmd.Run,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:    "db:drop",
		Hidden: true,
		Run:    dropCmd.Run,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:    "db:setup",
		Hidden: true,
		Run:    setupCmd.Run,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:    "db:migrate",
		Hidden: true,
		Run:    migrateCmd.Run,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:    "db:rollback",
		Hidden: true,
		Run:    rollbackCmd.Run,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:    "db:reset",
		Hidden: true,
		Run:    resetCmd.Run,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:    "db:schema:dump",
		Hidden: true,
		Run:    schemaCmd.Run,
	})
}

// connectDB connects to the database using config/database.yml
func connectDB() (*sql.DB, error) {
	// Read database config
	configPath := filepath.Join("config", "database.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read database config: %v", err)
	}

	// Parse YAML
	var config map[string]map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse database config: %v", err)
	}

	// Get environment (default to development)
	env := os.Getenv("GORT_ENV")
	if env == "" {
		env = "development"
	}

	dbConfig, ok := config[env]
	if !ok {
		return nil, fmt.Errorf("database config not found for environment: %s", env)
	}

	adapter := dbConfig["adapter"].(string)
	var dsn string

	switch adapter {
	case "postgresql":
		host := dbConfig["host"].(string)
		port := int(dbConfig["port"].(int))
		database := dbConfig["database"].(string)
		username := dbConfig["username"].(string)
		password := dbConfig["password"].(string)
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, username, password, database)
		return sql.Open("postgres", dsn)

	case "mysql":
		host := dbConfig["host"].(string)
		port := int(dbConfig["port"].(int))
		database := dbConfig["database"].(string)
		username := dbConfig["username"].(string)
		password := dbConfig["password"].(string)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			username, password, host, port, database)
		return sql.Open("mysql", dsn)

	case "sqlite3":
		database := dbConfig["database"].(string)
		return sql.Open("sqlite3", database)

	default:
		return nil, fmt.Errorf("unsupported database adapter: %s", adapter)
	}
}

// runMigrationUp runs the Up function for a migration
func runMigrationUp(db *sql.DB, version string) error {
	// This is a simplified version - in reality, you'd need to:
	// 1. Load the migration file
	// 2. Execute the Up function
	// For now, we'll read and execute SQL directly

	migrationFile := filepath.Join("db", "migrations", version+"_*.go")
	matches, err := filepath.Glob(migrationFile)
	if err != nil || len(matches) == 0 {
		return fmt.Errorf("migration file not found")
	}

	// Read the migration file
	content, err := os.ReadFile(matches[0])
	if err != nil {
		return err
	}

	// Extract SQL from Up function (this is hacky but works for now)
	// In a real implementation, you'd compile and run the Go code
	sqlQuery := extractSQL(string(content), "Up_"+version)
	if sqlQuery == "" {
		return fmt.Errorf("could not extract SQL from Up function")
	}

	_, err = db.Exec(sqlQuery)
	return err
}

// runMigrationDown runs the Down function for a migration
func runMigrationDown(db *sql.DB, version string) error {
	migrationFile := filepath.Join("db", "migrations", version+"_*.go")
	matches, err := filepath.Glob(migrationFile)
	if err != nil || len(matches) == 0 {
		return fmt.Errorf("migration file not found")
	}

	content, err := os.ReadFile(matches[0])
	if err != nil {
		return err
	}

	sqlQuery := extractSQL(string(content), "Down_"+version)
	if sqlQuery == "" {
		return fmt.Errorf("could not extract SQL from Down function")
	}

	_, err = db.Exec(sqlQuery)
	return err
}

// extractSQL extracts SQL from a migration function
func extractSQL(content, functionName string) string {
	// Find the function
	startMarker := "func " + functionName
	startIdx := indexOf(content, startMarker)
	if startIdx == -1 {
		return ""
	}

	// Find the SQL string (between backticks)
	sqlStart := indexOf(content[startIdx:], "`")
	if sqlStart == -1 {
		return ""
	}
	sqlStart += startIdx + 1

	sqlEnd := indexOf(content[sqlStart:], "`")
	if sqlEnd == -1 {
		return ""
	}
	sqlEnd += sqlStart

	return content[sqlStart:sqlEnd]
}

func indexOf(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// createDatabase creates the database
func createDatabase() error {
	config, env, err := loadDBConfig()
	if err != nil {
		return err
	}

	dbConfig := config[env]
	adapter := dbConfig["adapter"].(string)
	database := dbConfig["database"].(string)

	switch adapter {
	case "postgresql":
		host := dbConfig["host"].(string)
		port := int(dbConfig["port"].(int))
		username := dbConfig["username"].(string)
		password := dbConfig["password"].(string)

		// Connect to postgres database to create the target database
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
			host, port, username, password)
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return err
		}
		defer db.Close()

		// Create database
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", database))
		if err != nil {
			// Check if database already exists
			if indexOf(err.Error(), "already exists") != -1 {
				return fmt.Errorf("database %s already exists", database)
			}
			return err
		}

	case "mysql":
		host := dbConfig["host"].(string)
		port := int(dbConfig["port"].(int))
		username := dbConfig["username"].(string)
		password := dbConfig["password"].(string)

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
			username, password, host, port)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		defer db.Close()

		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", database))
		if err != nil {
			return err
		}

	case "sqlite3":
		// SQLite creates the database file automatically
		return nil

	default:
		return fmt.Errorf("unsupported database adapter: %s", adapter)
	}

	return nil
}

// dropDatabase drops the database
func dropDatabase() error {
	config, env, err := loadDBConfig()
	if err != nil {
		return err
	}

	dbConfig := config[env]
	adapter := dbConfig["adapter"].(string)
	database := dbConfig["database"].(string)

	switch adapter {
	case "postgresql":
		host := dbConfig["host"].(string)
		port := int(dbConfig["port"].(int))
		username := dbConfig["username"].(string)
		password := dbConfig["password"].(string)

		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
			host, port, username, password)
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return err
		}
		defer db.Close()

		_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", database))
		return err

	case "mysql":
		host := dbConfig["host"].(string)
		port := int(dbConfig["port"].(int))
		username := dbConfig["username"].(string)
		password := dbConfig["password"].(string)

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
			username, password, host, port)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		defer db.Close()

		_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", database))
		return err

	case "sqlite3":
		// Delete the SQLite file
		return os.Remove(database)

	default:
		return fmt.Errorf("unsupported database adapter: %s", adapter)
	}
}

// loadDBConfig loads the database configuration
func loadDBConfig() (map[string]map[string]interface{}, string, error) {
	configPath := filepath.Join("config", "database.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read database config: %v", err)
	}

	var config map[string]map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, "", fmt.Errorf("failed to parse database config: %v", err)
	}

	env := os.Getenv("GORT_ENV")
	if env == "" {
		env = "development"
	}

	if _, ok := config[env]; !ok {
		return nil, "", fmt.Errorf("database config not found for environment: %s", env)
	}

	return config, env, nil
}
