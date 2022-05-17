package router

import (
	"strings"

	"go.sancus.dev/web"
	"go.sancus.dev/web/context"
	"go.sancus.dev/web/errors"
)

func (m *Mux) findBestNode(path string) (string, string, *node) {
	if s, v, ok := m.trie.LongestPrefix(path); ok {
		if h, ok := v.(*node); !ok {
			// wtf, how did this get in the trie?
			panic(errors.New("bad node at %q (%T)", s, v))
		} else if s == path {
			// exact match
			return s, "", h
		} else if strings.HasSuffix(h.Pattern, "/*") {
			// match for foo/* patterns

			// `/*` special case
			if s == "/" {
				s = ""
			}

			if l := len(s); path[l] == '/' {
				// good match
				return s, path[l:], h
			}
		}
	}

	// fail
	return "", "", nil
}

// Resolve finds the best handler for a path and returns the corresponding RouteContext
func (m *Mux) Resolve(path string, rctx *context.RoutingContext) (web.Handler, *context.RoutingContext, bool) {
	s0, s1, h := m.findBestNode(path)

	if h == nil {
		return nil, nil, false
	}

	if s0 != "/" && s1 == "" && strings.HasSuffix(h.Pattern, "/*") {
		// redirect to the root of the subrouter
		return errors.NewPermanentRedirect("%s/", rctx.Path()), rctx, true
	}

	if rctx != nil {
		rctx = rctx.Step(s0, nil)
	} else {
		rctx = context.NewRouteContext(s0, s1)
	}

	return h, rctx, true
}
