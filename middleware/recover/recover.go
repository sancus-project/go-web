package recover

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

func Recover(h web.ErrorHandlerFunc) web.MiddlewareHandlerFunc {
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
