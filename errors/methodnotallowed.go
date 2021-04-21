package errors

import (
	"net/http"
	"strings"

	"go.sancus.dev/web"
)

type MethodNotAllowedError struct {
	Method  string
	Allowed []string
}

func MethodNotAllowed(method string, allowed ...string) web.Error {
	err := &MethodNotAllowedError{
		Method:  method,
		Allowed: allowed,
	}
	return err
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

func (err *MethodNotAllowedError) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	code := err.Status()

	methods := append(err.Allowed, "OPTIONS")
	w.Header().Set("Allow", strings.Join(methods, ", "))

	if err.Method == "OPTIONS" || code == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, ErrorText(code), code)
	}
}
