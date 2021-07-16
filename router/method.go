package router

import (
	"net/http"
	"sort"
	"strings"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
	"go.sancus.dev/web/intercept"
)

type MethodHandler struct {
	handler map[string]web.Handler
	allowed []string
}

func (m *MethodHandler) init() {
	m.handler = make(map[string]web.Handler, 2)
}

func (m *MethodHandler) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	if h, ok := m.handler[strings.ToUpper(r.Method)]; ok {
		return h.TryServeHTTP(w, r)
	} else if h, ok := m.handler["*"]; ok {
		return h.TryServeHTTP(w, r)
	} else {
		return m.Options(w, r)
	}
}

func (m *MethodHandler) Options(w http.ResponseWriter, r *http.Request) error {
	allowed := m.allowed

	if allowed == nil {
		// memoize
		allowed = make([]string, 0, len(m.handler))
		for k, _ := range m.handler {
			if k != "*" {
				allowed = append(allowed, k)
			}
		}

		if _, ok := m.handler["*"]; ok {
			// catch all
			all := []string{"GET", "HEAD", "PUT", "DELETE"}
			for _, k := range all {
				if _, ok := m.handler[k]; !ok {
					allowed = append(allowed, k)
				}
			}
		}
		sort.Strings(allowed)
		m.allowed = allowed
	}

	if err := errors.MethodNotAllowed(r.Method, allowed...); err.Method == "OPTIONS" {
		// render
		err.ServeHTTP(w, r)
		return nil
	} else {
		// Method Not Allowed
		return err
	}
}

func (m *MethodHandler) Method(method string, h http.Handler) {
	h2, ok := h.(web.Handler)
	if !ok {
		h2 = intercept.Intercept(h)
	}
	m.TryMethod(method, h2)
}

func (m *MethodHandler) TryMethod(method string, h web.Handler) {
	if m.handler == nil {
		m.init()
	}

	if h != nil {
		method = strings.ToUpper(method)
		m.handler[method] = h
		if method == "GET" {
			if _, ok := m.handler["HEAD"]; !ok {
				m.handler["HEAD"] = h
			}
		}
	}
}

//
// Add Method Handler to Router
//
func (m *Mux) getMethodHandler(path string, chain ...web.MiddlewareHandlerFunc) *MethodHandler {
	var h *MethodHandler

	if h2 := m.findHandler(path); h2 != nil {
		var ok bool

		// node exists
		h, ok = h2.(*MethodHandler)
		if !ok {
			// but not a MethodHandler, wrap it
			h = &MethodHandler{}
			h.TryMethod("*", h2)

			m.tryHandle(path, CompileTryChain(chain, h))
		}
	} else {
		// new node
		h = &MethodHandler{}
		m.tryHandle(path, CompileTryChain(chain, h))
	}

	return h
}

func (m *Mux) Method(method string, path string, h http.Handler) {
	m.getMethodHandler(path).Method(method, h)
}

func (m *Mux) MethodFunc(method string, path string, h http.HandlerFunc) {
	m.getMethodHandler(path).Method(method, h)
}

func (m *Mux) TryMethod(method string, path string, h web.Handler) {
	m.getMethodHandler(path).TryMethod(method, h)
}

func (m *Mux) TryMethodFunc(method string, path string, h web.HandlerFunc) {
	m.getMethodHandler(path).TryMethod(method, h)
}

//
// Add Method Handler to Chain
//
func (m *Chain) getMethodHandler(path string) *MethodHandler {
	return m.mux.getMethodHandler(path, m.chain...)
}

func (m *Chain) Method(method string, path string, h http.Handler) {
	m.getMethodHandler(path).Method(method, h)
}

func (m *Chain) TryMethod(method string, path string, h web.Handler) {
	m.getMethodHandler(path).TryMethod(method, h)
}

func (m *Chain) MethodFunc(method string, path string, h http.HandlerFunc) {
	m.getMethodHandler(path).Method(method, h)
}

func (m *Chain) TryMethodFunc(method string, path string, h web.HandlerFunc) {
	m.getMethodHandler(path).TryMethod(method, h)
}
