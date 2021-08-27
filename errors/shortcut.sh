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
	NotImplemented \
	; do

	cat <<EOT

func Err$x(s string, args ...interface{}) error {
	return errors.Err$x(s, args...)
}
EOT
done

if ! cmp -s "$F" "$F~"; then
	mv "$F~" "$F"
fi