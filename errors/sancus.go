package errors

import (
	"go.sancus.dev/core/errors"
)

func New(s string, args ...interface{}) error {
	return errors.New(s, args...)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Here() *errors.Frame {
	return errors.StackFrame(1)
}
