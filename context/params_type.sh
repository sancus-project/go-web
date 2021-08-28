#!/bin/sh

set -eu

F="${0%.sh}.go"
trap "rm -f '$F~'" EXIT
exec > "$F~"

cat <<EOT
package context

//go:generate $0 $*

import (
	"go.sancus.dev/core/typeconv"
)
EOT

for x; do

	n=${x%:*}
	t=${x#*:}
	if [ "x$t" = "x$x" ]; then
		t=$(echo "$x" | tr 'A-Z' 'a-z')
	fi

	case "$t" in
	int*) extra=", 10" ;;
	*) extra=
	esac

	if [ -n "$n" ]; then
	cat <<EOT

// Get $t parameter from RouteContext
func (rctx *RoutingContext) Get$n(key string) (x $t, err error, ok bool) {
	v, _, ok := rctx.Get(key)
	if ok {
		x, err = typeconv.As$n(v)
	}
	return
}
EOT
	fi
	cat <<EOT

// Get slice of $t parameters from RouteContext
func (rctx *RoutingContext) Get${n}Slice(key string) (x []$t, err error, ok bool) {
	v, _, ok := rctx.Get(key)
	if ok {
		x, err = typeconv.As${n}Slice(v)
	}
	return
}
EOT
done

if ! diff -u "$F" "$F~" >&2; then
	mv "$F~" "$F"
fi
