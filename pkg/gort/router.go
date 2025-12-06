package gort

import (
	"context"
	"net/http"

	"github.com/AndrewNgKF/gort/internal/router"
)

// Router is the main routing engine
type Router = router.Router

// Route represents a single route
type Route = router.Route

// MiddlewareFunc is a function that wraps an http.Handler
type MiddlewareFunc = router.MiddlewareFunc

// NewRouter creates a new Router
func NewRouter() *Router {
	return router.New()
}

// GetParam retrieves a route parameter from the request
func GetParam(r *http.Request, name string) string {
	return router.GetParam(r.Context(), name)
}

// GetParams retrieves all route parameters from the request
func GetParams(r *http.Request) map[string]string {
	return router.GetParams(r.Context())
}

// Middleware helpers

// Logger is a middleware that logs requests
func Logger() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

// Recovery is a middleware that recovers from panics
func Recovery() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// MethodOverride is a middleware that allows _method parameter to override HTTP method
func MethodOverride() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				// Parse form to check for _method parameter
				r.ParseForm()
				if method := r.PostFormValue("_method"); method != "" {
					r.Method = method
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CORS is a middleware that adds CORS headers
func CORS() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// WithParams is a helper to add params to context (for testing)
func WithParams(ctx context.Context, params map[string]string) context.Context {
	return router.WithParams(ctx, params)
}
