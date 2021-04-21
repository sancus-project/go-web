package errors

import (
	"fmt"
	"net/http"

	"go.sancus.dev/web"
)

func ErrorText(code int) string {
	text := http.StatusText(code)

	if len(text) == 0 {
		text = fmt.Sprintf("Unknown Error %d", code)
	} else if code >= 400 {
		text = fmt.Sprintf("%s (Error %d)", text, code)
	}

	return text
}

func HandleMiddlewareError(w http.ResponseWriter, r *http.Request, err error, next http.Handler) {
	if err != nil {
		if next != nil {
			// middleware
			if e, ok := err.(web.Error); ok && e.Status() == http.StatusNotFound {
				next.ServeHTTP(w, r)
				return
			}
		}

		// does the error know how to render itself?
		h, ok := err.(http.Handler)
		if !ok {
			var code int

			// but if it doesn't, wrap it in HandlerError{}
			if e, ok := err.(web.Error); ok {
				code = e.Status()
			} else {
				code = http.StatusInternalServerError
			}

			h = &errors.HandlerError{
				Code: code,
				Err:  err,
			}
		}

		h.ServeHTTP(w, r)
	}

}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	HandleMiddlewareError(w, r, err, nil)
}
