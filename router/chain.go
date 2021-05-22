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
	h := CompileChain(m.chain, handler)
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

func CompileChain(chain []web.MiddlewareHandlerFunc, h http.Handler) web.Handler {
	l := len(chain)
	for l > 0 {
		l -= 1
		h = chain[l](h)
	}
	return intercept.Intercept(h)
}

func CompileTryChain(chain []web.MiddlewareHandlerFunc, h web.Handler) web.Handler {
	l := len(chain)
	if l > 0 {
		// use fallback error handler to minimize writing
		h2 := intercept.Resolve(h, nil)
		h = CompileChain(chain, h2)
	}

	return h
}
