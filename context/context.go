package context

import (
	"context"
	"log"
	"path/filepath"
	"strings"

	"go.sancus.dev/web/errors"
)

type Context struct {
	RoutePrefix  string
	RoutePath    string
	RoutePattern string
}

// Clone() creates a copy of a routing Context object
func (rctx Context) Clone() *Context {
	log.Printf("%+n: %#v", errors.Here(0), rctx)
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

	log.Printf("%+n: prefix:%q + path:%q -> %#v",
		errors.Here(0), prefix, path, rctx)
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

	log.Printf("%+n: prefix:%q + path:%q -> %q",
		errors.Here(0), prefix, path, s)
	return s
}

func (rctx *Context) Step(prefix string) *Context {
	var path string
	log.Printf("%+n.%v: %#v + prefix:%q",
		errors.Here(0), 1, rctx, prefix)

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
			errors.Here(0), prefix, rctx)
		log.Fatal(err)
	} else {
		// offset... resuming nested routing
		l := n + len(prefix)

		prefix = rctx.RoutePath[:l]
		path = rctx.RoutePath[l:]

		if prefix+s != rctx.RoutePath {
			err := errors.New("%+n: BUG: prefix:%q + path:%q != %q",
				errors.Here(0), prefix, path, rctx.RoutePath)
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

	log.Printf("%+n.%v: %#v", errors.Here(0), 2, next)
	return next
}

func (rctx *Context) Next() (*Context, string) {

	log.Printf("%+n.%v: %#v", errors.Here(0), 1, rctx)

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

		log.Printf("%+n.%v: s:%q -> %#v",
			errors.Here(0), 2, s, next)

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
