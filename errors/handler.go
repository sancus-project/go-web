package errors

import (
	"net/http"

	"go.sancus.dev/web"
)

var (
	// 404
	ErrNotFound = &HandlerError{Code: http.StatusNotFound}

	// interfaces
	_ http.Handler = (*HandlerError)(nil)
	_ web.Handler  = (*HandlerError)(nil)
	_ web.Error    = (*HandlerError)(nil)
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

// Serve Error as HTTP Response
func (err HandlerError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveHTTP(err, w, r)
}

func serveHTTP(err web.Error, w http.ResponseWriter, r *http.Request) {
	AsDescriptor(err).ServeHTTP(w, r)
}

// Serve Error as HTTP Response if it's not an actual error
func (err HandlerError) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return tryServeHTTP(err, w, r)
}

func tryServeHTTP(err web.Error, rw http.ResponseWriter, req *http.Request) error {
	code := err.Status()

	switch {
	case code < 300:
		// Success
		AsDescriptor(err).ServeHTTP(rw, req)
		return nil
	default:
		// Error, avoid transforming to ErrorDescriptor unnecessarily
		return err
	}
}
