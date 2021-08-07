package resource

import (
	"net/http"

	"go.sancus.dev/web/context"
)

func (_ Resource) RoutePath(req *http.Request) string {
	if rctx := context.RouteContext(req.Context()); rctx != nil {
		return rctx.Path()
	}
	return ""
}
