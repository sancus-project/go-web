package middleware

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

// Handler creates a middleware that tries to use the web.Handler
// interface, and pass any error to the provided ErrorHandler.
// If no handler is provided, our errors.HandleError is used.
func Handler(h web.ErrorHandlerFunc) web.MiddlewareHandlerFunc {
	if h == nil {
		h = errors.HandleError
	}

	return func(next http.Handler) http.Handler {

		if next == nil {
			fn := func(w http.ResponseWriter, r *http.Request) {
				h(w, r, errors.ErrNotFound)
			}
			return http.HandlerFunc(fn)

		} else if wh, ok := next.(web.Handler); ok {

			fn := func(w http.ResponseWriter, r *http.Request) {
				if err := wh.TryServeHTTP(w, r); err != nil {
					h(w, r, err)
				}
			}
			return http.HandlerFunc(fn)

		} else {

			return next
		}
	}
}
