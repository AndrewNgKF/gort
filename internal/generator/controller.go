package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// GenerateController creates a new controller file with specified actions
func GenerateController(name string, actions []string) error {
	// Check if we're in a Gort app directory
	if _, err := os.Stat("app/controllers"); os.IsNotExist(err) {
		return fmt.Errorf("not in a Gort application directory (app/controllers not found)")
	}

	// Convert name to proper format
	controllerName := ToPascalCase(name)
	if !strings.HasSuffix(controllerName, "Controller") {
		controllerName += "Controller"
	}

	fileName := ToSnakeCase(name) + "_controller.go"
	filePath := filepath.Join("app", "controllers", fileName)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("controller already exists: %s", filePath)
	}

	// Generate controller content
	content := generateControllerContent(controllerName, name, actions)

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create controller: %w", err)
	}

	// Create view directories and files for each action
	viewDir := filepath.Join("app", "views", ToSnakeCase(name))
	if err := os.MkdirAll(viewDir, 0755); err != nil {
		return fmt.Errorf("failed to create view directory: %w", err)
	}

	for _, action := range actions {
		if isViewAction(action) {
			viewFile := filepath.Join(viewDir, action+".html")
			viewContent := generateViewContent(name, action)
			if err := os.WriteFile(viewFile, []byte(viewContent), 0644); err != nil {
				fmt.Printf("Warning: failed to create view %s: %v\n", viewFile, err)
			}
		}
	}

	fmt.Printf("  Created: %s\n", filePath)
	if len(actions) > 0 {
		fmt.Printf("  Created: %s/ (view directory)\n", viewDir)
	}

	return nil
}

func generateControllerContent(controllerName, baseName string, actions []string) string {
	var sb strings.Builder

	sb.WriteString("package controllers\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"net/http\"\n")
	sb.WriteString("\t\"github.com/AndrewNgKF/gort/pkg/gort\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString(fmt.Sprintf("type %s struct {\n", controllerName))
	sb.WriteString("\tgort.Controller\n")
	sb.WriteString("}\n\n")

	// Generate methods for each action
	if len(actions) == 0 {
		// If no actions specified, create an empty controller
		sb.WriteString(fmt.Sprintf("// Add your %s methods here\n", controllerName))
	} else {
		for _, action := range actions {
			sb.WriteString(generateActionMethod(controllerName, baseName, action))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func generateActionMethod(controllerName, baseName, action string) string {
	var sb strings.Builder

	actionName := ToPascalCase(action)
	httpMethod := getHTTPMethod(action)
	path := getActionPath(baseName, action)

	sb.WriteString(fmt.Sprintf("// %s handles %s %s\n", actionName, httpMethod, path))
	sb.WriteString(fmt.Sprintf("func (c *%s) %s(w http.ResponseWriter, r *http.Request) {\n", controllerName, actionName))

	switch action {
	case "index":
		sb.WriteString("\t// List all resources\n")
		sb.WriteString("\tdata := gort.H{\n")
		sb.WriteString(fmt.Sprintf("\t\t\"Title\": \"%s\",\n", baseName))
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tc.Render(w, r, \"%s/index\", data)\n", ToSnakeCase(baseName)))

	case "show":
		sb.WriteString("\t// Show a single resource\n")
		sb.WriteString("\tid := c.Param(r, \"id\")\n")
		sb.WriteString("\tdata := gort.H{\n")
		sb.WriteString(fmt.Sprintf("\t\t\"Title\": \"%s Details\",\n", baseName))
		sb.WriteString("\t\t\"ID\":    id,\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tc.Render(w, r, \"%s/show\", data)\n", ToSnakeCase(baseName)))

	case "new":
		sb.WriteString("\t// Show form for creating a new resource\n")
		sb.WriteString("\tdata := gort.H{\n")
		sb.WriteString(fmt.Sprintf("\t\t\"Title\": \"New %s\",\n", baseName))
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tc.Render(w, r, \"%s/new\", data)\n", ToSnakeCase(baseName)))

	case "edit":
		sb.WriteString("\t// Show form for editing a resource\n")
		sb.WriteString("\tid := c.Param(r, \"id\")\n")
		sb.WriteString("\tdata := gort.H{\n")
		sb.WriteString(fmt.Sprintf("\t\t\"Title\": \"Edit %s\",\n", baseName))
		sb.WriteString("\t\t\"ID\":    id,\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tc.Render(w, r, \"%s/edit\", data)\n", ToSnakeCase(baseName)))

	case "create":
		sb.WriteString("\t// Create a new resource\n")
		sb.WriteString("\t// TODO: Parse form data and save to database\n")
		sb.WriteString("\t\n")
		sb.WriteString("\tc.JSON(w, gort.H{\"message\": \"Created successfully\"})\n")
		sb.WriteString(fmt.Sprintf("\t// c.Redirect(w, r, \"/%s\")\n", ToSnakeCase(baseName)))

	case "update":
		sb.WriteString("\t// Update a resource\n")
		sb.WriteString("\tid := c.Param(r, \"id\")\n")
		sb.WriteString("\t// TODO: Parse form data and update in database\n")
		sb.WriteString("\t\n")
		sb.WriteString("\tc.JSON(w, gort.H{\"message\": \"Updated successfully\", \"id\": id})\n")
		sb.WriteString(fmt.Sprintf("\t// c.Redirect(w, r, \"/%s/\" + id)\n", ToSnakeCase(baseName)))

	case "destroy", "delete":
		sb.WriteString("\t// Delete a resource\n")
		sb.WriteString("\tid := c.Param(r, \"id\")\n")
		sb.WriteString("\t// TODO: Delete from database\n")
		sb.WriteString("\t\n")
		sb.WriteString("\tc.JSON(w, gort.H{\"message\": \"Deleted successfully\", \"id\": id})\n")
		sb.WriteString(fmt.Sprintf("\t// c.Redirect(w, r, \"/%s\")\n", ToSnakeCase(baseName)))

	default:
		sb.WriteString(fmt.Sprintf("\t// TODO: Implement %s action\n", action))
		sb.WriteString("\tdata := gort.H{\n")
		sb.WriteString(fmt.Sprintf("\t\t\"Title\": \"%s %s\",\n", baseName, actionName))
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tc.Render(w, r, \"%s/%s\", data)\n", ToSnakeCase(baseName), action))
	}

	sb.WriteString("}\n")

	return sb.String()
}

func generateViewContent(baseName, action string) string {
	title := ToPascalCase(baseName) + " " + ToPascalCase(action)

	content := fmt.Sprintf(`{{ define "content" }}
<div class="container">
    <h1>%s</h1>
    
    <p>This is the %s view for %s.</p>
    
    <!-- Add your content here -->
</div>
{{ end }}
`, title, action, baseName)

	return content
}

func isViewAction(action string) bool {
	viewActions := []string{"index", "show", "new", "edit"}
	for _, va := range viewActions {
		if action == va {
			return true
		}
	}
	return false
}

func getHTTPMethod(action string) string {
	switch action {
	case "create":
		return "POST"
	case "update":
		return "PUT"
	case "destroy", "delete":
		return "DELETE"
	default:
		return "GET"
	}
}

func getActionPath(baseName, action string) string {
	base := "/" + ToSnakeCase(baseName)

	switch action {
	case "index", "create":
		return base
	case "new":
		return base + "/new"
	case "show", "edit", "update", "destroy", "delete":
		return base + "/:id"
	default:
		return base + "/" + action
	}
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// Split by underscore or space
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == ' ' || r == '-'
	})

	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, "")
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	var result strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	return strings.ReplaceAll(result.String(), "-", "_")
}
