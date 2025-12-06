package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Config holds database configuration
type Config struct {
	Adapter  string
	Database string
	Host     string
	Port     string
	Username string
	Password string
	SSLMode  string
}

// Connect establishes a database connection
func Connect(config Config) error {
	var dsn string
	var driver string

	switch config.Adapter {
	case "postgresql", "postgres":
		driver = "postgres"
		sslMode := config.SSLMode
		if sslMode == "" {
			sslMode = "disable"
		}
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.Username, config.Password, config.Database, sslMode)

	case "mysql":
		driver = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			config.Username, config.Password, config.Host, config.Port, config.Database)

	case "sqlite", "sqlite3":
		driver = "sqlite3"
		dsn = config.Database

	default:
		return fmt.Errorf("unsupported database adapter: %s", config.Adapter)
	}

	var err error
	db, err = sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// DB returns the database connection
func DB() *sql.DB {
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// IsConnected checks if database is connected
func IsConnected() bool {
	if db == nil {
		return false
	}
	return db.Ping() == nil
}
