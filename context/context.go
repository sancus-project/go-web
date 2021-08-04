package context

import (
	"go.sancus.dev/core/context"
)

type Context struct {
	RoutePrefix  string
	RoutePath    string
	RoutePattern string
}

// Clone() creates a copy of a routing Context object
func (rctx Context) Clone() *Context {
	return &rctx
}

// RouteContext returns mix's routing Context object from a
// http.Request Context.
func RouteContext(ctx context.Context) *Context {

	if rctx, ok := ctx.Value(RouteCtxKey).(*Context); ok {
		return rctx
	}
	return nil
}

// NewRouteContext returns a new routing Context object.
func NewRouteContext(prefix, path string) *Context {
	rctx := &Context{}
	rctx.Init(prefix, path)
	return rctx
}

// WithRouteContext returns a new http.Request Context with
// a given mix routing Context object connected to it, so it
// can later be extracted using RouteContext()
func WithRouteContext(ctx context.Context, rctx *Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if rctx == nil {
		rctx = &Context{}
	}
	return context.WithValue(ctx, RouteCtxKey, rctx)
}

var (
	// RouteCtxKey is the context.Context key to store the request context.
	RouteCtxKey = context.NewContextKey("RouteContext")
)
