package errors

import (
	"net/http"
)

var (
	ErrNotFound = &HandlerError{Code: http.StatusNotFound}
)
