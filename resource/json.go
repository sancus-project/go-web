package resource

import (
	"encoding/json"
	"net/http"
)

func (_ Resource) WriteJSON(w http.ResponseWriter, prefix string, indent string, d interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	e := json.NewEncoder(w)
	e.SetIndent(prefix, indent)
	return e.Encode(d)
}
