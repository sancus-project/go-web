package errors

import (
	"fmt"
	"net/http"
	"strings"
)

type BadRequestError struct {
	Errors []error
}

func (err *BadRequestError) AppendError(e error) {
	err.Errors = append(err.Errors, e)
}

func (err *BadRequestError) AppendErrorString(s string) {
	err.AppendError(New(s))
}

func (err *BadRequestError) AppendErrorf(s string, args ...interface{}) {
	err.AppendError(New(s, args...))
}

func (err *BadRequestError) Ok() bool {
	return len(err.Errors) == 0
}

func (err *BadRequestError) Status() int {
	if len(err.Errors) == 0 {
		return http.StatusOK
	} else {
		return http.StatusBadRequest
	}
}

func (err *BadRequestError) String() string {
	var errors []string
	for _, e := range err.Errors {
		errors = append(errors, e.Error())
	}
	return strings.Join(errors, "\n")
}

func (err *BadRequestError) Error() string {
	return ErrorText(err.Status())
}

func (err *BadRequestError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := err.Status()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	fmt.Fprintln(w, ErrorText(code))
	if code != http.StatusOK {
		fmt.Fprintln(w)
		for _, e := range err.Errors {
			fmt.Fprintln(w, e.Error())
		}
	}
}
