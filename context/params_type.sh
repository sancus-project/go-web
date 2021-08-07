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
func (rctx *RoutingContext) Get$n(key string) ($t, bool) {
	var zero $t

	if v, ok := rctx.Get(key); ok {
		return typeconv.As$n(v)
	}

	return zero, false
}
EOT
	fi
	cat <<EOT

// Get slice of $t parameters from RouteContext
func (rctx *RoutingContext) Get${n}Slice(key string) ([]$t, bool) {
	if v, ok := rctx.Get(key); ok {
		return typeconv.As${n}Slice(v)
	}

	return nil, false
}
EOT
done

if ! cmp -s "$F" "$F~"; then
	mv "$F~" "$F"
fi
