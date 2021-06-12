package errors

import (
	"go.sancus.dev/core/context"
	"go.sancus.dev/web"
)

func ErrorContext(ctx context.Context) *web.Error {

	if out, ok := ctx.Value(ErrorCtxKey).(*web.Error); ok {
		return out
	}
	return nil
}

func WithErrorContext(ctx context.Context, out *web.Error) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if out != nil {
		ctx = context.WithValue(ctx, ErrorCtxKey, out)
	}

	return ctx
}

var (
	ErrorCtxKey = context.NewContextKey("WebErrorContext")
)
