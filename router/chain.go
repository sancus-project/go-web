package router

import (
	"net/http"

	"go.sancus.dev/web"
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
