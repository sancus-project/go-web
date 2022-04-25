package errors

import (
	"fmt"
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
	code := err.Status()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	fmt.Fprintln(w, ErrorText(code))
	if code != http.StatusOK {
		fmt.Fprintln(w)
		for _, e := range err.Errors() {
			fmt.Fprintln(w, e.Error())
		}
	}
}
