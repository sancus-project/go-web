package web

import (
	"net/http"
)

type ErrorHandlerFunc func(http.ResponseWriter, *http.Request, error)

type Error interface {
	Error() string
	Status() int
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (f HandlerFunc) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

type Handler interface {
	TryServeHTTP(http.ResponseWriter, *http.Request) error
}

type RendererFunc func(http.ResponseWriter, *http.Request) error

func (f RendererFunc) Render(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

type Renderer interface {
	Render(http.ResponseWriter, *http.Request) error
}

type MiddlewareHandlerFunc func(http.Handler) http.Handler

func (f MiddlewareHandlerFunc) Middleware(next http.Handler) http.Handler {
	return f(next)
}

type MiddlewareHandler interface {
	Middleware(http.Handler) http.Handler
}
