package errors

import (
	"fmt"
	"net/http"

	"go.sancus.dev/web"
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
	location string
	code     int
}

func (e RedirectError) Location() string {
	if len(e.location) == 0 {
		return "/"
	}
	return e.location
}

func (e RedirectError) Status() int {
	// 300 StatusMultipleChoices not supported
	if e.code > 300 && e.code < 400 {
		return e.code
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

func (e *RedirectError) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return e
}

func (e *RedirectError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := e.Status()
	location := e.Location()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Location", location)
	w.WriteHeader(code)

	fmt.Fprintln(w, "Redirected to %s", location)
}

// 301
func NewMovedPermanently(location string) *RedirectError {
	return &RedirectError{location, http.StatusMovedPermanently}
}

// 302
func NewFound(location string) *RedirectError {
	return &RedirectError{location, http.StatusFound}
}

// 303
func NewSeeOther(location string) *RedirectError {
	return &RedirectError{location, http.StatusSeeOther}
}

// 307
func NewTemporaryRedirect(location string) *RedirectError {
	return &RedirectError{location, http.StatusTemporaryRedirect}
}

// 308
func NewPermanentRedirect(location string) *RedirectError {
	return &RedirectError{location, http.StatusPermanentRedirect}
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
			// Friedly
			code := v.Status()

			if !CodeIsRedirect(code) {
				// Not a redirect
				return nil, false
			} else if e, ok := err.(*HandlerError); ok {
				// check http.Header
				if loc := e.Header.Get("Location"); loc != "" {
					// Redirect
					p := &RedirectError{
						location: loc,
						code:     code,
					}
					return p, true
				}
			} else if e, ok := err.(interface {
				Header() http.Header
			}); ok {
				// Redirect with `Header() http.Header` interface
				if loc := e.Header().Get("Location"); loc != "" {
					//  Redirect
					p := &RedirectError{
						location: loc,
						code:     code,
					}
					return p, true
				}
			} else if e, ok := err.(interface {
				Location() string
			}); ok {
				// Redirect with `Location() string` interface
				if loc := e.Location(); loc != "" {
					// Redirect
					p := &RedirectError{
						location: loc,
						code:     code,
					}
					return p, true
				}
			}
		}
		// fall through
	}

fail:
	return nil, false
}
