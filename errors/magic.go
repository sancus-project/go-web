package errors

import (
	"net/http"
)

func NewError(code int, headers http.Header, body []byte) error {

	if code < 100 {
		// 500
		code = http.StatusInternalServerError
	} else if code < 300 {
		// OK
		return nil
	}

	switch code {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect, http.StatusPermanentRedirect:
		location := headers.Get("Location")
		return &RedirectError{
			location: location,
			code: code,
		}
	case http.StatusMethodNotAllowed:
		allowed := headers.Get("Allowed")
		return &MethodNotAllowedError{
			Allowed: strings.Split(allowed, ", "),
		}
	default:
		return &HandlerError{
			Code: code,
			Headers: headers,
		}
	}
}
