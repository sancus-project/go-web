#!/bin/sh

set -eu

if [ -n "${GOFILE:-}" ]; then
	exec > "$GOFILE~"
	trap "rm -f '$GOFILE~'" EXIT
fi

cat <<EOT
package ${GOPACKAGE:-undefined}

//go:generate $0${*:+ $*}
EOT

gen() {
	local K="$1"
	local T="${K}Error"
	local S="Status${K}"

cat <<EOT

type $T struct {
	errors.ErrorStack
}

func (err *$T) AsError() error {
	if err.Ok() {
		return nil
	} else {
		return err
	}
}

func (err *$T) Status() int {
	if err.Ok() {
		return http.StatusOK
	} else {
		return http.$S
	}
}

func (err *$T) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveHTTP(err, w, r)
}

func (err *$T) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return tryServeHTTP(err, w, r)
}

func $K(errs... error) *$T {
	return &$T{
		ErrorStack: errors.NewErrorStack(errs...),
	}
}
EOT
}

if [ $# -gt 0 ]; then

	cat <<EOT

import (
	"net/http"

	"go.sancus.dev/core/errors"
)

var (
EOT
	for x; do
		cat <<EOT
	// Constant http.Status$x HandlerError
	Err$x = &HandlerError{Code: http.Status$x}
EOT
	done
	echo ")"

	for x; do
		gen "$x"
	done

	if [ -n "${GOFILE:-}" ]; then
		mv "$GOFILE~" "$GOFILE"
	fi
fi
