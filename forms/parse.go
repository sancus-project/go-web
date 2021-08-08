package forms

import (
	"net/http"
	"strconv"
	"strings"

	"go.sancus.dev/web/errors"
)

const (
	KiB = 1 << 10
	MiB = 1 << 20

	MinimumFormSize = 4 * KiB
	DefaultFormSize = 1 * MiB
)

func ParseForm(req *http.Request, size int64) error {

	if req.Form != nil {
		// Form already parsed
		return nil
	}

	if size < MinimumFormSize {
		size = DefaultFormSize
	}

	t := strings.Split(req.Header.Get("Content-Type"), ";")[0]
	if t == "application/x-www-form-urlencoded" {
		return req.ParseForm()
	} else if strings.HasPrefix(t, "multipart/form-data") {
		return req.ParseMultipartForm(size)
	} else {
		var err errors.BadRequestError
		err.AppendErrorf("Invalid Content-Type %q", t)
		return err
	}

	return nil
}

func FormValue(req *http.Request, key string) (string, error, bool) {
	if req.Form == nil {
		if err := ParseForm(req, 0); err != nil {
			return "", err, false
		}
	}

	if v, ok := req.Form[key]; ok {
		return strings.TrimSpace(v[0]), nil, true
	} else {
		return "", nil, false
	}
}

func FormValueInt(req *http.Request, key string, base int, bitSize int) (int64, error, bool) {
	var n int64

	s, err, ok := FormValue(req, key)
	if ok && err == nil {
		n, err = strconv.ParseInt(s, base, bitSize)
	}
	return n, err, ok
}

func FormValueBool(req *http.Request, key string) (bool, error, bool) {
	var v bool

	s, err, ok := FormValue(req, key)
	if ok && err == nil {
		v, err = strconv.ParseBool(s)
	}
	return v, err, ok
}

func FormValueFloat(req *http.Request, key string, bitSize int) (float64, error, bool) {
	var v float64

	s, err, ok := FormValue(req, key)
	if ok && err == nil {
		v, err = strconv.ParseFloat(s, bitSize)
	}
	return v, err, ok
}
