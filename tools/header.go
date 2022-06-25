package tools

import (
	"fmt"
	"net/http"
	"strings"
)

func NewHeader(key, value string, args ...interface{}) http.Header {
	return SetHeader(nil, key, value, args...)
}

func SetHeader(hdr http.Header, key, value string, args ...interface{}) http.Header {
	if hdr == nil {
		hdr = make(map[string][]string, 1)
	}

	if len(key) > 0 {
		if len(args) > 0 {
			value = fmt.Sprintf(value, args...)
		}

		value = strings.TrimSpace(value)
		if len(value) > 0 {
			hdr.Set(key, value)
		} else {
			hdr.Del(key)
		}
	}

	return hdr
}
