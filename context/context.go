package context

import (
	"log"
	"path/filepath"
	"strings"

	"go.sancus.dev/core/context"
	"go.sancus.dev/core/errors"
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
	var s string

	prefix := rctx.RoutePrefix
	path := rctx.RoutePath

	if prefix == "/" {
		s = path
	} else if path == "" {
		s = prefix
	} else {
		s = prefix + path
	}

	return s
}

func (rctx *Context) Step(prefix string) *Context {
	var path string

	pattern := strings.TrimSuffix("/*", rctx.RoutePattern)

	if prefix == rctx.RoutePath {
		// prefix is the whole RoutePath
		pattern += prefix
		prefix = rctx.Path()
	} else if s := strings.TrimPrefix(rctx.RoutePath, prefix); s != rctx.RoutePath {
		// prefix is part of the RoutePath
		pattern += prefix + "/*"
		if rctx.RoutePrefix != "/" {
			prefix = rctx.RoutePrefix + prefix
		}

		path = s
	} else if n := strings.Index(rctx.RoutePath, prefix); n < 0 {
		err := errors.New("%+n: prefix:%q incompatible (%#v)",
			errors.Here(), prefix, rctx)
		log.Fatal(err)
	} else {
		// offset... resuming nested routing
		l := n + len(prefix)

		prefix = rctx.RoutePath[:l]
		path = rctx.RoutePath[l:]

		if prefix+s != rctx.RoutePath {
			err := errors.New("%+n: BUG: prefix:%q + path:%q != %q",
				errors.Here(), prefix, path, rctx.RoutePath)
			log.Fatal(err)
		}

		pattern += prefix + "/*"
		if rctx.RoutePrefix != "/" {
			prefix = rctx.RoutePrefix + prefix
		}
	}

	next := &Context{
		RoutePath:    path,
		RoutePrefix:  prefix,
		RoutePattern: pattern,
	}

	return next
}

func (rctx *Context) Next() (*Context, string) {

	path := rctx.RoutePath
	if len(path) > 1 {
		var prefix string

		s := path[1:]
		if n := strings.IndexRune(s, '/'); n < 0 {
			prefix = path
		} else {
			s = s[:n]
			prefix = path[:n+1]
		}

		next := rctx.Step(prefix)

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
	RouteCtxKey = context.NewContextKey("RouteContext")
)
