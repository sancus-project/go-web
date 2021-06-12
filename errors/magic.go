package errors

import (
	"net/http"
	"strings"

	"go.sancus.dev/web"
)

func NewError(code int, headers http.Header, body []byte) error {

	if code < 100 {
		// 500
		code = http.StatusInternalServerError
	} else if code < 300 {
		// OK
		return nil
	}

	if CodeIsRedirect(code) {
		location := headers.Get("Location")
		return &RedirectError{
			location: location,
			code:     code,
		}
	} else if code == http.StatusMethodNotAllowed {
		allowed := headers.Get("Allowed")
		return &MethodNotAllowedError{
			Allowed: strings.Split(allowed, ", "),
		}
	} else {
		return &HandlerError{
			Code:   code,
			Header: headers,
		}
	}
}

func AsWebError(err error) web.Error {
	var p web.Error
	var ok bool

	if err == nil {
		// Ignore
	} else if p, ok = err.(web.Error); ok {
		// Ready
	} else {
		// Wrap
		p = &HandlerError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return p
}

func NewFromError(err error) error {

	if err == nil {
		// Ignore
		return nil
	} else if p, ok := AsRedirect(err); ok {
		// Redirect
		return p
	}

	var code int

	switch v := err.(type) {
	case http.Handler:
		// if it can render itself it might know better
		goto done
	case *PanicError, *HandlerError, *RedirectError:
		// Ours
		goto done
	case web.Error:
		// Friedly
		code = v.Status()
	default:
		code = http.StatusInternalServerError
	}

	// Wrap
	err = &HandlerError{
		Code: code,
		Err:  err,
	}
done:
	return err
}
