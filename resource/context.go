package resource

import (
	"go.sancus.dev/web/context"
	"go.sancus.dev/web/errors"
)

// check request is a exact match (leaf, not intermediate router)
// and makes sure there is no trailing slash
func DefaultResourceChecker(ctx context.Context) (context.Context, error) {
	if rctx := context.RouteContext(ctx); rctx != nil {
		switch rctx.RoutePath {
		case "":
			break
		case "/":
			return nil, errors.NewSeeOther(rctx.RoutePrefix)
		default:
			return nil, errors.ErrNotFound
		}
	}
	return ctx, nil
}
