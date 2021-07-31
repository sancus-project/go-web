package router

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/intercept"
)

// the entry node is special node, it only takes Use()
type entry struct {
	h     web.Handler
	node  *node
	mux   *Mux
	chain []web.MiddlewareHandlerFunc
}

func (n *node) initEntry(mux *Mux, h web.HandlerFunc) {
	n.Handler = &entry{
		h:    h,
		node: n,
		mux:  mux,
	}
}

func (n *entry) compile() Handler {
	h := NewHandler(n.h, n.chain, n.mux.errorHandler)
	n.node.Handler = h
	return h
}

func (n *entry) TryServeHTTP(rw http.ResponseWriter, req *http.Request) error {
	return n.compile().TryServeHTTP(rw, req)
}

func (n *entry) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	n.compile().ServeHTTP(rw, req)
}

func (n *entry) use(f web.MiddlewareHandlerFunc) {
	if f != nil {
		n.chain = append(n.chain, f)
	}
}

// a rawNode is a node that hasn't been compiled into a final Handler yet
type rawNode struct {
	h     web.Handler
	node  *node
	chain []web.MiddlewareHandlerFunc
}

func (n *node) initRaw(mux *Mux) {
	n.Handler = &rawNode{
		node: n,
	}
}

func (n *rawNode) compile() Handler {
	h := NewHandler(n.h, nil, nil)
	n.node.Handler = h
	return h
}

func (n *rawNode) TryServeHTTP(rw http.ResponseWriter, req *http.Request) error {
	return n.compile().TryServeHTTP(rw, req)
}

func (n *rawNode) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	n.compile().ServeHTTP(rw, req)
}

func (n *rawNode) handle(h2 http.Handler) {
	n.tryHandle(n.asHandler(h2))
}

func (n *rawNode) method(method string, h http.Handler) {
	n.tryMethod(method, n.asHandler(h))
}

func (n *rawNode) tryHandle(h web.Handler) {
	if v, ok := n.h.(*MethodHandler); ok {
		v.set("*", h, n.chain...)
	} else {
		n.h = NewHandler(h, n.chain, nil)
	}
}

func (n *rawNode) tryMethod(method string, h web.Handler) {
	v, ok := n.h.(*MethodHandler)
	if !ok {
		v = NewMethodHandler(n.h)
		n.h = v
	}

	v.set(method, h, n.chain...)
}

func (n *rawNode) with(chain ...web.MiddlewareHandlerFunc) {
	n.chain = chain
}

func (_ rawNode) asHandler(h2 http.Handler) web.Handler {
	if h2 == nil {
		return nil
	} else if h, ok := h2.(web.Handler); ok {
		return h
	} else {
		return intercept.Intercept(h2)
	}
}
