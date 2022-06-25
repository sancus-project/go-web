package errors

import (
	"io/ioutil"
	"net/http"
	"strings"

	"go.sancus.dev/web"
)

func NewErrorFromResponse(res *http.Response) error {
	if res == nil {
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	return newWebError(res.StatusCode, res.Header, body, err)
}

func NewWebError(code int, headers http.Header, body []byte) web.Error {
	return newWebError(code, headers, body, nil)
}

func newWebError(code int, headers http.Header, body []byte, readError error) web.Error {
	if code < 100 {
		// 500
		code = http.StatusInternalServerError
	} else if code < 300 {
		// OK
		return nil
	}

	if CodeIsRedirect(code) {
		location := headers.Get("Location")
		return newRedirect(code, location)
	} else if code == http.StatusMethodNotAllowed {
		allowed := headers.Get("Allowed")
		return &MethodNotAllowedError{
			Allowed: strings.Split(allowed, ", "),
		}
	} else {
		return &HandlerError{
			Code:   code,
			Err:    readError,
			Header: headers,
		}
	}
}

func AsWebError(err error) web.Error {
	var p web.Error
	var ok bool

	if err == nil {
		// Ignore
	} else if p, ok = err.(web.Error); ok {
		// Ready
	} else {
		// Wrap
		p = &HandlerError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return p
}

// Returns http.Handler capable web.Error
func NewFromError(err error) error {

	if err == nil {
		// Ignore
		return nil
	} else if p, ok := AsRedirect(err); ok {
		// Redirect
		return p
	}

	var code int

	switch v := err.(type) {
	case http.Handler:
		// if it can render itself it might know better
		goto done
	case web.Error:
		// Friedly
		code = v.Status()
	default:
		code = http.StatusInternalServerError
	}

	// Wrap
	err = &HandlerError{
		Code: code,
		Err:  err,
	}
done:
	return err
}
