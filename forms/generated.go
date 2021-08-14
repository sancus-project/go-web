package forms

//go:generate ./generated.sh

import (
	"net/http"
	"strconv"

	"go.sancus.dev/core/typeconv"
)

// Code generated by ./generated.sh DO NOT EDIT

func formValueFloat(req *http.Request, key string, bitsize int) (float64, error, bool) {
	var v float64

	s, err, ok := FormValue(req, key)
	if ok && err == nil {
		v, err = strconv.ParseFloat(s, bitsize)
	}

	return v, err, ok
}

func FormValueFloat32(req *http.Request, key string) (float32, error, bool) {
	v, err, ok := formValueFloat(req, key, 32)
	return float32(v), err, ok
}

func formValueInt(req *http.Request, key string, base int, bitsize int) (int64, error, bool) {
	var v int64

	s, err, ok := FormValue(req, key)
	if ok && err == nil {
		v, err = strconv.ParseInt(s, base, bitsize)
	}

	return v, err, ok
}

func FormValueInt(req *http.Request, key string, base int) (int, error, bool) {
	v, err, ok := formValueInt(req, key, base, typeconv.IntSize)
	return int(v), err, ok
}

func FormValueInt8(req *http.Request, key string, base int) (int8, error, bool) {
	v, err, ok := formValueInt(req, key, base, 8)
	return int8(v), err, ok
}

func FormValueInt16(req *http.Request, key string, base int) (int16, error, bool) {
	v, err, ok := formValueInt(req, key, base, 16)
	return int16(v), err, ok
}

func FormValueInt32(req *http.Request, key string, base int) (int32, error, bool) {
	v, err, ok := formValueInt(req, key, base, 32)
	return int32(v), err, ok
}

func formValueUint(req *http.Request, key string, base int, bitsize int) (uint64, error, bool) {
	var v uint64

	s, err, ok := FormValue(req, key)
	if ok && err == nil {
		v, err = strconv.ParseUint(s, base, bitsize)
	}

	return v, err, ok
}

func FormValueUint(req *http.Request, key string, base int) (uint, error, bool) {
	v, err, ok := formValueUint(req, key, base, typeconv.IntSize)
	return uint(v), err, ok
}

func FormValueUint8(req *http.Request, key string, base int) (uint8, error, bool) {
	v, err, ok := formValueUint(req, key, base, 8)
	return uint8(v), err, ok
}

func FormValueUint16(req *http.Request, key string, base int) (uint16, error, bool) {
	v, err, ok := formValueUint(req, key, base, 16)
	return uint16(v), err, ok
}

func FormValueUint32(req *http.Request, key string, base int) (uint32, error, bool) {
	v, err, ok := formValueUint(req, key, base, 32)
	return uint32(v), err, ok
}
