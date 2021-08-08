package context

//go:generate ./params_type.sh :interface{} Int String

import (
	"go.sancus.dev/core/typeconv"
)

// Get slice of interface{} parameters from RouteContext
func (rctx *RoutingContext) GetSlice(key string) (x []interface{}, err error, ok bool) {
	v, _, ok := rctx.Get(key)
	if ok {
		x, err = typeconv.AsSlice(v)
	}
	return
}

// Get int parameter from RouteContext
func (rctx *RoutingContext) GetInt(key string) (x int, err error, ok bool) {
	v, _, ok := rctx.Get(key)
	if ok {
		x, err = typeconv.AsInt(v)
	}
	return
}

// Get slice of int parameters from RouteContext
func (rctx *RoutingContext) GetIntSlice(key string) (x []int, err error, ok bool) {
	v, _, ok := rctx.Get(key)
	if ok {
		x, err = typeconv.AsIntSlice(v)
	}
	return
}

// Get string parameter from RouteContext
func (rctx *RoutingContext) GetString(key string) (x string, err error, ok bool) {
	v, _, ok := rctx.Get(key)
	if ok {
		x, err = typeconv.AsString(v)
	}
	return
}

// Get slice of string parameters from RouteContext
func (rctx *RoutingContext) GetStringSlice(key string) (x []string, err error, ok bool) {
	v, _, ok := rctx.Get(key)
	if ok {
		x, err = typeconv.AsStringSlice(v)
	}
	return
}
