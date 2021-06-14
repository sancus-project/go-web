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

func compileChain(chain []web.MiddlewareHandlerFunc, h web.Handler, eh web.ErrorHandlerFunc) http.Handler {
	return CompileChain(chain, intercept.Resolve(h, eh))
}

func CompileTryChain(chain []web.MiddlewareHandlerFunc, h web.Handler) web.Handler {
	if len(chain) > 0 {
		h2 := compileChain(chain, h, nil)

		h = intercept.Intercept(h2)
	}

	return h
}

// Entry handlers
type entry struct {
	h1 web.Handler
	h2 http.Handler
}

func (h entry) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.h2.ServeHTTP(w, r)
}

func (h entry) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return h.h1.TryServeHTTP(w, r)
}

func (m *Mux) compile(h web.HandlerFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.entry == nil {

		if h == nil {
			// No handler defined, but we still want to process the middleware
			f := func(w http.ResponseWriter, r *http.Request) error {
				return errors.ErrNotFound
			}

			h = web.HandlerFunc(f)
		}

		m.entry = newEntry(m.chain, h, m.errorHandler)
	}
}

func newEntry(chain []web.MiddlewareHandlerFunc, h web.Handler, eh web.ErrorHandlerFunc) Handler {
	var h2 http.Handler

	if len(chain) > 0 {
		// got a middleware chain we have to process before
		// calling the error handler
		h2 = compileChain(chain, h, eh)
		h = intercept.Intercept(h2)
	} else {
		// no middleware chain, only add resolver for standard handler
		h2 = intercept.Resolve(h, eh)
	}
	return entry{h, h2}
}
