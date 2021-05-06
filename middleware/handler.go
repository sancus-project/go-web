package middleware

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

// Handler creates a middleware that tries to use the web.Handler
// interface, and pass any error to the provided ErrorHandler.
// This middleware will also catch panic(), superseeding Recoverer().
// If no handler is provided, our errors.HandleError is used.
func Handler(h web.ErrorHandlerFunc) web.MiddlewareHandlerFunc {
	if h == nil {
		h = errors.HandleError
	}

	return func(next http.Handler) http.Handler {
		var fn http.HandlerFunc

		if next == nil {
			fn = func(w http.ResponseWriter, r *http.Request) {
				h(w, r, errors.ErrNotFound)
			}

		} else if wh, ok := next.(web.Handler); ok {

			fn = func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if err := errors.Recover(); err != nil {
						h(w, r, err)
					}
				}()

				if err := wh.TryServeHTTP(w, r); err != nil {
					h(w, r, err)
				}
			}

		} else {

			fn = func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if err := errors.Recover(); err != nil {
						h(w, r, err)
					}
				}()

				next.ServeHTTP(w, r)
			}
		}

		return fn
	}
}
