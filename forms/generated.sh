#!/bin/sh

set -eu

F="${0%.sh}.go"
trap "rm -f '$F~'" EXIT
exec > "$F~"

cat <<EOT
// Code generated DO NOT EDIT
package forms

//go:generate $0

import (
	"net/http"
)
EOT

generate() {
	local N="$1" S="$2"
	local n="$N$S"
	local t="$(echo "$n" | tr A-Z a-z)"
	local extra_args= extra=
	shift 2

	if [ -n "${1:-}" ]; then
		extra_args=", $1"
		extra="$(echo "$1" | tr ',' '\n' | sed -e 's|^ \+||' -e 's| \+$||' | cut -d' ' -f1 | sed -e 's|^|, |' | tr -d '\n')"
	fi

cat <<EOT

func FormValue$n(req *http.Request, key string$extra_args) ($t, error) {
	var zero $t
	if v, err := FormValue$N(req, key$extra, $S); err != nil {
		return zero, err
	} else {
		return $t(v), nil
	}
}
EOT
}

generate Float 32
generate Int 16 "base int"
generate Int 32 "base int"

if ! cmp -s "$F" "$F~"; then
	mv "$F~" "$F"
fi
