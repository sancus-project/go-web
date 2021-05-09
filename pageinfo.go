package web

import (
	"net/http"
)

// Router can check for the existance of the requested resource
type RouterPageInfo interface {
	PageInfo(*http.Request) (interface{}, bool)
}
