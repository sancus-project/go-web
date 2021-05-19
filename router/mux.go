package router

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

type Mux struct{}

func NewRouter() Router {
	return &Mux{}
}

// http.Handler
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := m.TryServeHTTP(w, r); err != nil {
		panic(err)
	}
}

func (m *Mux) Handle(path string, handler http.Handler) {
}

func (m *Mux) HandleFunc(path string, handler http.HandlerFunc) {
	m.Handle(path, handler)
}

// web.Handler
func (m *Mux) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return errors.ErrNotFound
}

func (m *Mux) TryHandle(path string, handler web.Handler) error {
	return nil
}

func (m *Mux) TryHandleFunc(path string, handler web.HandlerFunc) error {
	return m.TryHandle(path, handler)
}
