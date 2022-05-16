package router

import (
	"strings"

	"go.sancus.dev/web"
	"go.sancus.dev/web/context"
	"go.sancus.dev/web/errors"
)

// Resolve finds the best handler for a path and returns the corresponding RouteContext
func (m *Mux) Resolve(path string, rctx *context.RoutingContext) (web.Handler, *context.RoutingContext, bool) {
	s0, s1, h := m.findBestNode(path)

	if h == nil {
		return nil, nil, false
	}

	if s0 != "/" && s1 == "" && strings.HasSuffix(h.Pattern, "/*") {
		// redirect to the root of the subrouter
		return errors.NewPermanentRedirect("%s/", rctx.Path()), rctx, true
	}

	if rctx != nil {
		rctx = rctx.Step(s0)
	} else {
		rctx = context.NewRouteContext(s0, s1)
	}

	return h, rctx, true
}
