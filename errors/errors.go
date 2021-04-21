package errors

import (
	"fmt"
	"log"
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

func HandleError(w http.ResponseWriter, r *http.Request, err error, next http.Handler) {
	if err == nil {
		// served
	} else if h, ok := err.(web.Error); ok && h.Status() == http.StatusNotFound {

		if next != nil {
			next.ServeHTTP(w, r)
		} else if h, ok := err.(http.Handler); ok {
			h.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}

	} else if h, ok := err.(http.Handler); ok {
		h.ServeHTTP(w, r)
	} else {
		log.Fatal(err)
	}
}
