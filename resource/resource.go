package resource

import (
	"net/http"
	"sort"
	"strings"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

type Resource struct {
	h     map[string]web.HandlerFunc
	eh    web.ErrorHandlerFunc
	check ContextChecker
}

func (m *Resource) TryServeHTTP(rw http.ResponseWriter, req *http.Request) error {
	if m.check != nil {
		ctx0 := req.Context()
		if ctx1, err := m.check(ctx0); err != nil {
			// incorrect context, fail
			return err
		} else if ctx0 != ctx1 {
			// new context, update request
			req = req.WithContext(ctx1)
		}
	}

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

func (m *Resource) Init(v interface{}, eh web.ErrorHandlerFunc, check ContextChecker) {
	if eh == nil {
		eh = errors.HandleError
	}

	if check == nil {
		if p, ok := v.(Checker); ok {
			check = p.Check
		} else {
			check = DefaultResourceChecker
		}
	}

	*m = Resource{
		h:     make(map[string]web.HandlerFunc),
		eh:    eh,
		check: check,
	}

	// GET
	if p, ok := v.(Getter); ok {
		m.h["GET"] = p.Get
	}

	// HEAD
	if p, ok := v.(Peeker); ok {
		m.h["HEAD"] = p.Head
	} else if fn, ok := m.h["GET"]; ok {
		m.h["HEAD"] = fn
	}

	// POST
	if p, ok := v.(Poster); ok {
		m.h["POST"] = p.Post
	}

	// DELETE
	if p, ok := v.(Deleter); ok {
		m.h["DELETE"] = p.Delete
	}

	// OPTIONS
	if p, ok := v.(Optioner); ok {
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

func NewResource(v interface{}, eh web.ErrorHandlerFunc, check ContextChecker) *Resource {
	m := &Resource{}
	m.Init(v, eh, check)
	return m
}
