package router

import (
	"net/http"
	"regexp"

	"github.com/armon/go-radix"

	"go.sancus.dev/web"
	"go.sancus.dev/web/context"
	"go.sancus.dev/web/errors"
)

type Mux struct {
	router
	node

	trie         *radix.Tree
	pattern      map[*regexp.Regexp]*node
	errorHandler web.ErrorHandlerFunc
}

func NewRouter(h web.ErrorHandlerFunc) Router {
	if h == nil {
		h = errors.HandleError
	}

	m := &Mux{
		trie:         radix.New(),
		pattern:      make(map[*regexp.Regexp]*node),
		errorHandler: h,
	}

	// set entrypoint, but wait for middleware
	m.node.initEntry(m, m.tryServeHTTP)

	m.router.getNode = m.getNode
	return m
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

	if p, err := m.parsePath(path); err != nil {
		panic(err)
	} else if p.Literal() {
		// reuse node when there is a match
		path = p.Path()
		if _, s1, h := m.findBestNode(path, nil); h != nil && len(s1) == 0 {
			return h
		}

		// or create a new one
		n := &node{
			Pattern: p.Pattern(),
		}

		n.initRaw(m)
		m.trie.Insert(path, n)
		return n
	} else {
		// reuse node when there is a match
		pattern := p.Pattern()
		for _, n := range m.pattern {
			if n.Pattern == pattern {
				return n
			}
		}

		// or create a new one
		re, err := p.Compile()
		if err != nil {
			panic(err)
		}

		n := &node{
			Pattern: pattern,
		}
		n.initRaw(m)
		m.pattern[re] = n

		return n
	}
}

// tryServeHTTP is the Router's entrypoint
func (m *Mux) tryServeHTTP(w http.ResponseWriter, r *http.Request) error {

	// get (or create) RoutingContext and the corresponding Route Path
	ctx, rctx, path := context.GetRouteContextPath(r)

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
