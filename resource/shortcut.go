// Code generated DO NOT EDIT
package resource

//go:generate ./shortcut.sh

import (
	"net/http"

	"go.sancus.dev/web/forms"
)

func (_ Resource) FormValue(req *http.Request, key string) (string, error) {
	return forms.FormValue(req, key)
}

func (_ Resource) FormValueInt16(req *http.Request, key string) (int16, error) {
	return forms.FormValueInt16(req, key, 10)
}

func (_ Resource) FormValueInt32(req *http.Request, key string) (int32, error) {
	return forms.FormValueInt32(req, key, 10)
}

func (_ Resource) FormValueFloat32(req *http.Request, key string) (float32, error) {
	return forms.FormValueFloat32(req, key)
}

func (_ Resource) FormValueBool(req *http.Request, key string) (bool, error) {
	return forms.FormValueBool(req, key)
}
