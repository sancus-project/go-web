package middleware

import (
	"net/http"
)

// Do-Nothing middleware
func NOP(next http.Handler) http.Handler {
	return next
}
