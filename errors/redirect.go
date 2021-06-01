package errors

import (
	"fmt"
	"net/http"
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
