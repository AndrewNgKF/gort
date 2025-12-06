package generator

import (
	"fmt"
	"strings"
)

func createMainGo(name string) error {
	moduleName := strings.ReplaceAll(name, "-", "_")
	content := fmt.Sprintf(`package main

import (
	"fmt"
	"log"
	"net/http"

	"%s/config"
)

func main() {
	// Initialize configuration
	config.Init()
	
	// Initialize database (non-fatal, just warn if it fails)
	if err := config.InitDatabase(); err != nil {
		log.Printf("⚠️  Database connection failed: %%v\n", err)
		log.Println("   Application will run without database support")
		log.Println("   Configure your database in config/database.yml and restart")
	}

	// Set up routes
	router := config.Routes()

	// Print routes in development
	if config.Get("environment") == "development" {
		router.PrintRoutes()
	}

	// Start server
	port := config.Get("port", "3000")
	fmt.Printf("🚀 Starting %s on http://localhost:%%s\n", port)
	fmt.Println("   Press Ctrl+C to stop")
	
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
`, moduleName, name)

	return writeFile(name, "main.go", content)
}

func createGoMod(name string) error {
	moduleName := strings.ReplaceAll(name, "-", "_")
	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/AndrewNgKF/gort latest
)
`, moduleName)

	return writeFile(name, "go.mod", content)
}

func createReadme(name string) error {
	title := strings.Title(strings.ReplaceAll(name, "_", " "))
	content := fmt.Sprintf(`# %s

A Gort application.

## Getting Started

### Setup

bash
# Install dependencies
go mod tidy

# Create database
gort db:create

# Run migrations
gort db:migrate
`, title)

	return writeFile(name, "README.md", content)
}

func createGitignore(name string) error {
	content := `# Binaries
*.exe
*.dll
*.so
*.dylib
bin/

# Test binary
*.test

# Coverage
*.out

# Database
*.db
*.sqlite

# Logs
*.log
tmp/

# Environment
.env

# OS files
.DS_Store

# IDE
.vscode/
.idea/
`

	return writeFile(name, ".gitignore", content)
}
