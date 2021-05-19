package router

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

type Mux struct {
	chain []web.MiddlewareHandlerFunc
	entry web.Handler
}

func NewRouter() Router {
	return &Mux{}
}

// http.Handler
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := m.TryServeHTTP(w, r); err != nil {
		panic(err)
	}
}

func (m *Mux) handle(path string, handler http.Handler) {
}

func (m *Mux) Handle(path string, handler http.Handler) {
	if m.entry == nil {
		// squash Use()ed middleware
		h := web.HandlerFunc(m.tryServeHTTP)
		m.entry = CompileTryChain(m.chain, h)
	}

	m.handle(path, handler)
}

func (m *Mux) HandleFunc(path string, handler http.HandlerFunc) {
	m.Handle(path, handler)
}

// web.Handler
func (m *Mux) tryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return errors.ErrNotFound
}

func (m *Mux) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	if m.entry == nil {
		return errors.ErrNotFound
	}

	return m.entry.TryServeHTTP(w, r)
}

func (m *Mux) tryHandle(path string, handler web.Handler) error {
	return nil
}

func (m *Mux) TryHandle(path string, handler web.Handler) error {
	if m.entry == nil {
		h := web.HandlerFunc(m.tryServeHTTP)
		m.entry = CompileTryChain(m.chain, h)
	}

	return m.tryHandle(path, handler)
}

func (m *Mux) TryHandleFunc(path string, handler web.HandlerFunc) error {
	return m.TryHandle(path, handler)
}

// web.MiddlewareHandlerFunc
func (m *Mux) Use(f web.MiddlewareHandlerFunc) Router {
	if m.entry != nil {
		panic(errors.New("Can't call Router.Use() after Router.Handle()"))
	}

	if f != nil {
		m.chain = append(m.chain, f)
	}
	return m
}

func (m *Mux) With(f web.MiddlewareHandlerFunc) MiniRouter {
	chain := &Chain{mux: m}
	return chain.With(f)
}
