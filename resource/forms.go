package resource

import (
	"net/http"

	"go.sancus.dev/web/forms"
)

func (_ Resource) ParseForm(req *http.Request, size int64) error {
	if req.Form == nil {
		return forms.ParseForm(req, size)
	}
	return nil
}
