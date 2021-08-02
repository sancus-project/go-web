package resource

import (
	"net/http"
	"sort"
	"strings"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

type Resource struct {
	h  map[string]web.HandlerFunc
	eh web.ErrorHandlerFunc
}

func (m *Resource) TryServeHTTP(rw http.ResponseWriter, req *http.Request) error {
	h, ok := m.h[req.Method]
	if !ok {
		h = m.h["OPTIONS"]
	}

	return h(rw, req)
}

func (m *Resource) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if err := m.TryServeHTTP(rw, req); err != nil {
		m.eh(rw, req, err)
	}
}

func (m *Resource) Init(v interface{}, eh web.ErrorHandlerFunc) {
	if eh == nil {
		eh = errors.HandleError
	}

	*m = Resource{
		h:  make(map[string]web.HandlerFunc),
		eh: eh,
	}

	// GET
	if p, ok := v.(interface {
		Get(http.ResponseWriter, *http.Request) error
	}); ok {
		m.h["GET"] = p.Get
	}

	// HEAD
	if p, ok := v.(interface {
		Head(http.ResponseWriter, *http.Request) error
	}); ok {
		m.h["HEAD"] = p.Head
	} else if fn, ok := m.h["GET"]; ok {
		m.h["HEAD"] = fn
	}

	// POST
	if p, ok := v.(interface {
		Post(http.ResponseWriter, *http.Request) error
	}); ok {
		m.h["POST"] = p.Post
	}

	// OPTIONS
	if p, ok := v.(interface {
		Options(http.ResponseWriter, *http.Request) error
	}); ok {
		m.h["OPTIONS"] = p.Options
	} else {
		var allowed []string

		for k := range m.h {
			allowed = append(allowed, k)
		}
		allowed = append(allowed, "OPTIONS")

		sort.Strings(allowed)

		fn := func(rw http.ResponseWriter, req *http.Request) error {

			if strings.ToUpper(req.Method) == "OPTIONS" {
				// render
				rw.Header().Set("Allow", strings.Join(allowed, ", "))
				rw.WriteHeader(http.StatusNoContent)
				return nil
			} else {
				return errors.MethodNotAllowed(req.Method, allowed...)
			}
		}

		m.h["OPTIONS"] = fn
	}
}

func NewResource(v interface{}, eh web.ErrorHandlerFunc) *Resource {
	m := &Resource{}
	m.Init(v, eh)
	return m
}
