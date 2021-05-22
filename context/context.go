package context

import (
	"context"
	"path/filepath"
	"strings"
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

func (rctx *Context) Init(prefix, path string) {
	var pattern string

	if prefix == "" {
		prefix = "/"
	}

	if path == "" {
		pattern = prefix
	} else if n := strings.IndexRune(path[1:], '/'); n < 0 {
		pattern = filepath.Join(prefix, path)
	} else {
		pattern = filepath.Join(prefix, "*")
	}

	*rctx = Context{
		RoutePrefix:  prefix,
		RoutePattern: pattern,
		RoutePath:    path,
	}
}

func (rctx *Context) Path() string {
	prefix := rctx.RoutePrefix
	path := rctx.RoutePath

	if prefix == "/" {
		return path
	} else if path == "" {
		return prefix
	} else {
		return prefix + path
	}
}

func (rctx *Context) Next() (*Context, string) {

	path := rctx.RoutePath

	if len(path) > 1 {

		s := path[1:]
		prefix := rctx.RoutePrefix
		pattern := strings.TrimSuffix("/*", rctx.RoutePattern)

		if prefix == "/" {
			prefix = ""
		}

		if n := strings.IndexRune(s, '/'); n < 0 {
			prefix += path
			pattern += path
			path = ""

		} else {
			s = s[:n]
			prefix += path[:n+1]
			pattern += path[:n+1] + "/*"
			path = path[n+1:]
		}

		next := &Context{
			RoutePath:    path,
			RoutePrefix:  prefix,
			RoutePattern: pattern,
		}

		return next, s
	}

	return nil, ""
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
	RouteCtxKey = &contextKey{"RouteContext"}
)

// contextKey is a value for use with context.WithValue
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "mix context value " + k.name
}
