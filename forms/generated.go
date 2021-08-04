// Code generated DO NOT EDIT
package forms

//go:generate ./generated.sh

import (
	"net/http"
)

func FormValueFloat32(req *http.Request, key string) (float32, error) {
	var zero float32
	if v, err := FormValueFloat(req, key, 32); err != nil {
		return zero, err
	} else {
		return float32(v), nil
	}
}

func FormValueInt16(req *http.Request, key string, base int) (int16, error) {
	var zero int16
	if v, err := FormValueInt(req, key, base, 16); err != nil {
		return zero, err
	} else {
		return int16(v), nil
	}
}

func FormValueInt32(req *http.Request, key string, base int) (int32, error) {
	var zero int32
	if v, err := FormValueInt(req, key, base, 32); err != nil {
		return zero, err
	} else {
		return int32(v), nil
	}
}
