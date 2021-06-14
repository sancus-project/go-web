package router

import (
	"log"
	"net/http"
	"sync"

	"github.com/armon/go-radix"

	"go.sancus.dev/web"
	"go.sancus.dev/web/context"
	"go.sancus.dev/web/errors"
	"go.sancus.dev/web/intercept"
)

type Mux struct {
	mu    sync.Mutex
	chain []web.MiddlewareHandlerFunc
	trie  *radix.Tree
	entry Handler

	errorHandler web.ErrorHandlerFunc
}

func NewRouter(h web.ErrorHandlerFunc) Router {
	if h == nil {
		h = errors.HandleError
	}

	return &Mux{
		trie:         radix.New(),
		errorHandler: h,
	}
}

func (m *Mux) GetRoutePath(r *http.Request) string {
	if rctx := context.RouteContext(r.Context()); rctx != nil {
		return rctx.RoutePath
	} else {
		return r.URL.Path
	}
}

func (m *Mux) resolve(v interface{}, rctx *context.Context, prefix, path string) (web.Handler, *context.Context, bool) {
	if h, ok := v.(web.Handler); ok {

		if rctx != nil {
			rctx = rctx.Step(prefix)
		} else {
			rctx = context.NewRouteContext(prefix, path)
		}

		return h, rctx, true
	}
	return nil, nil, false
}

func (m *Mux) Resolve(path string, rctx *context.Context) (web.Handler, *context.Context, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s, v, ok := m.trie.LongestPrefix(path); ok {
		if s == path {
			return m.resolve(v, rctx, s, "")
		}

		l := len(s)
		if s[l-1] == '/' {
			return m.resolve(v, rctx, s, path[l-1:])
		} else if path[l] == '/' {
			return m.resolve(v, rctx, s, path[l:])
		}
	}

	return nil, nil, false
}

// http.Handler
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.entry == nil {
		m.compile(nil)
	}

	m.entry.ServeHTTP(w, r)
}

func (m *Mux) Handle(path string, handler http.Handler) {
	var h web.Handler

	if handler != nil {
		var ok bool
		h, ok = handler.(web.Handler)
		if !ok {
			h = intercept.Intercept(handler)
		}
	}

	if err := m.TryHandle(path, h); err != nil {
		log.Fatal(err)
	}
}

func (m *Mux) HandleFunc(path string, handler http.HandlerFunc) {
	m.Handle(path, handler)
}

// web.Handler
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

func (m *Mux) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	if m.entry == nil {
		m.compile(nil)
	}

	return m.entry.TryServeHTTP(w, r)
}

func (m *Mux) tryHandle(path string, handler web.Handler) error {
	m.trie.Insert(path, handler)
	return nil
}

func (m *Mux) TryHandle(path string, handler web.Handler) error {
	if m.entry == nil {
		m.compile(m.tryServeHTTP)
	}

	if handler != nil {
		return m.tryHandle(path, handler)
	}

	return nil
}

func (m *Mux) TryHandleFunc(path string, handler web.HandlerFunc) error {
	return m.TryHandle(path, handler)
}

// web.MiddlewareHandlerFunc
func (m *Mux) Use(f web.MiddlewareHandlerFunc) Router {
	if m.entry != nil {
		err := errors.New("Can't call Router.Use() after Router.Handle()")
		log.Fatal(err)
	}

	if f != nil {
		m.chain = append(m.chain, f)
	}
	return m
}

func (m *Mux) With(f web.MiddlewareHandlerFunc) MiniRouter {
	chain := &Chain{mux: m}
	return chain.With(f)
}
