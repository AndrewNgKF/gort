package gort

import (
	"github.com/AndrewNgKF/gort/internal/database"
)

// DB provides access to database operations
var DB = &dbWrapper{}

type dbWrapper struct{}

// DBConfig holds database configuration
type DBConfig struct {
	Adapter  string
	Database string
	Host     string
	Port     string
	Username string
	Password string
	SSLMode  string
}

// Model creates a new model query for the given table
func (d *dbWrapper) Model(tableName string) *database.Model {
	return database.NewModel(tableName)
}

// Connect connects to the database
func (d *dbWrapper) Connect(config DBConfig) error {
	dbConf := database.Config{
		Adapter:  config.Adapter,
		Database: config.Database,
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		SSLMode:  config.SSLMode,
	}
	return database.Connect(dbConf)
}

// Close closes the database connection
func (d *dbWrapper) Close() error {
	return database.Close()
}

// IsConnected checks if database is connected
func (d *dbWrapper) IsConnected() bool {
	return database.IsConnected()
}

// Query builder shortcuts

// Find finds a record by ID
func Find(tableName string, id interface{}, dest interface{}) error {
	return database.NewModel(tableName).Find(id, dest)
}

// FindAll finds all records
func FindAll(tableName string, dest interface{}) error {
	return database.NewModel(tableName).FindAll(dest)
}

// Create creates a new record
func Create(tableName string, data interface{}) error {
	return database.NewModel(tableName).Create(data)
}

// Update updates a record
func Update(tableName string, id interface{}, data interface{}) error {
	return database.NewModel(tableName).Update(id, data)
}

// Delete deletes a record
func Delete(tableName string, id interface{}) error {
	return database.NewModel(tableName).Delete(id)
}

// Where creates a query builder with a WHERE clause
func Where(tableName string, condition string, args ...interface{}) *database.QueryBuilder {
	return database.NewModel(tableName).Where(condition, args...)
}
