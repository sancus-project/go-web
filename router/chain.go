package router

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/intercept"
)

type Chain struct {
	router

	mux   *Mux
	chain []web.MiddlewareHandlerFunc
}

func (m *Chain) init(mux *Mux) {
	m.mux = mux
	m.router.getNode = m.getNode
}

func (m *Chain) getNode(path string) *node {
	n := m.mux.getNode(path)
	n.with(m.chain...)
	return n
}

func (m *Chain) With(f web.MiddlewareHandlerFunc) MiniRouter {
	if f != nil {
		m2 := &Chain{
			chain: append(m.chain, f),
		}
		m2.init(m.mux)
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
