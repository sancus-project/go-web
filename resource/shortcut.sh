#!/bin/sh

set -eu

F="${0%.sh}.go"
trap "rm -f '$F~'" EXIT
exec > "$F~"

cat <<EOT
// Code generated DO NOT EDIT
package resource

//go:generate $0

import (
	"net/http"

	"go.sancus.dev/web/forms"
)
EOT

for x in :string Int16 Int32 Float32 Bool; do

	n=${x%:*}
	t=${x#*:}
	if [ "x$t" = "x$x" ]; then
		t=$(echo "$x" | tr 'A-Z' 'a-z')
	fi

	case "$t" in
	int*) extra=", 10" ;;
	*) extra=
	esac

	cat <<EOT

func (_ Resource) FormValue$n(req *http.Request, key string) ($t, error) {
	return forms.FormValue$n(req, key$extra)
}
EOT
done

if ! cmp -s "$F" "$F~"; then
	mv "$F~" "$F"
fi