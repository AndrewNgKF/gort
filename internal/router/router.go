package router

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Router is the main routing engine
type Router struct {
	routes     []*Route
	middleware []MiddlewareFunc
	mux        *http.ServeMux
}

// Route represents a single route
type Route struct {
	Method     string
	Path       string
	Pattern    *regexp.Regexp
	Handler    http.HandlerFunc
	RouteName  string
	ParamNames []string
}

// MiddlewareFunc is a function that wraps an http.Handler
type MiddlewareFunc func(http.Handler) http.Handler

// New creates a new Router
func New() *Router {
	return &Router{
		routes: make([]*Route, 0),
		mux:    http.NewServeMux(),
	}
}

// Get registers a GET route
func (r *Router) Get(path string, handler http.HandlerFunc) *Route {
	return r.addRoute("GET", path, handler)
}

// Post registers a POST route
func (r *Router) Post(path string, handler http.HandlerFunc) *Route {
	return r.addRoute("POST", path, handler)
}

// Put registers a PUT route
func (r *Router) Put(path string, handler http.HandlerFunc) *Route {
	return r.addRoute("PUT", path, handler)
}

// Patch registers a PATCH route
func (r *Router) Patch(path string, handler http.HandlerFunc) *Route {
	return r.addRoute("PATCH", path, handler)
}

// Delete registers a DELETE route
func (r *Router) Delete(path string, handler http.HandlerFunc) *Route {
	return r.addRoute("DELETE", path, handler)
}

// Use adds middleware to the router
func (r *Router) Use(middleware MiddlewareFunc) {
	r.middleware = append(r.middleware, middleware)
}

// addRoute adds a route to the router
func (r *Router) addRoute(method, path string, handler http.HandlerFunc) *Route {
	route := &Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}

	// Convert path to regex pattern and extract param names
	pattern, paramNames := pathToRegex(path)
	route.Pattern = pattern
	route.ParamNames = paramNames

	r.routes = append(r.routes, route)
	return route
}

// Name sets the name of the route
func (route *Route) Name(name string) *Route {
	route.RouteName = name
	return route
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Handle method override before routing
	if req.Method == "POST" {
		req.ParseForm()
		if method := req.PostFormValue("_method"); method != "" {
			req.Method = method
		}
	}

	// Apply middleware
	handler := r.match(req)

	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}

	handler.ServeHTTP(w, req)
}

// match finds the matching route and returns its handler
func (r *Router) match(req *http.Request) http.Handler {
	for _, route := range r.routes {
		if route.Method != req.Method {
			continue
		}

		matches := route.Pattern.FindStringSubmatch(req.URL.Path)
		if matches == nil {
			continue
		}

		// Extract params and add to context
		if len(route.ParamNames) > 0 {
			params := make(map[string]string)
			for i, name := range route.ParamNames {
				params[name] = matches[i+1]
			}

			// Add params to request context
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := WithParams(r.Context(), params)
				route.Handler(w, r.WithContext(ctx))
			})
		}

		return route.Handler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}

// GetRoutes returns all registered routes
func (r *Router) GetRoutes() []*Route {
	return r.routes
}

// pathToRegex converts a path pattern to a regex
// Converts /users/:id to ^/users/([^/]+)$ and extracts ["id"]
func pathToRegex(path string) (*regexp.Regexp, []string) {
	paramNames := []string{}

	// Find all :param patterns
	paramRegex := regexp.MustCompile(`:(\w+)`)
	matches := paramRegex.FindAllStringSubmatch(path, -1)

	for _, match := range matches {
		paramNames = append(paramNames, match[1])
	}

	// Replace :param with regex capture groups
	pattern := paramRegex.ReplaceAllString(path, `([^/]+)`)

	// Escape special regex characters except our capture groups
	pattern = strings.ReplaceAll(pattern, ".", `\.`)

	// Add anchors
	pattern = "^" + pattern + "$"

	return regexp.MustCompile(pattern), paramNames
}

// Root sets the root route (/)
func (r *Router) Root(handler http.HandlerFunc) *Route {
	return r.Get("/", handler)
}

// Resources creates RESTful routes for a resource
// Example: r.Resources("posts", postsController)
func (r *Router) Resources(name string, controller interface{}) {
	// RESTful routes:
	// GET    /posts          -> Index
	// GET    /posts/new      -> New
	// GET    /posts/:id      -> Show
	// GET    /posts/:id/edit -> Edit
	// POST   /posts          -> Create
	// PUT    /posts/:id      -> Update
	// DELETE /posts/:id      -> Destroy

	basePath := "/" + name

	// Index - GET /posts
	if handler := getControllerMethod(controller, "Index"); handler != nil {
		r.Get(basePath, handler).Name(name + "#index")
	}

	// New - GET /posts/new (must come before Show to avoid conflict)
	if handler := getControllerMethod(controller, "New"); handler != nil {
		r.Get(basePath+"/new", handler).Name(name + "#new")
	}

	// Show - GET /posts/:id
	if handler := getControllerMethod(controller, "Show"); handler != nil {
		r.Get(basePath+"/:id", handler).Name(name + "#show")
	}

	// Edit - GET /posts/:id/edit
	if handler := getControllerMethod(controller, "Edit"); handler != nil {
		r.Get(basePath+"/:id/edit", handler).Name(name + "#edit")
	}

	// Create - POST /posts
	if handler := getControllerMethod(controller, "Create"); handler != nil {
		r.Post(basePath, handler).Name(name + "#create")
	}

	// Update - PUT /posts/:id
	if handler := getControllerMethod(controller, "Update"); handler != nil {
		r.Put(basePath+"/:id", handler).Name(name + "#update")
	}

	// Destroy - DELETE /posts/:id
	if handler := getControllerMethod(controller, "Destroy"); handler != nil {
		r.Delete(basePath+"/:id", handler).Name(name + "#destroy")
	}
}

// getControllerMethod uses reflection to get controller methods
func getControllerMethod(controller interface{}, methodName string) http.HandlerFunc {
	// This is a simplified version - in production we'd use reflection
	// For now, we'll require controllers to implement specific interfaces

	// Type assertion for different controller methods
	type indexer interface {
		Index(http.ResponseWriter, *http.Request)
	}
	type newer interface {
		New(http.ResponseWriter, *http.Request)
	}
	type shower interface {
		Show(http.ResponseWriter, *http.Request)
	}
	type editor interface {
		Edit(http.ResponseWriter, *http.Request)
	}
	type creator interface {
		Create(http.ResponseWriter, *http.Request)
	}
	type updater interface {
		Update(http.ResponseWriter, *http.Request)
	}
	type destroyer interface {
		Destroy(http.ResponseWriter, *http.Request)
	}

	switch methodName {
	case "Index":
		if c, ok := controller.(indexer); ok {
			return c.Index
		}
	case "New":
		if c, ok := controller.(newer); ok {
			return c.New
		}
	case "Show":
		if c, ok := controller.(shower); ok {
			return c.Show
		}
	case "Edit":
		if c, ok := controller.(editor); ok {
			return c.Edit
		}
	case "Create":
		if c, ok := controller.(creator); ok {
			return c.Create
		}
	case "Update":
		if c, ok := controller.(updater); ok {
			return c.Update
		}
	case "Destroy":
		if c, ok := controller.(destroyer); ok {
			return c.Destroy
		}
	}

	return nil
}

// Namespace creates a route group with a prefix
func (r *Router) Namespace(prefix string, fn func(*Router)) {
	// Create a sub-router with the prefix
	subRouter := &Router{
		routes:     make([]*Route, 0),
		middleware: r.middleware,
	}

	// Execute the callback
	fn(subRouter)

	// Add all sub-routes with prefix
	for _, route := range subRouter.routes {
		route.Path = prefix + route.Path
		pattern, paramNames := pathToRegex(route.Path)
		route.Pattern = pattern
		route.ParamNames = paramNames
		r.routes = append(r.routes, route)
	}
}

// Group creates a route group (alias for Namespace)
func (r *Router) Group(prefix string, fn func(*Router)) {
	r.Namespace(prefix, fn)
}

// PrintRoutes prints all routes (useful for debugging)
func (r *Router) PrintRoutes() {
	fmt.Println("\nRegistered Routes:")
	fmt.Println("==================")
	for _, route := range r.routes {
		name := route.RouteName
		if name == "" {
			name = "-"
		}
		fmt.Printf("%-6s %-30s %s\n", route.Method, route.Path, name)
	}
	fmt.Println()
}
