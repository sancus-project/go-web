package middleware

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

// Recoverer creates a middleware to catch panic() and pass it to
// an ErrorHandler wrapped as *errors.PanicError.
// If no handler is provided, our errors.HandleError is used.
func Recoverer(h web.ErrorHandlerFunc) web.MiddlewareHandlerFunc {
	if h == nil {
		h = errors.HandleError
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := errors.Recover(); err != nil {
					h(w, r, err)
				}
			}()
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
