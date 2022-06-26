package errors

//go:generate ./errorstack.sh BadRequest NotAcceptable

import (
	"net/http"

	"go.sancus.dev/core/errors"
	"go.sancus.dev/web/tools"
)

var (
	// Constant http.StatusBadRequest HandlerError
	ErrBadRequest = &HandlerError{Code: http.StatusBadRequest}
	// Constant http.StatusNotAcceptable HandlerError
	ErrNotAcceptable = &HandlerError{Code: http.StatusNotAcceptable}
)

type BadRequestError struct {
	errors.ErrorStack

	Header http.Header
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

func (err *BadRequestError) Headers() http.Header {
	if err.Header == nil {
		err.Header = make(map[string][]string)
	}
	return err.Header
}

func (err *BadRequestError) WithHeaders(hdr http.Header) *BadRequestError {
	tools.CopyHeaders(err.Headers(), hdr)
	return err
}

func (err *BadRequestError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveHTTP(err, w, r)
}

func (err *BadRequestError) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return tryServeHTTP(err, w, r)
}

func BadRequest(errs ...error) *BadRequestError {
	return &BadRequestError{
		ErrorStack: errors.NewErrorStack(errs...),
	}
}

type NotAcceptableError struct {
	errors.ErrorStack

	Header http.Header
}

func (err *NotAcceptableError) AsError() error {
	if err.Ok() {
		return nil
	} else {
		return err
	}
}

func (err *NotAcceptableError) Status() int {
	if err.Ok() {
		return http.StatusOK
	} else {
		return http.StatusNotAcceptable
	}
}

func (err *NotAcceptableError) Headers() http.Header {
	if err.Header == nil {
		err.Header = make(map[string][]string)
	}
	return err.Header
}

func (err *NotAcceptableError) WithHeaders(hdr http.Header) *NotAcceptableError {
	tools.CopyHeaders(err.Headers(), hdr)
	return err
}

func (err *NotAcceptableError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveHTTP(err, w, r)
}

func (err *NotAcceptableError) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return tryServeHTTP(err, w, r)
}

func NotAcceptable(errs ...error) *NotAcceptableError {
	return &NotAcceptableError{
		ErrorStack: errors.NewErrorStack(errs...),
	}
}
