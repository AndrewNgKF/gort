package generator

import (
	"fmt"
	"strings"
)

func createConfigFiles(name string) error {
	if err := createApplicationConfig(name); err != nil {
		return err
	}
	if err := createDatabaseConfig(name); err != nil {
		return err
	}
	if err := createRoutesConfig(name); err != nil {
		return err
	}
	if err := createEnvConfigs(name); err != nil {
		return err
	}
	return nil
}

func createApplicationConfig(name string) error {
	content := fmt.Sprintf(`package config

import (
	"fmt"
	"os"
	"strconv"
	
	"github.com/AndrewNgKF/gort/pkg/gort"
	"gopkg.in/yaml.v3"
)

var config map[string]string

// Init initializes the application configuration
func Init() {
	config = make(map[string]string)
	
	config["port"] = getEnv("PORT", "3000")
	config["environment"] = getEnv("GO_ENV", "development")
	config["app_name"] = "%s"
	
	loadEnvironmentConfig()
}

// Get retrieves a configuration value
func Get(key string, defaultValue ...string) string {
	if val, ok := config[key]; ok {
		return val
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// Set sets a configuration value
func Set(key, value string) {
	config[key] = value
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadEnvironmentConfig() {
	env := config["environment"]
	
	switch env {
	case "development":
		loadDevelopmentConfig()
	case "production":
		loadProductionConfig()
	case "test":
		loadTestConfig()
	}
}

// InitDatabase initializes the database connection
func InitDatabase() error {
	data, err := os.ReadFile("config/database.yml")
	if err != nil {
		return fmt.Errorf("failed to read database config: %%w", err)
	}

	var dbConfig map[string]map[string]interface{}
	if err := yaml.Unmarshal(data, &dbConfig); err != nil {
		return fmt.Errorf("failed to parse database config: %%w", err)
	}

	env := config["environment"]
	envConfig, ok := dbConfig[env]
	if !ok {
		return fmt.Errorf("database config not found for environment: %%s", env)
	}

	dbConf := gort.DBConfig{
		Adapter:  getString(envConfig, "adapter"),
		Database: getString(envConfig, "database"),
		Host:     getString(envConfig, "host"),
		Port:     getStringFromInt(envConfig, "port"),
		Username: getString(envConfig, "username"),
		Password: getString(envConfig, "password"),
		SSLMode:  getString(envConfig, "sslmode"),
	}

	return gort.DB.Connect(dbConf)
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getStringFromInt(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		case int:
			return strconv.Itoa(v)
		case int64:
			return strconv.FormatInt(v, 10)
		}
	}
	return ""
}
`, name)

	return writeFile(name, "config/application.go", content)
}

func createDatabaseConfig(name string) error {
	dbName := strings.ReplaceAll(name, "-", "_")
	content := fmt.Sprintf(`development:
  adapter: postgresql
  database: %s_development
  host: localhost
  port: 5432
  username: postgres
  password: postgres

test:
  adapter: postgresql
  database: %s_test
  host: localhost
  port: 5432
  username: postgres
  password: postgres

production:
  adapter: postgresql
  database: %s_production
  host: ${DB_HOST}
  port: 5432
  username: ${DB_USER}
  password: ${DB_PASSWORD}
`, dbName, dbName, dbName)

	return writeFile(name, "config/database.yml", content)
}

func createRoutesConfig(name string) error {
	moduleName := strings.ReplaceAll(name, "-", "_")
	content := fmt.Sprintf(`package config

import (
	"github.com/AndrewNgKF/gort/pkg/gort"
	"%s/app/controllers"
)

// Routes sets up the application routes
func Routes() *gort.Router {
	r := gort.NewRouter()
	
	// Middleware
	r.Use(gort.Logger())
	r.Use(gort.Recovery())
	r.Use(gort.MethodOverride())
	
	// Root route
	r.Root(controllers.HomeIndex)
	
	// Example routes:
	// r.Get("/about", controllers.AboutPage)
	// r.Resources("posts", &controllers.PostsController{})
	
	return r
}
`, moduleName)

	return writeFile(name, "config/routes.go", content)
}

func createEnvConfigs(name string) error {
	dev := `package config

func loadDevelopmentConfig() {
	Set("log_level", "debug")
	Set("cache_enabled", "false")
}
`
	prod := `package config

func loadProductionConfig() {
	Set("log_level", "info")
	Set("cache_enabled", "true")
}
`
	test := `package config

func loadTestConfig() {
	Set("log_level", "error")
}
`

	if err := writeFile(name, "config/development.go", dev); err != nil {
		return err
	}
	if err := writeFile(name, "config/production.go", prod); err != nil {
		return err
	}
	return writeFile(name, "config/test.go", test)
}
