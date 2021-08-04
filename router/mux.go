package router

import (
	"net/http"

	"github.com/armon/go-radix"

	"go.sancus.dev/web"
	"go.sancus.dev/web/context"
	"go.sancus.dev/web/errors"
)

type Mux struct {
	router
	node

	trie         *radix.Tree
	errorHandler web.ErrorHandlerFunc
}

func NewRouter(h web.ErrorHandlerFunc) Router {
	if h == nil {
		h = errors.HandleError
	}

	m := &Mux{
		trie:         radix.New(),
		errorHandler: h,
	}

	// set entrypoint, but wait for middleware
	m.node.initEntry(m, m.tryServeHTTP)

	m.router.getNode = m.getNode
	return m
}

func (m *Mux) GetRoutePath(r *http.Request) string {
	if rctx := context.RouteContext(r.Context()); rctx != nil {
		return rctx.RoutePath
	} else {
		return r.URL.Path
	}
}

func (m *Mux) findBestNode(path string) (string, string, *node) {
	if s, v, ok := m.trie.LongestPrefix(path); !ok {
		goto fail
	} else if h, ok := v.(*node); !ok {
		goto fail
	} else {

		if s == path {
			return s, "", h
		}

		l := len(s)
		if s[l-1] == '/' {
			return s, path[l-1:], h
		} else if path[l] == '/' {
			return s, path[l:], h
		}
	}

fail:
	return "", "", nil
}

// getNode() is only called by Router methods to populate the tree
func (m *Mux) getNode(path string) *node {
	// to prevend semantic mess-ups we compile the entry point
	// so Use() can't be used after Handle()
	if v, ok := m.node.Handler.(interface {
		compile()
	}); ok {
		v.compile()
	}

	// reuse node when there is a match
	if _, s1, h := m.findBestNode(path); h != nil && len(s1) == 0 {
		return h
	}

	// or create a new one
	n := &node{}
	n.initRaw(m)
	m.trie.Insert(path, n)
	return n
}

// resolve updates the RouteContext for each http.Request
func (m *Mux) resolve(h web.Handler, rctx *context.RoutingContext, prefix, path string) (web.Handler, *context.RoutingContext, bool) {
	if rctx != nil {
		rctx = rctx.Step(prefix)
	} else {
		rctx = context.NewRouteContext(prefix, path)
	}

	return h, rctx, true
}

// Resolve finds the best handler for a path and returns the corresponding RouteContext, prefix, and path
func (m *Mux) Resolve(path string, rctx *context.RoutingContext) (web.Handler, *context.RoutingContext, bool) {
	if s0, s1, h := m.findBestNode(path); h != nil {
		return m.resolve(h, rctx, s0, s1)
	} else {
		return nil, nil, false
	}
}

// tryServeHTTP is the Router's entrypoint
func (m *Mux) tryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	path := m.GetRoutePath(r)
	ctx := r.Context()

	rctx := context.RouteContext(ctx)

	if h, rctx, ok := m.Resolve(path, rctx); ok {

		ctx = context.WithRouteContext(ctx, rctx)
		r = r.WithContext(ctx)

		return h.TryServeHTTP(w, r)
	}

	return errors.ErrNotFound
}

// Use appends middleware to the entrypoint of the Router
func (m *Mux) Use(f web.MiddlewareHandlerFunc) Router {
	if v, ok := m.node.Handler.(interface {
		use(web.MiddlewareHandlerFunc)
	}); ok {
		v.use(f)
	} else {
		m.node.toolate("Use")
	}

	return m
}

func (m *Mux) With(f web.MiddlewareHandlerFunc) MiniRouter {
	chain := &Chain{}
	chain.init(m)
	return chain.With(f)
}
