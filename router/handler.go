package router

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
	"go.sancus.dev/web/intercept"
)

type handler struct {
	h1 web.Handler
	h2 http.Handler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.h2.ServeHTTP(w, r)
}

func (h *handler) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return h.h1.TryServeHTTP(w, r)
}

func NewHandlerFunc(h web.HandlerFunc, chain []web.MiddlewareHandlerFunc, eh web.ErrorHandlerFunc) Handler {
	return NewHandler(h, chain, eh)
}

func NewHandler(h web.Handler, chain []web.MiddlewareHandlerFunc, eh web.ErrorHandlerFunc) Handler {
	var h2 http.Handler

	if h == nil {
		// No handler defined, but we still want to process the middleware
		f := func(w http.ResponseWriter, r *http.Request) error {
			return errors.ErrNotFound
		}

		h = web.HandlerFunc(f)
	}

	if len(chain) > 0 {
		// got a middleware chain we have to process before
		// calling the error handler
		h2 = compileChain(chain, h, eh)
		h = intercept.Intercept(h2)
	} else {
		// no middleware chain, only add resolver for standard handler
		h2 = intercept.Resolve(h, eh)
	}

	if h, ok := h.(http.Handler); ok && h == h2 {
		// don't wrapper handlers needlessly
		return h.(Handler)
	}

	return &handler{
		h1: h,
		h2: h2,
	}
}
