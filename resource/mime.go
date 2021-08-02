package resource

import (
	"net/http"

	"go.sancus.dev/web/errors"
	"go.sancus.dev/web/mimeparse"
)

func (_ Resource) BestMime(req *http.Request, supported ...string) (string, error) {
	if s := mimeparse.BestMatch(supported, req.Header.Get("Accept")); len(s) > 0 {
		return s, nil
	}

	return "", errors.ErrNotAcceptable
}
