.PHONY: all generate fmt build test deps

GO ?= go
GOPATH ?= $(CURDIR)
GOBIN ?= $(GOPATH)/bin
GOFMT ?= gofmt
GOFMT_FLAGS = -w -l -s
GOGENERATE_FLAGS = -v

all: generate fmt build

fmt:
	@find . -name '*.go' | xargs -r $(GOFMT) $(GOFMT_FLAGS)
	$(GO) mod tidy || true

generate: deps
	@git grep -l '^//go:generate' | sed -n -e 's|\(.*\)/[^/]\+\.go$$|\1|p' | sort -u | while read d; do \
		git grep -l '^//go:generate' "$$d"/*.go | xargs -r $(GO) generate $(GOGENERATE_FLAGS); \
	done

build:
	$(GO) get -v ./...

test:
	$(GO) test -v ./...

deps: $(GOBIN)/peg

$(GOBIN)/peg:
	$(GO) install github.com/pointlander/peg@latest
