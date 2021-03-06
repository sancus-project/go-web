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
	var err error

	if req.Form != nil {
		// Form already parsed
		return nil
	}

	if size < MinimumFormSize {
		size = DefaultFormSize
	}

	t := strings.Split(req.Header.Get("Content-Type"), ";")[0]
	if t == "application/x-www-form-urlencoded" {
		err = req.ParseForm()
	} else if strings.HasPrefix(t, "multipart/form-data") {
		err = req.ParseMultipartForm(size)
	} else {
		err = errors.New("Invalid Content-Type %q", t)
	}

	return errors.BadRequest(err).AsError()
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

func FormValueBool(req *http.Request, key string) (bool, error, bool) {
	var v bool

	s, err, ok := FormValue(req, key)
	if ok && err == nil {
		v, err = strconv.ParseBool(s)
	}
	return v, err, ok
}
