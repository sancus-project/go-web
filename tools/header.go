package tools

import (
	"fmt"
	"net/http"
)

func NewHeader(key, value string, args ...interface{}) http.Header {
	if len(args) > 0 {
		value = fmt.Sprintf(value, args...)
	}

	m := make(map[string][]string, 1)

	if len(key) > 0 {
		m[key] = []string{value}
	}

	return m
}
