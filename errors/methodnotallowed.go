package errors

import (
	"net/http"
	"strings"
)

type MethodNotAllowedError struct {
	Method  string
	Allowed []string
}

func MethodNotAllowed(method string, allowed ...string) *MethodNotAllowedError {
	return &MethodNotAllowedError{
		Method:  strings.ToUpper(method),
		Allowed: allowed,
	}
}

func (err *MethodNotAllowedError) Status() int {
	if err.Method == "OPTIONS" {
		return http.StatusOK
	} else {
		return http.StatusMethodNotAllowed
	}
}

func (err *MethodNotAllowedError) Error() string {
	return ErrorText(err.Status())
}

func (err *MethodNotAllowedError) Methods() []string {
	return err.Allowed
}

func (err *MethodNotAllowedError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)

	w.Header().Set("Allow", strings.Join(err.Allowed, ", "))

	if method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
	} else {
		code := http.StatusMethodNotAllowed

		http.Error(w, ErrorText(code), code)
	}
}

func (err *MethodNotAllowedError) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return MethodNotAllowed(r.Method, err.Allowed...)
}
