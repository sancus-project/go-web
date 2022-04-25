package errors

import (
	"net/http"

	"go.sancus.dev/core/errors"
)

type BadRequestError struct {
	errors.ErrorStack
}

func (err *BadRequestError) AsError() error {
	if err.Ok() {
		return nil
	} else {
		return err
	}
}

func (err *BadRequestError) Status() int {
	if err.Ok() {
		return http.StatusOK
	} else {
		return http.StatusBadRequest
	}
}

func (err *BadRequestError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveHTTP(err, w, r)
}

func (err *BadRequestError) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return tryServeHTTP(err, w, r)
}
