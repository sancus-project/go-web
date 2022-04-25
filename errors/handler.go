package errors

import (
	"fmt"
	"net/http"

	"go.sancus.dev/core/errors"
	"go.sancus.dev/web"
)

var (
	// 400
	ErrBadRequest = &HandlerError{Code: http.StatusBadRequest}
	// 404
	ErrNotFound = &HandlerError{Code: http.StatusNotFound}
	// 406
	ErrNotAcceptable = &HandlerError{Code: http.StatusNotAcceptable}
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

func (err HandlerError) Headers() http.Header {
	return err.Header
}

func (err HandlerError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveHTTP(err, w, r)
}

func serveHTTP(err web.Error, w http.ResponseWriter, r *http.Request) {
	code := err.Status()

	// Headers
	if he, ok := err.(interface {
		Headers() http.Header
	}); ok {
		for k, v := range he.Headers() {
			switch k {
			case "Context-Type", "X-Context-Type-Options":
				// skip
			default:
				for _, s := range v {
					w.Header().Add(k, s)
				}
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

		if p, ok := err.(interface {
			Recovered() error
		}); ok {
			// Panic
			fmt.Fprintln(w)
			fmt.Fprintln(w, "panic:", p.Recovered())
		} else if p, ok := err.(interface {
			Errors() []error
		}); ok {
			// Validator
			fmt.Fprintln(w)
			for _, err := range p.Errors() {
				fmt.Fprintln(w, err.Error())
			}

		} else if p := errors.Unwrap(err); p != nil {
			// Wrapped
			fmt.Fprintln(w)
			fmt.Fprintln(w, p.Error())
		}

		// StackTrace
		if p, ok := errors.AsStackTracer(err); ok {
			fmt.Fprintln(w)
			fmt.Fprintf(w, "%#+v", p.StackTrace())
		}
	}
}

func (err HandlerError) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return tryServeHTTP(err, w, r)
}

func tryServeHTTP(err web.Error, rw http.ResponseWriter, req *http.Request) error {
	code := err.Status()

	switch code {
	case http.StatusOK, http.StatusNoContent:
		serveHTTP(err, rw, req)
		return nil
	default:
		return err
	}
}
