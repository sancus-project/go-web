#!/bin/sh

set -eu

F="${0%.sh}.go"
trap "rm -f '$F~'" EXIT
exec > "$F~"

cat <<EOT
package errors

//go:generate $0

import (
	"go.sancus.dev/core/errors"
)

// Code generated by $0 DO NOT EDIT
EOT

for x in \
	InvalidValue \
	InvalidArgument \
	MissingField \
	NotImplemented \
	; do

	cat <<EOT

func Err$x(s string, args ...interface{}) error {
	return errors.Err$x(s, args...)
}

func As${x}Error(err error, s string, args ...interface{}) error {
	return errors.As${x}Error(err, s, args...)
}
EOT
done

if ! diff -u "$F" "$F~" >&2; then
	mv "$F~" "$F"
fi
