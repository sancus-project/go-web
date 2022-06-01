package intercept

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

const (
	DefaultReadBufferSize = 4096
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

func (m *Interceptor) tryServeHTTP(w http.ResponseWriter, r *http.Request, out *errors.Panic) {
	defer func() {
		if err := errors.Recover(); err != nil {
			*out = err
		}
	}()

	m.h.ServeHTTP(w, r)
}

func (m *Interceptor) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	var pee errors.Panic
	var err web.Error

	// error context and intercept.Writer
	ctx := errors.WithErrorContext(r.Context(), &err)
	r2 := r.WithContext(ctx)
	w2 := NewWriter(w, r2.Method)

	// try/recover
	m.tryServeHTTP(w2.Writer(), r2, &pee)

	// panic?
	if pee != nil {
		if err, ok := errors.Unwrap(pee).(web.Error); ok && err != nil {
			// return wrapped web.Error
			return err
		} else {
			// or just the panic
			return pee
		}
	}

	if err == nil {
		// no error via context, ask the Writer
		err = w2.Error()
	}

	return err
}
