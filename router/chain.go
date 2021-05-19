package router

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

type Chain struct {
	mux   *Mux
	chain []web.MiddlewareHandlerFunc
}

// http.Handler
func (m *Chain) Handle(path string, handler http.Handler) {
	h := CompileChain(m.chain, handler)
	m.mux.Handle(path, h)
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

func CompileChain(chain []web.MiddlewareHandlerFunc, h http.Handler) http.Handler {
	l := len(chain)
	for l > 0 {
		l -= 1
		h = chain[l](h)
	}
	return h
}

func CompileTryChain(chain []web.MiddlewareHandlerFunc, h0 web.Handler) web.Handler {

	if len(chain) > 0 {
		h, ok := h0.(http.Handler)
		if !ok {
			h = errors.PanicMaker{Handler: h0}
		}

		h = CompileChain(chain, h)

		h0, ok = h.(web.Handler)
		if !ok {
			h0 = errors.PanicInterceptor{Handler: h}
		}
	}

	return h0
}
