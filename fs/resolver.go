package fs

import (
	"net/http"

	"go.sancus.dev/web/context"
)

type ResolverFunc func(req *http.Request) string

func DefaultResolver(req *http.Request) string {
	return context.GetRoutePath(req)
}
