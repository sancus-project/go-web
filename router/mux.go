package router

import (
	"log"
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
	"go.sancus.dev/web/intercept"
)

type Mux struct {
	chain []web.MiddlewareHandlerFunc
	entry web.Handler

	errorHandler web.ErrorHandlerFunc
}

func NewRouter(h web.ErrorHandlerFunc) Router {
	if h == nil {
		h = errors.HandleError
	}

	return &Mux{
		errorHandler: h,
	}
}

// http.Handler
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := m.TryServeHTTP(w, r); err != nil {
		m.errorHandler(w, r, err)
	}
}

func (m *Mux) Handle(path string, handler http.Handler) {
	if handler != nil {
		var h web.Handler

		h, ok := handler.(web.Handler)
		if !ok || h == nil {
			h = intercept.Intercept(handler)
		}

		if err := m.TryHandle(path, h); err != nil {
			log.Fatal(err)
		}
	}
}

func (m *Mux) HandleFunc(path string, handler http.HandlerFunc) {
	m.Handle(path, handler)
}

// web.Handler
func (m *Mux) tryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	// TODO: resolve from trie
	return errors.ErrNotFound
}

func (m *Mux) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	if m.entry == nil {
		// no handlers
		return errors.ErrNotFound
	}

	return m.entry.TryServeHTTP(w, r)
}

func (m *Mux) tryHandle(path string, handler web.Handler) error {
	// TODO: add handler to trie
	return nil
}

func (m *Mux) TryHandle(path string, handler web.Handler) error {
	if m.entry == nil {
		// squash Use()ed middleware
		h := web.HandlerFunc(m.tryServeHTTP)
		m.entry = CompileTryChain(m.chain, h)
	}

	if handler != nil {
		return m.tryHandle(path, handler)
	}

	return nil
}

func (m *Mux) TryHandleFunc(path string, handler web.HandlerFunc) error {
	return m.TryHandle(path, handler)
}

// web.MiddlewareHandlerFunc
func (m *Mux) Use(f web.MiddlewareHandlerFunc) Router {
	if m.entry != nil {
		err := errors.New("Can't call Router.Use() after Router.Handle()")
		log.Fatal(err)
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
