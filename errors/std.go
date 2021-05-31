package errors

import (
	"errors"
	"fmt"
)

func New(s string, args ...interface{}) error {
	if len(args) > 0 {
		return fmt.Errorf(s, args...)
	}

	return errors.New(s)
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
