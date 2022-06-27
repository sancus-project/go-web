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

func CopyHeaders(dst http.Header, src http.Header, except ...string) {
	for key, values := range src {
		for _, k := range except {
			if strings.EqualFold(key, k) {
				// skip
				continue
			}
		}
		for _, value := range values {
			// TODO: deduplicate
			dst.Add(key, value)
		}
	}
}
