package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/AndrewNgKF/gort/internal/router"
)

// Base provides common controller functionality
type Base struct{}

// Render renders an HTML template
func (c *Base) Render(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	layoutPath := filepath.Join("app", "views", "layouts", "application.html")
	templatePath := filepath.Join("app", "views", templateName+".html")

	// Parse layout and main template
	tmpl, err := template.ParseFiles(layoutPath, templatePath)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse partials (files starting with _) from the same directory as the template
	templateDir := filepath.Dir(templatePath)
	partials, _ := filepath.Glob(filepath.Join(templateDir, "_*.html"))
	if len(partials) > 0 {
		tmpl, err = tmpl.ParseFiles(partials...)
		if err != nil {
			http.Error(w, "Partial template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}

// JSON renders JSON response
func (c *Base) JSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// JSONStatus renders JSON with custom status code
func (c *Base) JSONStatus(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Text renders plain text
func (c *Base) Text(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(text))
}

// Redirect redirects to a URL
func (c *Base) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusFound)
}

// NotFound renders 404
func (c *Base) NotFound(w http.ResponseWriter) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

// BadRequest renders 400
func (c *Base) BadRequest(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

// Unauthorized renders 401
func (c *Base) Unauthorized(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

// Forbidden renders 403
func (c *Base) Forbidden(w http.ResponseWriter) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}

// UnprocessableEntity renders 422
func (c *Base) UnprocessableEntity(w http.ResponseWriter, data interface{}) {
	c.JSONStatus(w, http.StatusUnprocessableEntity, data)
}

// InternalServerError renders 500
func (c *Base) InternalServerError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}

// Params retrieves route parameters
func (c *Base) Params(r *http.Request) map[string]string {
	return router.GetParams(r.Context())
}

// Param retrieves a single route parameter
func (c *Base) Param(r *http.Request, name string) string {
	return router.GetParam(r.Context(), name)
}

// BindJSON binds JSON request body to a struct
func (c *Base) BindJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
