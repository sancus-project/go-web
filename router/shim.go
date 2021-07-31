package router

import (
	"net/http"

	"go.sancus.dev/core/errors"
	"go.sancus.dev/web"
)

// a router is a trampoline for most Router methods
type router struct {
	getNode func(path string) *node
}

func (r *router) Handle(path string, h http.Handler) {
	r.getNode(path).handle(h)
}

func (r *router) HandleFunc(path string, h http.HandlerFunc) {
	r.getNode(path).handle(h)
}

func (r *router) TryHandle(path string, h web.Handler) {
	r.getNode(path).tryHandle(h)
}

func (r *router) TryHandleFunc(path string, h web.HandlerFunc) {
	r.getNode(path).tryHandle(h)
}

func (r *router) Method(method string, path string, h http.Handler) {
	r.getNode(path).method(method, h)
}

func (r *router) MethodFunc(method string, path string, h http.HandlerFunc) {
	r.getNode(path).method(method, h)
}

func (r *router) TryMethod(method string, path string, h web.Handler) {
	r.getNode(path).tryMethod(method, h)
}

func (r *router) TryMethodFunc(method string, path string, h web.HandlerFunc) {
	r.getNode(path).tryMethod(method, h)
}

func (r *router) Route(path string, fn func(Router)) Router {
	return r.getNode(path).route(fn)
}

// just a node on the trie but allowing it to
// migrate between raw and ready states
type node struct {
	Handler
}

func (n *node) toolate(fn string) {
	err := errors.New("Can't call %s() anymore", fn)
	panic(errors.WithStackTrace(2, err))
}

func (n *node) handle(h http.Handler) {
	if v, ok := n.Handler.(interface {
		handle(http.Handler)
	}); ok {
		v.handle(h)
	} else {
		n.toolate("Handle")
	}
}

func (n *node) tryHandle(h web.Handler) {
	if v, ok := n.Handler.(interface {
		tryHandle(web.Handler)
	}); ok {
		v.tryHandle(h)
	} else {
		n.toolate("TryHandle")
	}
}

func (n *node) method(method string, h http.Handler) {
	if v, ok := n.Handler.(interface {
		method(string, http.Handler)
	}); ok {
		v.method(method, h)
	} else {
		n.toolate("Method")
	}
}

func (n *node) tryMethod(method string, h web.Handler) {
	if v, ok := n.Handler.(interface {
		tryMethod(string, web.Handler)
	}); ok {
		v.tryMethod(method, h)
	} else {
		n.toolate("TryMethod")
	}
}

func (n *node) with(middleware ...web.MiddlewareHandlerFunc) {
	if v, ok := n.Handler.(interface {
		with(...web.MiddlewareHandlerFunc)
	}); ok {
		v.with(middleware...)
	} else {
		n.toolate("With")
	}
}

func (n *node) route(fn func(Router)) Router {
	if v, ok := n.Handler.(interface {
		route(func(Router)) Router
	}); ok {
		return v.route(fn)
	} else {
		n.toolate("Route")
		return nil // not reached.
	}
}
