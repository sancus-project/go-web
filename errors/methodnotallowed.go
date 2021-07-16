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

func (err *MethodNotAllowedError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := err.Status()

	methods := append(err.Allowed, "OPTIONS")
	w.Header().Set("Allow", strings.Join(methods, ", "))

	if err.Method == "OPTIONS" || code == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, ErrorText(code), code)
	}
}

func (err *MethodNotAllowedError) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return err
}
