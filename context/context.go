package context

import (
	"go.sancus.dev/core/context"
)

type RoutingContext struct {
	RoutePrefix  string
	RoutePath    string
	RoutePattern string
}

// Clone() creates a copy of a RoutingContext object
func (rctx RoutingContext) Clone() *RoutingContext {
	return &rctx
}

// RouteContext returns a RoutingContext object from a
// http.Request Context.
func RouteContext(ctx context.Context) *RoutingContext {

	if rctx, ok := ctx.Value(RouteCtxKey).(*RoutingContext); ok {
		return rctx
	}
	return nil
}

// NewRouteContext returns a new RoutingContext object.
func NewRouteContext(prefix, path string) *RoutingContext {
	rctx := &RoutingContext{}
	rctx.Init(prefix, path)
	return rctx
}

// WithRouteContext returns a new http.Request Context with
// a given mix routing Context object connected to it, so it
// can later be extracted using RouteContext()
func WithRouteContext(ctx context.Context, rctx *RoutingContext) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if rctx == nil {
		rctx = &RoutingContext{}
	}
	return context.WithValue(ctx, RouteCtxKey, rctx)
}

var (
	// RouteCtxKey is the context.Context key to store the request context.
	RouteCtxKey = context.NewContextKey("RouteContext")
)
