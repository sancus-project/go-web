package router

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
	"go.sancus.dev/web/intercept"
)

type Chain struct {
	mux   *Mux
	chain []web.MiddlewareHandlerFunc
}

// http.Handler
func (m *Chain) Handle(path string, handler http.Handler) {
	handler = CompileChain(m.chain, handler)
	h, ok := handler.(web.Handler)
	if !ok {
		h = intercept.Intercept(handler)
	}
	m.mux.TryHandle(path, h)
}

func (m *Chain) HandleFunc(path string, handler http.HandlerFunc) {
	m.Handle(path, handler)
}

// web.Handler
func (m *Chain) TryHandle(path string, handler web.Handler) error {
	h := CompileTryChain(m.chain, handler)
	return m.mux.TryHandle(path, h)
}

func (m *Chain) TryHandleFunc(path string, handler web.HandlerFunc) error {
	return m.TryHandle(path, handler)
}

// web.MiddlewareHandlerFunc
func (m *Chain) With(f web.MiddlewareHandlerFunc) MiniRouter {
	if f != nil {
		m.chain = append(m.chain, f)
	}
	return m
}

// Squash middleware chain
func CompileChain(chain []web.MiddlewareHandlerFunc, h http.Handler) http.Handler {
	l := len(chain)
	for l > 0 {
		l -= 1
		h = chain[l](h)
	}
	return h
}

func CompileTryChain(chain []web.MiddlewareHandlerFunc, h web.Handler) web.Handler {
	l := len(chain)
	if l > 0 {
		h2, ok := h.(http.Handler)
		if !ok {
			// use fallback error handler to minimize writing
			h2 = intercept.Resolve(h, nil)
		}
		h2 = CompileChain(chain, h2)

		h, ok = h2.(web.Handler)
		if !ok {
			h = intercept.Intercept(h2)
		}
	}

	return h
}

// Entry handlers
type Entry struct {
	web.Handler

	m *Mux // for the error handler
}

func (e Entry) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := e.TryServeHTTP(w, r); err != nil {
		e.m.errorHandler(w, r, err)
	}
}

func (m *Mux) compileFunc(h web.HandlerFunc) {
	m.compile(&Entry{h, m})
}

func (m *Mux) compile(h Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.entry == nil {
		var ok bool

		if h == nil {
			// No handler defined, but we still want to process the middleware
			h = errors.ErrNotFound
		}

		// squash middleware
		h2 := CompileTryChain(m.chain, h)

		// and prepare Mux's entry handler
		if h, ok = h2.(Handler); !ok {
			h = &Entry{h2, m}
		}

		m.entry = h
	}
}
