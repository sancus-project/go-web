package resource

import (
	"context"
	"net/http"
)

type ContextChecker func(ctx context.Context) (context.Context, error)

type Checker interface {
	Check(ctx context.Context) (context.Context, error)
}

// GET
type Getter interface {
	Get(rw http.ResponseWriter, req *http.Request) error
}

// HEAD
type Peeker interface {
	Head(rw http.ResponseWriter, req *http.Request) error
}

// POST
type Poster interface {
	Post(rw http.ResponseWriter, req *http.Request) error
}

// OPTIONS
type Optioner interface {
	Options(rw http.ResponseWriter, req *http.Request) error
}
