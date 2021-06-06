package errors

import (
	"fmt"
	"net/http"

	"go.sancus.dev/core/errors"
)

type PanicError struct {
	errors.Panic
}

func (_ PanicError) Status() int {
	return http.StatusInternalServerError
}

func (p PanicError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := p.Status()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	fmt.Fprintf(w, "%s (Error %v)\n\n", http.StatusText(code), code)
	fmt.Fprintln(w, "panic:", p.Recovered())
	fmt.Fprintln(w)

	fmt.Fprintf(w, "%#+v", p.StackTrace())
}

func Recover() errors.Panic {
	if err := errors.Recover(); err != nil {
		return &PanicError{err}
	}
	return nil
}
