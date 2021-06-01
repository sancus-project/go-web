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

func NewFromError(err error) error {

	if err != nil {
		var code int

		switch v := err.(type) {
		case http.Handler:
			// if it can render itself it might know better
			return err
		case *PanicError, *HandlerError, *RedirectError:
			// Ours
			return err
		case web.Error:
			// Friedly
			code = v.Status()

			if CodeIsRedirect(code) {
				// Redirect with `Header() http.Header` interface
				if e, ok := err.(interface {
					Header() http.Header
				}); ok {
					if loc := e.Header().Get("Location"); loc != "" {
						//  Redirect
						return &RedirectError{
							location: loc,
							code:     code,
						}
					}

				}
				// Redirect with `Location() string` interface
				if e, ok := err.(interface {
					Location() string
				}); ok {
					if loc := e.Location(); loc != "" {
						// Redirect
						return &RedirectError{
							location: loc,
							code:     code,
						}
					}
				}
			}
			// fall through
		default:
			code = http.StatusInternalServerError
		}

		return &HandlerError{
			Code: code,
			Err:  err,
		}

	}

	return nil
}
