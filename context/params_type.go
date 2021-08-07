package context

//go:generate ./params_type.sh :interface{} Int String

import (
	"go.sancus.dev/core/typeconv"
)

// Get slice of interface{} parameters from RouteContext
func (rctx *RoutingContext) GetSlice(key string) ([]interface{}, bool) {
	if v, ok := rctx.Get(key); ok {
		return typeconv.AsSlice(v)
	}

	return nil, false
}

// Get int parameter from RouteContext
func (rctx *RoutingContext) GetInt(key string) (int, bool) {
	var zero int

	if v, ok := rctx.Get(key); ok {
		return typeconv.AsInt(v)
	}

	return zero, false
}

// Get slice of int parameters from RouteContext
func (rctx *RoutingContext) GetIntSlice(key string) ([]int, bool) {
	if v, ok := rctx.Get(key); ok {
		return typeconv.AsIntSlice(v)
	}

	return nil, false
}

// Get string parameter from RouteContext
func (rctx *RoutingContext) GetString(key string) (string, bool) {
	var zero string

	if v, ok := rctx.Get(key); ok {
		return typeconv.AsString(v)
	}

	return zero, false
}

// Get slice of string parameters from RouteContext
func (rctx *RoutingContext) GetStringSlice(key string) ([]string, bool) {
	if v, ok := rctx.Get(key); ok {
		return typeconv.AsStringSlice(v)
	}

	return nil, false
}
