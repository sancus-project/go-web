package errors

import (
	"net/http"
	"strings"

	"go.sancus.dev/web"
)

type ErrMethodNotAllowed struct {
	Method  string
	Allowed []string
}

func MethodNotAllowed(method string, allowed ...string) web.Error {
	err := &ErrMethodNotAllowed{
		Method:  method,
		Allowed: allowed,
	}
	return err
}

func (err *ErrMethodNotAllowed) Status() int {
	if err.Method == "OPTIONS" {
		return http.StatusOK
	} else {
		return http.StatusMethodNotAllowed
	}
}

func (err *ErrMethodNotAllowed) Error() string {
	return ErrorText(err.Status())
}

func (err *ErrMethodNotAllowed) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	code := err.Status()

	methods := append(err.Allowed, "OPTIONS")
	w.Header().Set("Allow", strings.Join(methods, ", "))

	if err.Method == "OPTIONS" || code == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, ErrorText(code), code)
	}
}
