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

type Handler interface {
	TryServeHTTP(http.ResponseWriter, *http.Request) error
}

type MiddlewareHandlerFunc func(http.Handler) http.Handler
