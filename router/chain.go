package router

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

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
