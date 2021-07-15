package router

import (
	"log"
	"net/http"

	"github.com/armon/go-radix"

	"go.sancus.dev/web"
	"go.sancus.dev/web/context"
	"go.sancus.dev/web/errors"
	"go.sancus.dev/web/intercept"
)

type Mux struct {
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

func (m *Mux) findBestNode(path string) (string, string, web.Handler) {
	if s, v, ok := m.trie.LongestPrefix(path); !ok {
		goto fail
	} else if h, ok := v.(web.Handler); !ok {
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

func (m *Mux) resolve(h web.Handler, rctx *context.Context, prefix, path string) (web.Handler, *context.Context, bool) {
	if rctx != nil {
		rctx = rctx.Step(prefix)
	} else {
		rctx = context.NewRouteContext(prefix, path)
	}

	return h, rctx, true
}

func (m *Mux) Resolve(path string, rctx *context.Context) (web.Handler, *context.Context, bool) {
	if s0, s1, h := m.findBestNode(path); h != nil {
		return m.resolve(h, rctx, s0, s1)
	} else {
		return nil, nil, false
	}
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

	m.TryHandle(path, h)
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

func (m *Mux) tryHandle(path string, handler web.Handler) {
	if m.entry == nil {
		m.compile(m.tryServeHTTP)
	}

	if handler != nil {
		m.trie.Insert(path, handler)
	}
}

func (m *Mux) TryHandle(path string, handler web.Handler) {
	m.tryHandle(path, handler)
}

func (m *Mux) TryHandleFunc(path string, handler web.HandlerFunc) {
	m.tryHandle(path, handler)
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
