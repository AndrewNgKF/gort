package generator

func createControllerFiles(name string) error {
	baseController := `package controllers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type BaseController struct{}

func (c *BaseController) Render(w http.ResponseWriter, templateName string, data interface{}) {
	layoutPath := filepath.Join("app", "views", "layouts", "application.html")
	templatePath := filepath.Join("app", "views", templateName+".html")
	
	tmpl, err := template.ParseFiles(layoutPath, templatePath)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}
`

	homeController := `package controllers

import (
	"net/http"
	"runtime"
)

func HomeIndex(w http.ResponseWriter, r *http.Request) {
	controller := &BaseController{}
	
	data := map[string]interface{}{
		"Title":       "Welcome to Gort!",
		"Message":     "Your Rails-inspired Go framework is ready!",
		"GoVersion":   runtime.Version(),
		"GortVersion": "0.1.0",
	}
	
	controller.Render(w, "home/index", data)
}
`

	if err := writeFile(name, "app/controllers/application_controller.go", baseController); err != nil {
		return err
	}
	return writeFile(name, "app/controllers/home_controller.go", homeController)
}
