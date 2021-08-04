package context

import (
	"log"
	"path/filepath"
	"strings"

	"go.sancus.dev/core/errors"
)

func (rctx *RoutingContext) Init(prefix, path string) {
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

	*rctx = RoutingContext{
		RoutePrefix:  prefix,
		RoutePattern: pattern,
		RoutePath:    path,
	}
}

func (rctx *RoutingContext) Path() string {
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

func (rctx *RoutingContext) Step(prefix string) *RoutingContext {
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

	next := &RoutingContext{
		RoutePath:    path,
		RoutePrefix:  prefix,
		RoutePattern: pattern,
	}

	return next
}

func (rctx *RoutingContext) Next() (*RoutingContext, string) {

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
