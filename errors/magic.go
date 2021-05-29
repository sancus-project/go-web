package errors

import (
	"net/http"
	"strings"

	"go.sancus.dev/web"
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
			code:     code,
		}
	case http.StatusMethodNotAllowed:
		allowed := headers.Get("Allowed")
		return &MethodNotAllowedError{
			Allowed: strings.Split(allowed, ", "),
		}
	default:
		return &HandlerError{
			Code:    code,
			Headers: headers,
		}
	}
}

func NewFromError(err error) error {
	var code int

	if err == nil {
		return nil
	} else if _, ok := err.(http.Handler); ok {
		return err
	} else if e, ok := err.(web.Error); ok {
		code = e.Status()
	} else {
		code = http.StatusInternalServerError
	}

	return &HandlerError{
		Code: code,
		Err:  err,
	}
}
