package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateScaffold generates a complete CRUD resource (model, migration, controller, views, routes)
func GenerateScaffold(name string, fields []string) error {
	fmt.Printf("Generating scaffold for %s...\n", name)

	// 1. Generate model
	fmt.Print("  Creating model...")
	if err := GenerateModel(name, fields); err != nil {
		return fmt.Errorf("failed to generate model: %v", err)
	}
	fmt.Println(" ✓")

	// Note: Migration is already created by GenerateModel

	// 2. Generate controller with all RESTful actions
	fmt.Print("  Creating controller...")
	if err := generateScaffoldController(name, fields); err != nil {
		return fmt.Errorf("failed to generate controller: %v", err)
	}
	fmt.Println(" ✓")

	// 3. Generate all views
	fmt.Print("  Creating views...")
	if err := generateScaffoldViews(name, fields); err != nil {
		return fmt.Errorf("failed to generate views: %v", err)
	}
	fmt.Println(" ✓")

	// 4. Update routes
	fmt.Print("  Updating routes...")
	if err := addScaffoldRoute(name); err != nil {
		return fmt.Errorf("failed to update routes: %v", err)
	}
	fmt.Println(" ✓")

	pluralName := ToPlural(name)
	fmt.Printf("\n✨ Scaffold generated successfully!\n\n")
	fmt.Printf("Next steps:\n")
	fmt.Printf("  1. Run migrations: gort db:migrate\n")
	fmt.Printf("  2. Start server: go run main.go\n")
	fmt.Printf("  3. Visit: http://localhost:3000/%s\n\n", strings.ToLower(pluralName))

	return nil
}

func generateScaffoldController(name string, fields []string) error {
	controllerName := ToPascalCase(name) + "Controller"
	pluralName := ToPlural(name)
	singularName := strings.ToLower(name)
	tableName := strings.ToLower(pluralName)

	// Get module name from go.mod
	moduleName, err := getModuleName()
	if err != nil {
		return err
	}

	// Parse fields to extract names and types
	fieldList := parseFieldsForScaffold(fields)

	content := fmt.Sprintf(`package controllers

import (
	"net/http"

	"github.com/AndrewNgKF/gort/pkg/gort"
	"%s/app/models"
)

type %s struct {
	gort.Controller
}

// Index - GET /%s
func (c *%s) Index(w http.ResponseWriter, r *http.Request) {
	var %s []models.%s
	
	if err := gort.FindAll("%s", &%s); err != nil {
		c.InternalServerError(w, err.Error())
		return
	}

	c.Render(w, r, "%s/index", map[string]interface{}{
		"%s": %s,
	})
}

// Show - GET /%s/:id
func (c *%s) Show(w http.ResponseWriter, r *http.Request) {
	id := c.Param(r, "id")
	var %s models.%s

	if err := gort.Find("%s", id, &%s); err != nil {
		c.NotFound(w)
		return
	}

	c.Render(w, r, "%s/show", map[string]interface{}{
		"%s": %s,
	})
}

// New - GET /%s/new
func (c *%s) New(w http.ResponseWriter, r *http.Request) {
	c.Render(w, r, "%s/new", map[string]interface{}{
		"%s": models.%s{},
	})
}

// Create - POST /%s
func (c *%s) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.BadRequest(w, "Invalid form data")
		return
	}

	%s := models.%s{
%s	}

	if err := gort.Create("%s", &%s); err != nil {
		c.InternalServerError(w, err.Error())
		return
	}

	c.Redirect(w, r, "/%s")
}

// Edit - GET /%s/:id/edit
func (c *%s) Edit(w http.ResponseWriter, r *http.Request) {
	id := c.Param(r, "id")
	var %s models.%s

	if err := gort.Find("%s", id, &%s); err != nil {
		c.NotFound(w)
		return
	}

	c.Render(w, r, "%s/edit", map[string]interface{}{
		"%s": %s,
	})
}

// Update - PUT/PATCH /%s/:id
func (c *%s) Update(w http.ResponseWriter, r *http.Request) {
	id := c.Param(r, "id")
	
	if err := r.ParseForm(); err != nil {
		c.BadRequest(w, "Invalid form data")
		return
	}

	var %s models.%s
	if err := gort.Find("%s", id, &%s); err != nil {
		c.NotFound(w)
		return
	}

%s
	if err := gort.Update("%s", id, &%s); err != nil {
		c.InternalServerError(w, err.Error())
		return
	}

	c.Redirect(w, r, "/%s")
}

// Destroy - DELETE /%s/:id
func (c *%s) Destroy(w http.ResponseWriter, r *http.Request) {
	id := c.Param(r, "id")

	if err := gort.Delete("%s", id); err != nil {
		c.InternalServerError(w, err.Error())
		return
	}

	c.Redirect(w, r, "/%s")
}
`,
		moduleName,
		controllerName,
		strings.ToLower(pluralName), controllerName,
		strings.ToLower(pluralName), ToPascalCase(name),
		tableName, strings.ToLower(pluralName),
		strings.ToLower(pluralName),
		strings.ToLower(pluralName), strings.ToLower(pluralName),
		strings.ToLower(pluralName), controllerName,
		singularName, ToPascalCase(name),
		tableName, singularName,
		strings.ToLower(pluralName),
		singularName, singularName,
		strings.ToLower(pluralName), controllerName,
		strings.ToLower(pluralName),
		singularName, ToPascalCase(name),
		strings.ToLower(pluralName), controllerName,
		singularName, ToPascalCase(name),
		generateCreateFields(fieldList),
		tableName, singularName,
		strings.ToLower(pluralName),
		strings.ToLower(pluralName), controllerName,
		singularName, ToPascalCase(name),
		tableName, singularName,
		strings.ToLower(pluralName),
		singularName, singularName,
		strings.ToLower(pluralName), controllerName,
		singularName, ToPascalCase(name),
		tableName, singularName,
		generateUpdateFields(fieldList, singularName),
		tableName, singularName,
		strings.ToLower(pluralName),
		strings.ToLower(pluralName), controllerName,
		tableName,
		strings.ToLower(pluralName),
	)

	controllerPath := filepath.Join("app", "controllers", strings.ToLower(pluralName)+"_controller.go")
	return os.WriteFile(controllerPath, []byte(content), 0644)
}

func generateScaffoldViews(name string, fields []string) error {
	pluralName := ToPlural(name)
	singularName := strings.ToLower(name)
	viewsDir := filepath.Join("app", "views", strings.ToLower(pluralName))

	if err := os.MkdirAll(viewsDir, 0755); err != nil {
		return err
	}

	fieldList := parseFieldsForScaffold(fields)

	// Index view
	indexContent := generateIndexView(singularName, pluralName, fieldList)
	if err := os.WriteFile(filepath.Join(viewsDir, "index.html"), []byte(indexContent), 0644); err != nil {
		return err
	}

	// Show view
	showContent := generateShowView(singularName, pluralName, fieldList)
	if err := os.WriteFile(filepath.Join(viewsDir, "show.html"), []byte(showContent), 0644); err != nil {
		return err
	}

	// New view
	newContent := generateNewView(singularName, pluralName, fieldList)
	if err := os.WriteFile(filepath.Join(viewsDir, "new.html"), []byte(newContent), 0644); err != nil {
		return err
	}

	// Edit view
	editContent := generateEditView(singularName, pluralName, fieldList)
	if err := os.WriteFile(filepath.Join(viewsDir, "edit.html"), []byte(editContent), 0644); err != nil {
		return err
	}

	// Form partial
	formContent := generateFormPartial(singularName, fieldList)
	if err := os.WriteFile(filepath.Join(viewsDir, "_form.html"), []byte(formContent), 0644); err != nil {
		return err
	}

	return nil
}

func generateIndexView(singular, plural string, fields []scaffoldField) string {
	headers := ""
	rows := ""

	for _, field := range fields {
		if field.Name != "id" {
			headers += fmt.Sprintf("          <th class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">%s</th>\n", ToPascalCase(field.Name))
			rows += fmt.Sprintf("          <td class=\"px-6 py-4 whitespace-nowrap text-sm text-gray-900\">{{ .%s }}</td>\n", ToPascalCase(field.Name))
		}
	}

	return fmt.Sprintf(`{{ define "content" }}
<div class="max-w-7xl mx-auto">
  <div class="flex justify-between items-center mb-6">
    <h1 class="text-3xl font-bold text-gray-900">%s</h1>
    <a href="/%s/new" class="bg-purple-600 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded transition duration-150">New %s</a>
  </div>
  
  <div class="bg-white shadow-md rounded-lg overflow-hidden">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
        <tr>
%s          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        {{ range .%s }}
        <tr class="hover:bg-gray-50">
%s          <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
            <a href="/%s/{{ .ID }}" class="text-indigo-600 hover:text-indigo-900 mr-3">Show</a>
            <a href="/%s/{{ .ID }}/edit" class="text-green-600 hover:text-green-900 mr-3">Edit</a>
            <form method="POST" action="/%s/{{ .ID }}" style="display:inline">
              <input type="hidden" name="_method" value="DELETE">
              <button type="submit" onclick="return confirm('Are you sure?')" class="text-red-600 hover:text-red-900">Delete</button>
            </form>
          </td>
        </tr>
        {{ end }}
      </tbody>
    </table>
  </div>
</div>
{{ end }}
`, ToPascalCase(plural), strings.ToLower(plural), ToPascalCase(singular), headers, strings.ToLower(plural), rows, strings.ToLower(plural), strings.ToLower(plural), strings.ToLower(plural))
}

func generateShowView(singular, plural string, fields []scaffoldField) string {
	fieldDisplay := ""

	for _, field := range fields {
		fieldDisplay += fmt.Sprintf("    <div class=\"border-b border-gray-200 py-4\">\n      <dt class=\"text-sm font-medium text-gray-500 mb-1\">%s</dt>\n      <dd class=\"text-sm text-gray-900\">{{ .%s.%s }}</dd>\n    </div>\n", ToPascalCase(field.Name), strings.ToLower(singular), ToPascalCase(field.Name))
	}

	return fmt.Sprintf(`{{ define "content" }}
<div class="max-w-3xl mx-auto">
  <div class="bg-white shadow-md rounded-lg overflow-hidden">
    <div class="px-6 py-4 bg-gray-50 border-b border-gray-200">
      <h1 class="text-2xl font-bold text-gray-900">%s Details</h1>
    </div>
    <div class="px-6 py-4">
      <dl>
%s      </dl>
    </div>
    <div class="px-6 py-4 bg-gray-50 flex gap-3">
      <a href="/%s" class="bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded transition duration-150">Back to List</a>
      <a href="/%s/{{ .%s.ID }}/edit" class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition duration-150">Edit</a>
    </div>
  </div>
</div>
{{ end }}
`, ToPascalCase(singular), fieldDisplay, strings.ToLower(plural), strings.ToLower(plural), strings.ToLower(singular))
}

func generateNewView(singular, plural string, fields []scaffoldField) string {
	return fmt.Sprintf(`{{ define "content" }}
<div class="max-w-2xl mx-auto">
  <div class="bg-white shadow-md rounded-lg overflow-hidden">
    <div class="px-6 py-4 bg-gray-50 border-b border-gray-200">
      <h1 class="text-2xl font-bold text-gray-900">New %s</h1>
    </div>
    <div class="px-6 py-4">
      <form method="POST" action="/%s">
        {{ template "_form" .%s }}
        <div class="flex gap-3 mt-6">
          <button type="submit" class="bg-purple-600 hover:bg-purple-700 text-white font-bold py-2 px-6 rounded transition duration-150">Create %s</button>
          <a href="/%s" class="bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-6 rounded transition duration-150">Back</a>
        </div>
      </form>
    </div>
  </div>
</div>
{{ end }}
`, ToPascalCase(singular), strings.ToLower(plural), strings.ToLower(singular), ToPascalCase(singular), strings.ToLower(plural))
}

func generateEditView(singular, plural string, fields []scaffoldField) string {
	return fmt.Sprintf(`{{ define "content" }}
<div class="max-w-2xl mx-auto">
  <div class="bg-white shadow-md rounded-lg overflow-hidden">
    <div class="px-6 py-4 bg-gray-50 border-b border-gray-200">
      <h1 class="text-2xl font-bold text-gray-900">Edit %s</h1>
    </div>
    <div class="px-6 py-4">
      <form method="POST" action="/%s/{{ .%s.ID }}">
        <input type="hidden" name="_method" value="PUT">
        {{ template "_form" .%s }}
        <div class="flex gap-3 mt-6">
          <button type="submit" class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-6 rounded transition duration-150">Update %s</button>
          <a href="/%s" class="bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-6 rounded transition duration-150">Back</a>
        </div>
      </form>
    </div>
  </div>
</div>
{{ end }}
`, ToPascalCase(singular), strings.ToLower(plural), strings.ToLower(singular), strings.ToLower(singular), ToPascalCase(singular), strings.ToLower(plural))
}

func generateFormPartial(singular string, fields []scaffoldField) string {
	formFields := "{{ define \"_form\" }}\n"

	for _, field := range fields {
		if field.Name == "id" || field.Name == "created_at" || field.Name == "updated_at" {
			continue
		}

		inputType := "text"
		if field.Type == "text" {
			formFields += fmt.Sprintf(`  <div class="mb-4">
    <label for="%s" class="block text-sm font-medium text-gray-700 mb-2">%s</label>
    <textarea name="%s" id="%s" rows="4" class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500">{{ .%s }}</textarea>
  </div>
`, field.Name, ToPascalCase(field.Name), field.Name, field.Name, ToPascalCase(field.Name))
			continue
		} else if field.Type == "boolean" {
			formFields += fmt.Sprintf(`  <div class="mb-4">
    <label class="flex items-center">
      <input type="checkbox" name="%s" value="true" {{ if .%s }}checked{{ end }} class="rounded border-gray-300 text-purple-600 shadow-sm focus:border-purple-300 focus:ring focus:ring-purple-200 focus:ring-opacity-50">
      <span class="ml-2 text-sm text-gray-700">%s</span>
    </label>
  </div>
`, field.Name, ToPascalCase(field.Name), ToPascalCase(field.Name))
			continue
		} else if field.Type == "integer" {
			inputType = "number"
		}

		formFields += fmt.Sprintf(`  <div class="mb-4">
    <label for="%s" class="block text-sm font-medium text-gray-700 mb-2">%s</label>
    <input type="%s" name="%s" id="%s" value="{{ .%s }}" class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500">
  </div>
`, field.Name, ToPascalCase(field.Name), inputType, field.Name, field.Name, ToPascalCase(field.Name))
	}

	formFields += "{{ end }}\n"
	return formFields
}

func addScaffoldRoute(name string) error {
	pluralName := ToPlural(name)
	routesPath := filepath.Join("config", "routes.go")

	content, err := os.ReadFile(routesPath)
	if err != nil {
		return err
	}

	controllerPkg := fmt.Sprintf("&controllers.%sController{}", ToPascalCase(name))
	routeLine := fmt.Sprintf("\tr.Resources(\"%s\", %s)", strings.ToLower(pluralName), controllerPkg)

	// Check if route already exists
	if strings.Contains(string(content), routeLine) {
		return nil
	}

	// Insert route before the return statement
	lines := strings.Split(string(content), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "return ") {
			// Insert before the return
			lines = append(lines[:i], append([]string{"", routeLine}, lines[i:]...)...)
			break
		}
	}

	return os.WriteFile(routesPath, []byte(strings.Join(lines, "\n")), 0644)
}

type scaffoldField struct {
	Name string
	Type string
}

func parseFieldsForScaffold(fields []string) []scaffoldField {
	var result []scaffoldField

	for _, field := range fields {
		parts := strings.Split(field, ":")
		if len(parts) >= 2 {
			result = append(result, scaffoldField{
				Name: parts[0],
				Type: parts[1],
			})
		}
	}

	return result
}

func generateCreateFields(fields []scaffoldField) string {
	var lines []string

	for _, field := range fields {
		if field.Name == "id" || field.Name == "created_at" || field.Name == "updated_at" {
			continue
		}

		fieldName := ToPascalCase(field.Name)
		formField := field.Name

		if field.Type == "boolean" {
			lines = append(lines, fmt.Sprintf("\t\t%s: r.FormValue(\"%s\") == \"true\",", fieldName, formField))
		} else if field.Type == "integer" {
			lines = append(lines, fmt.Sprintf("\t\t%s: parseInt(r.FormValue(\"%s\")),", fieldName, formField))
		} else {
			lines = append(lines, fmt.Sprintf("\t\t%s: r.FormValue(\"%s\"),", fieldName, formField))
		}
	}

	return strings.Join(lines, "\n") + "\n"
}

func generateUpdateFields(fields []scaffoldField, varName string) string {
	if len(fields) == 0 {
		return ""
	}

	var lines []string

	for _, field := range fields {
		if field.Name == "id" || field.Name == "created_at" || field.Name == "updated_at" {
			continue
		}

		fieldName := ToPascalCase(field.Name)
		formField := field.Name

		if field.Type == "boolean" {
			lines = append(lines, fmt.Sprintf("\t%s.%s = r.FormValue(\"%s\") == \"true\"", varName, fieldName, formField))
		} else if field.Type == "integer" {
			lines = append(lines, fmt.Sprintf("\t%s.%s = parseInt(r.FormValue(\"%s\"))", varName, fieldName, formField))
		} else {
			lines = append(lines, fmt.Sprintf("\t%s.%s = r.FormValue(\"%s\")", varName, fieldName, formField))
		}
	}

	return strings.Join(lines, "\n") + "\n"
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

func getModuleName() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module "), nil
		}
	}

	return "", fmt.Errorf("module name not found in go.mod")
}
