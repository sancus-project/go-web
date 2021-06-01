package errors

import (
	"fmt"
	"net/http"
)

var (
	ErrNotFound = &HandlerError{Code: http.StatusNotFound}
)

// Reference Handler error
type HandlerError struct {
	Code   int
	Err    error
	Header http.Header
}

func (err HandlerError) Status() int {
	var code int

	if err.Code != 0 {
		code = err.Code
	} else if err.Err == nil {
		code = http.StatusOK
	} else {
		code = http.StatusInternalServerError
	}

	return code
}

func (err HandlerError) Unwrap() error {
	return err.Err
}

func (err HandlerError) String() string {
	return ErrorText(err.Status())
}

func (err HandlerError) Error() string {
	return ErrorText(err.Status())
}

func (err HandlerError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := err.Status()

	if err.Header != nil {
		err.Header.Del("Context-Type")
		err.Header.Del("X-Context-Type-Options")

		for k, v := range err.Header {
			for _, s := range v {
				w.Header().Add(k, s)
			}
		}
	}

	switch code {
	case http.StatusOK, http.StatusNoContent:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(code)

		fmt.Fprintln(w, ErrorText(code))
	}
}

func (err HandlerError) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	code := err.Status()

	switch code {
	case http.StatusOK, http.StatusNoContent:
		w.WriteHeader(http.StatusNoContent)
		return nil
	default:
		return err
	}
}
