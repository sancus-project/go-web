package intercept

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

type Interceptor struct {
	h http.Handler
}

func Intercept(h http.Handler) web.Handler {
	if h != nil {
		return &Interceptor{h: h}
	}
	return nil
}

func (m *Interceptor) tryServeHTTP(w http.ResponseWriter, r *http.Request, out **errors.PanicError) {
	defer func() {
		if err := errors.Recover(); err != nil {
			*out = err
		}
	}()

	m.h.ServeHTTP(w, r)
}

func (m *Interceptor) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	var pee *errors.PanicError

	w2 := NewWriter(w, r.Method)
	m.tryServeHTTP(w2.Writer(), r, &pee)

	// panic?
	if pee != nil {
		if err, ok := pee.Unwrap().(web.Error); ok && err != nil {
			// return wrapped web.Error
			return err
		} else {
			// or just the panic
			return pee
		}
	}

	// the web.Error from the Writer
	return w2.Error()
}
