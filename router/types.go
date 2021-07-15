package router

import (
	"net/http"

	"go.sancus.dev/web"
)

type Handler interface {
	http.Handler
	web.Handler
}

type Router interface {
	Handler
	MiniRouter

	Use(web.MiddlewareHandlerFunc) Router
}

type MiniRouter interface {
	Handle(path string, handler http.Handler)
	HandleFunc(path string, handler http.HandlerFunc)

	TryHandle(path string, handler web.Handler)
	TryHandleFunc(path string, handler web.HandlerFunc)

	With(web.MiddlewareHandlerFunc) MiniRouter
}
