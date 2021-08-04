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

func FormValue(req *http.Request, key string) (string, error) {
	if req.Form == nil {
		if err := ParseForm(req, 0); err != nil {
			return "", err
		}
	}

	if v, ok := req.Form[key]; ok {
		return strings.TrimSpace(v[0]), nil
	} else {
		return "", nil
	}
}

func FormValueInt(req *http.Request, key string, base int, bitSize int) (int64, error) {
	if s, err := FormValue(req, key); err != nil {
		return 0, err
	} else {
		return strconv.ParseInt(s, base, bitSize)
	}
}

func FormValueBool(req *http.Request, key string) (bool, error) {
	if s, err := FormValue(req, key); err != nil {
		return false, err
	} else {
		return strconv.ParseBool(s)
	}
}

func FormValueFloat(req *http.Request, key string, bitSize int) (float64, error) {
	if s, err := FormValue(req, key); err != nil {
		return 0., err
	} else {
		return strconv.ParseFloat(s, bitSize)
	}
}
