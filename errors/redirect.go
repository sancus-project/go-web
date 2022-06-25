package errors

import (
	"fmt"
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/tools"
)

func CodeIsRedirect(code int) bool {
	switch code {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect, http.StatusPermanentRedirect:
		return true
	default:
		return false
	}
}

type RedirectError struct {
	HandlerError
}

func (e RedirectError) Location() string {
	location := e.Header.Get("Location")
	if len(location) == 0 {
		return "/"
	}
	return location
}

func (e RedirectError) Status() int {
	// 300 StatusMultipleChoices not supported
	code := e.HandlerError.Code
	if code > 300 && code < 400 {
		return code
	} else {
		return http.StatusTemporaryRedirect
	}
}

func (e RedirectError) Temporary() bool {
	switch e.Status() {
	case http.StatusMovedPermanently, http.StatusPermanentRedirect:
		return false
	default:
		return true
	}
}

func (e RedirectError) Error() string {
	return fmt.Sprintf("%v redirect: %q", e.Status(), e.Location())
}

func newRedirect(code int, location string, args ...interface{}) *RedirectError {
	if len(args) > 0 {
		location = fmt.Sprintf(location, args...)
	}

	return &RedirectError{
		HandlerError: HandlerError{
			Code:   code,
			Header: tools.NewHeader("Location", location),
		},
	}
}

// 301
func NewMovedPermanently(location string, args ...interface{}) *RedirectError {
	return newRedirect(http.StatusMovedPermanently, location, args...)
}

// 302
func NewFound(location string, args ...interface{}) *RedirectError {
	return newRedirect(http.StatusFound, location, args...)
}

// 303
func NewSeeOther(location string, args ...interface{}) *RedirectError {
	return newRedirect(http.StatusSeeOther, location, args...)
}

// 307
func NewTemporaryRedirect(location string, args ...interface{}) *RedirectError {
	return newRedirect(http.StatusTemporaryRedirect, location, args...)
}

// 308
func NewPermanentRedirect(location string, args ...interface{}) *RedirectError {
	return newRedirect(http.StatusPermanentRedirect, location, args...)
}

// Attempts to convert a given error into a RedirectError
func AsRedirect(err error) (*RedirectError, bool) {

	if err != nil {
		switch v := err.(type) {
		case *RedirectError:
			// Ours, Redirect already
			return v, true
		case *PanicError:
			// Ours, not Redirect
			goto fail
		case web.Error:
			// Friendly
			code := v.Status()

			if !CodeIsRedirect(code) {
				// Not a redirect
				return nil, false
			} else if e, ok := err.(*HandlerError); ok {
				// check http.Header
				if loc := e.Header.Get("Location"); loc != "" {
					// Redirect
					p := newRedirect(code, loc)
					return p, true
				}
			} else if e, ok := err.(interface {
				Header() http.Header
			}); ok {
				// Redirect with `Header() http.Header` interface
				if loc := e.Header().Get("Location"); loc != "" {
					//  Redirect
					p := newRedirect(code, loc)
					return p, true
				}
			} else if e, ok := err.(interface {
				Location() string
			}); ok {
				// Redirect with `Location() string` interface
				if loc := e.Location(); loc != "" {
					// Redirect
					p := newRedirect(code, loc)
					return p, true
				}
			}
		}
		// fall through
	}

fail:
	return nil, false
}
