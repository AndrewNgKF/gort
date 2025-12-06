package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Model is the base model that provides common functionality
type Model struct {
	tableName string
	db        *sql.DB
}

// NewModel creates a new model instance
func NewModel(tableName string) *Model {
	return &Model{
		tableName: tableName,
		db:        DB(),
	}
}

// Table sets the table name
func (m *Model) Table(name string) *Model {
	m.tableName = name
	return m
}

// QueryBuilder builds database queries
type QueryBuilder struct {
	model        *Model
	whereClauses []string
	whereArgs    []interface{}
	orderBy      string
	limit        int
	offset       int
}

// Where adds a WHERE clause
func (m *Model) Where(condition string, args ...interface{}) *QueryBuilder {
	return &QueryBuilder{
		model:        m,
		whereClauses: []string{condition},
		whereArgs:    args,
	}
}

// Find finds a record by ID
func (m *Model) Find(id interface{}, dest interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", m.tableName)

	row := m.db.QueryRow(query, id)
	return scanStruct(row, dest)
}

// FindAll finds all records
func (m *Model) FindAll(dest interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s", m.tableName)

	rows, err := m.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scanSlice(rows, dest)
}

// Create inserts a new record
func (m *Model) Create(data interface{}) error {
	if m.db == nil {
		return fmt.Errorf("database not connected")
	}

	// Set timestamps if they exist
	setTimestamps(data, true)

	fields, values, placeholders := getFieldsAndValues(data)

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		m.tableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "))

	// Execute and get the ID back
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	idField := v.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() {
		var id int64
		err := m.db.QueryRow(query, values...).Scan(&id)
		if err != nil {
			return fmt.Errorf("create failed: %w", err)
		}
		idField.SetInt(id)
		return nil
	}

	_, err := m.db.Exec(query, values...)
	return err
}

// Update updates a record
func (m *Model) Update(id interface{}, data interface{}) error {
	// Update timestamps
	setTimestamps(data, false)

	fields, values, _ := getFieldsAndValues(data)

	setClauses := make([]string, len(fields))
	for i, field := range fields {
		setClauses[i] = fmt.Sprintf("%s = $%d", field, i+1)
	}

	values = append(values, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		m.tableName,
		strings.Join(setClauses, ", "),
		len(values))

	_, err := m.db.Exec(query, values...)
	return err
}

// Delete deletes a record
func (m *Model) Delete(id interface{}) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", m.tableName)
	_, err := m.db.Exec(query, id)
	return err
}

// QueryBuilder methods

// Where adds another WHERE clause
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.whereClauses = append(qb.whereClauses, condition)
	qb.whereArgs = append(qb.whereArgs, args...)
	return qb
}

// Order adds ORDER BY clause
func (qb *QueryBuilder) Order(orderBy string) *QueryBuilder {
	qb.orderBy = orderBy
	return qb
}

// Limit sets the LIMIT
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset sets the OFFSET
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// First finds the first matching record
func (qb *QueryBuilder) First(dest interface{}) error {
	query := qb.buildQuery()
	query += " LIMIT 1"

	row := qb.model.db.QueryRow(query, qb.whereArgs...)
	return scanStruct(row, dest)
}

// All finds all matching records
func (qb *QueryBuilder) All(dest interface{}) error {
	query := qb.buildQuery()

	rows, err := qb.model.db.Query(query, qb.whereArgs...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scanSlice(rows, dest)
}

func (qb *QueryBuilder) buildQuery() string {
	query := fmt.Sprintf("SELECT * FROM %s", qb.model.tableName)

	if len(qb.whereClauses) > 0 {
		query += " WHERE " + strings.Join(qb.whereClauses, " AND ")
	}

	if qb.orderBy != "" {
		query += " ORDER BY " + qb.orderBy
	}

	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}

	if qb.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offset)
	}

	return query
}

// Helper functions

func getFieldsAndValues(data interface{}) ([]string, []interface{}, []string) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	var fields []string
	var values []interface{}
	var placeholders []string

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip ID and unexported fields
		if field.Name == "ID" || !field.IsExported() {
			continue
		}

		// Get field name from tag or use field name
		fieldName := field.Tag.Get("db")
		if fieldName == "" {
			fieldName = toSnakeCase(field.Name)
		}

		if fieldName == "-" {
			continue
		}

		fields = append(fields, fieldName)
		values = append(values, value.Interface())
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(placeholders)+1))
	}

	return fields, values, placeholders
}

func scanStruct(row *sql.Row, dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}

	v = v.Elem()

	// Build list of field pointers to scan into
	var fields []interface{}
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanSet() {
			fields = append(fields, v.Field(i).Addr().Interface())
		}
	}

	return row.Scan(fields...)
}

func scanSlice(rows *sql.Rows, dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dest must be a pointer to a slice")
	}

	sliceVal := v.Elem()
	elemType := sliceVal.Type().Elem()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		// Create new element
		elem := reflect.New(elemType).Elem()

		// Build scan destinations
		scanDest := make([]interface{}, len(columns))
		for i := 0; i < elem.NumField(); i++ {
			if i < len(scanDest) && elem.Field(i).CanSet() {
				scanDest[i] = elem.Field(i).Addr().Interface()
			}
		}

		// Fill remaining with dummy variables
		for i := elem.NumField(); i < len(scanDest); i++ {
			var dummy interface{}
			scanDest[i] = &dummy
		}

		if err := rows.Scan(scanDest...); err != nil {
			return err
		}

		sliceVal = reflect.Append(sliceVal, elem)
	}

	v.Elem().Set(sliceVal)
	return rows.Err()
}

func setTimestamps(data interface{}, isCreate bool) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	now := time.Now()

	if isCreate {
		if createdAt := v.FieldByName("CreatedAt"); createdAt.IsValid() && createdAt.CanSet() {
			createdAt.Set(reflect.ValueOf(now))
		}
	}

	if updatedAt := v.FieldByName("UpdatedAt"); updatedAt.IsValid() && updatedAt.CanSet() {
		updatedAt.Set(reflect.ValueOf(now))
	}
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// Timestamps adds created_at and updated_at fields
type Timestamps struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
