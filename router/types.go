package router

import (
	"net/http"

	"go.sancus.dev/web"
)

type Router interface {
	http.Handler
	web.Handler

	MiniRouter

	Use(web.MiddlewareHandlerFunc) Router
}

type MiniRouter interface {
	Handle(path string, handler http.Handler)
	HandleFunc(path string, handler http.HandlerFunc)

	TryHandle(path string, handler web.Handler) error
	TryHandleFunc(path string, handler web.HandlerFunc) error

	With(web.MiddlewareHandlerFunc) MiniRouter
}
