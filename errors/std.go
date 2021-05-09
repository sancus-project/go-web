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
