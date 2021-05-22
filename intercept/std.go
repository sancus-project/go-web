package intercept

import (
	"net/http"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

func Resolve(h web.Handler, eh web.ErrorHandlerFunc) http.Handler {
	if h == nil {
		return nil
	} else if h2, ok := h.(http.Handler); ok {
		return h2
	} else {

		if eh == nil {
			eh = errors.HandleError
		}

		fn := func(w http.ResponseWriter, r *http.Request) {
			if err := h.TryServeHTTP(w, r); err != nil {
				eh(w, r, err)
			}
		}

		return http.HandlerFunc(fn)
	}
}
