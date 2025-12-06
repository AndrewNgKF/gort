package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GenerateModel creates a new model file with specified fields
func GenerateModel(name string, fields []string) error {
	// Check if we're in a Gort app directory
	if _, err := os.Stat("app/models"); os.IsNotExist(err) {
		return fmt.Errorf("not in a Gort application directory (app/models not found)")
	}

	// Convert name to proper format
	modelName := ToPascalCase(name)
	tableName := ToPlural(ToSnakeCase(name))

	fileName := ToSnakeCase(name) + ".go"
	filePath := filepath.Join("app", "models", fileName)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("model already exists: %s", filePath)
	}

	// Parse fields
	modelFields := parseFields(fields)

	// Generate model content
	content := generateModelContent(modelName, tableName, modelFields)

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	fmt.Printf("  Created: %s\n", filePath)

	// Generate migration
	migrationName := fmt.Sprintf("create_%s", tableName)
	if err := GenerateMigration(migrationName, modelFields, tableName); err != nil {
		fmt.Printf("Warning: failed to create migration: %v\n", err)
	}

	return nil
}

type Field struct {
	Name       string
	Type       string
	GoType     string
	DBType     string
	Nullable   bool
	Unique     bool
	Index      bool
	References string
}

func parseFields(fields []string) []Field {
	var result []Field

	for _, fieldStr := range fields {
		parts := strings.Split(fieldStr, ":")
		if len(parts) < 2 {
			continue
		}

		fieldName := parts[0]
		fieldType := parts[1]

		field := Field{
			Name:   ToPascalCase(fieldName),
			Type:   fieldType,
			GoType: sqlTypeToGoType(fieldType),
			DBType: goTypeToSQLType(fieldType),
		}

		// Parse modifiers
		if len(parts) > 2 {
			for _, modifier := range parts[2:] {
				switch modifier {
				case "unique":
					field.Unique = true
				case "index":
					field.Index = true
				case "nullable", "null":
					field.Nullable = true
				default:
					if strings.HasPrefix(modifier, "ref:") {
						field.References = strings.TrimPrefix(modifier, "ref:")
					}
				}
			}
		}

		result = append(result, field)
	}

	return result
}

func generateModelContent(modelName, tableName string, fields []Field) string {
	var sb strings.Builder

	sb.WriteString("package models\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"time\"\n")
	sb.WriteString(")\n\n")

	// Struct definition
	sb.WriteString(fmt.Sprintf("type %s struct {\n", modelName))
	sb.WriteString("\tID        int64     `db:\"id\"`\n")

	for _, field := range fields {
		tag := fmt.Sprintf("`db:\"%s\"`", ToSnakeCase(field.Name))
		padding := 10 - len(field.Name)
		if padding < 1 {
			padding = 1
		}
		sb.WriteString(fmt.Sprintf("\t%s%s%s %s\n", field.Name, strings.Repeat(" ", padding), field.GoType, tag))
	}

	sb.WriteString("\tCreatedAt time.Time `db:\"created_at\"`\n")
	sb.WriteString("\tUpdatedAt time.Time `db:\"updated_at\"`\n")
	sb.WriteString("}\n\n")

	// Table name method
	sb.WriteString(fmt.Sprintf("// TableName returns the table name for this model\n"))
	sb.WriteString(fmt.Sprintf("func (%s) TableName() string {\n", modelName))
	sb.WriteString(fmt.Sprintf("\treturn \"%s\"\n", tableName))
	sb.WriteString("}\n\n")

	// Add comment about where to add validations and associations
	sb.WriteString("// Add validations, associations, and callbacks here\n")
	sb.WriteString("// Example:\n")
	sb.WriteString(fmt.Sprintf("// func (m *%s) Validate() error {\n", modelName))
	sb.WriteString("//     return gort.Validate(m,\n")
	sb.WriteString("//         gort.ValidatePresence(\"FieldName\"),\n")
	sb.WriteString("//         gort.ValidateLength(\"FieldName\", gort.Min(3)),\n")
	sb.WriteString("//     )\n")
	sb.WriteString("// }\n")
	sb.WriteString("//\n")
	sb.WriteString("// Associations:\n")
	sb.WriteString(fmt.Sprintf("// type %s struct {\n", modelName))
	sb.WriteString("//     ...\n")
	sb.WriteString("//     RelatedModel []RelatedModel `gort:\"has_many\"`\n")
	sb.WriteString("//     ParentModel  ParentModel   `gort:\"belongs_to\"`\n")
	sb.WriteString("// }\n")

	return sb.String()
}

func sqlTypeToGoType(sqlType string) string {
	switch strings.ToLower(sqlType) {
	case "string", "text", "varchar":
		return "string"
	case "integer", "int":
		return "int"
	case "bigint", "int64":
		return "int64"
	case "boolean", "bool":
		return "bool"
	case "float", "decimal":
		return "float64"
	case "datetime", "timestamp":
		return "time.Time"
	case "date":
		return "time.Time"
	case "uuid":
		return "string"
	case "json", "jsonb":
		return "interface{}"
	case "array":
		return "[]string"
	default:
		return "string"
	}
}

func goTypeToSQLType(goType string) string {
	switch strings.ToLower(goType) {
	case "string":
		return "VARCHAR(255)"
	case "text":
		return "TEXT"
	case "integer", "int":
		return "INTEGER"
	case "bigint", "int64":
		return "BIGINT"
	case "boolean", "bool":
		return "BOOLEAN"
	case "float", "decimal", "float64":
		return "DECIMAL(10,2)"
	case "datetime", "timestamp":
		return "TIMESTAMP"
	case "date":
		return "DATE"
	case "uuid":
		return "UUID"
	case "json":
		return "JSON"
	case "jsonb":
		return "JSONB"
	case "array":
		return "TEXT[]"
	default:
		return "VARCHAR(255)"
	}
}

// ToPlural converts a word to plural (simple version)
func ToPlural(s string) string {
	if strings.HasSuffix(s, "s") {
		return s + "es"
	}
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	return s + "s"
}

// GenerateMigration creates a migration file for the model
func GenerateMigration(name string, fields []Field, tableName string) error {
	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll("db/migrations", 0755); err != nil {
		return err
	}

	// Generate timestamp-based filename
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s.go", timestamp, name)
	filePath := filepath.Join("db", "migrations", fileName)

	content := generateMigrationContent(name, fields, tableName, timestamp)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}

	fmt.Printf("  Created: %s\n", filePath)
	return nil
}

func generateMigrationContent(name string, fields []Field, tableName, timestamp string) string {
	var sb strings.Builder

	sb.WriteString("package migrations\n\n")
	sb.WriteString("import \"database/sql\"\n\n")

	// Up function
	sb.WriteString(fmt.Sprintf("func Up_%s(db *sql.DB) error {\n", timestamp))
	sb.WriteString(fmt.Sprintf("\tquery := `\n"))
	sb.WriteString(fmt.Sprintf("\t\tCREATE TABLE %s (\n", tableName))
	sb.WriteString("\t\t\tid SERIAL PRIMARY KEY,\n")

	for _, field := range fields {
		nullable := ""
		if !field.Nullable {
			nullable = " NOT NULL"
		}
		unique := ""
		if field.Unique {
			unique = " UNIQUE"
		}

		sb.WriteString(fmt.Sprintf("\t\t\t%s %s%s%s,\n",
			ToSnakeCase(field.Name), field.DBType, nullable, unique))
	}

	sb.WriteString("\t\t\tcreated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,\n")
	sb.WriteString("\t\t\tupdated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP\n")
	sb.WriteString("\t\t)\n\t`\n")
	sb.WriteString("\t_, err := db.Exec(query)\n")
	sb.WriteString("\treturn err\n")
	sb.WriteString("}\n\n")

	// Down function
	sb.WriteString(fmt.Sprintf("func Down_%s(db *sql.DB) error {\n", timestamp))
	sb.WriteString(fmt.Sprintf("\tquery := `DROP TABLE IF EXISTS %s`\n", tableName))
	sb.WriteString("\t_, err := db.Exec(query)\n")
	sb.WriteString("\treturn err\n")
	sb.WriteString("}\n")

	return sb.String()
}
