package errors

import (
	"go.sancus.dev/core/errors"
)

type (
	Frame = errors.Frame
	Stack = errors.Stack
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

// errors.Frame
func Here() *errors.Frame {
	return errors.StackFrame(1)
}

func StackFrame(skip int) *errors.Frame {
	return errors.StackFrame(skip + 1)
}

func BackTrace(skip int) errors.Stack {
	return errors.BackTrace(skip + 1)
}

func StackTrace(err error) errors.Stack {
	return errors.StackTrace(err)
}
