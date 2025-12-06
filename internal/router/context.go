package router

import (
	"context"
)

type contextKey string

const paramsKey contextKey = "params"

// WithParams adds params to the request context
func WithParams(ctx context.Context, params map[string]string) context.Context {
	return context.WithValue(ctx, paramsKey, params)
}

// GetParams retrieves params from the request context
func GetParams(ctx context.Context) map[string]string {
	if params, ok := ctx.Value(paramsKey).(map[string]string); ok {
		return params
	}
	return make(map[string]string)
}

// GetParam retrieves a single param from the request context
func GetParam(ctx context.Context, name string) string {
	params := GetParams(ctx)
	return params[name]
}
