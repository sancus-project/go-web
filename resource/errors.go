package resource

import (
	"fmt"
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

func (_ Resource) SeeOther(location string, args ...interface{}) *errors.RedirectError {
	if len(args) > 0 {
		location = fmt.Sprintf(location, args...)
	}
	return errors.NewSeeOther(path.Clean(location))
}
