package router

import (
	"net/http"
	"sort"
	"strings"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

type MethodHandler struct {
	handler map[string]web.Handler
}

func NewMethodHandler(fallback web.Handler) *MethodHandler {
	m := &MethodHandler{}

	if fallback != nil {
		m.set("*", fallback)
	}

	return m
}

func (m *MethodHandler) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	var err error

	method := strings.ToUpper(r.Method)

	if h, ok := m.handler[method]; ok {
		err = h.TryServeHTTP(w, r)
	} else if h, ok := m.handler["*"]; ok {
		err = h.TryServeHTTP(w, r)
	} else {
		err = m.MethodNotAllowed(r)
	}

	if wee := errors.NewFromError(err); wee == nil {
		return wee
	}
	return nil
}

func (m *MethodHandler) MethodNotAllowed(r *http.Request) error {
	var allowed []string

	// this method could be called in parallel, so we
	// play it safe and try to avoid to get "OPTIONS"
	// duplicated
	h := m.handler["OPTIONS"]

	if v, ok := h.(*errors.MethodNotAllowedError); ok {
		allowed = v.Allowed
	} else if v, ok := h.(interface {
		Methods() []string
	}); ok {
		allowed = v.Methods()
	} else {
		allowed = make([]string, 0, len(m.handler)+1)
		allowed = append(allowed, "OPTIONS")

		for k := range m.handler {
			if k != "*" && k != "OPTIONS" {
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
	}

	v := errors.MethodNotAllowed(r.Method, allowed...)
	if h == nil {
		// memoize
		m.handler["OPTIONS"] = v
	}

	return v
}

func (m *MethodHandler) set(method string, h0 web.Handler, chain ...web.MiddlewareHandlerFunc) {
	if m.handler == nil {
		m.handler = make(map[string]web.Handler, 2)
	}

	if h0 != nil {
		h := NewHandler(h0, chain, nil)

		method = strings.ToUpper(method)
		m.handler[method] = h
		if method == "GET" {
			if _, ok := m.handler["HEAD"]; !ok {
				m.handler["HEAD"] = h
			}
		}
	}
}
