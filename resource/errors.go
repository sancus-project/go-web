package resource

import (
	"path"

	"go.sancus.dev/web/errors"
)

func (_ Resource) BadRequest(err error) *errors.BadRequestError {
	var valid errors.BadRequestError
	if err != nil {
		valid.AppendError(err)
	}
	return &valid
}

func (_ Resource) BadRequestf(s string, args ...interface{}) *errors.BadRequestError {
	var valid errors.BadRequestError
	if len(s) > 0 {
		valid.AppendErrorf(s, args...)
	}
	return &valid
}

func (_ Resource) SeeOther(location string) *errors.RedirectError {
	return errors.NewSeeOther(path.Clean(location))
}
