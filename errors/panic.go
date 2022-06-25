package errors

import (
	"net/http"

	"go.sancus.dev/core/errors"
)

type PanicError struct {
	errors.Panic
}

func (PanicError) Status() int {
	return http.StatusInternalServerError
}

func (p PanicError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveHTTP(p, w, r)
}

func Recover() errors.Panic {
	if err := errors.Recover(); err != nil {
		return &PanicError{err}
	}
	return nil
}
