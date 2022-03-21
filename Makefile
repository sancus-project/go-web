.PHONY: all generate fmt build install get test deps

GO ?= go
GOPATH ?= $(CURDIR)
GOBIN ?= $(GOPATH)/bin
GOFMT ?= gofmt
GOFMT_FLAGS = -w -l -s
GOGENERATE_FLAGS = -v

all: generate fmt get build

fmt:
	@find . -name '*.go' | xargs -r $(GOFMT) $(GOFMT_FLAGS)
	$(GO) mod tidy

generate: deps
	@git grep -l '^//go:generate' | sed -n -e 's|\(.*\)/[^/]\+\.go$$|\1|p' | sort -u | while read d; do \
		git grep -l '^//go:generate' "$$d"/*.go | xargs -r $(GO) generate $(GOGENERATE_FLAGS); \
	done

get:
	$(GO) get -v ./...

install:
	$(GO) install -v ./cmd/...

build:
	$(GO) build -v ./...

test:
	$(GO) test -v ./...

deps: $(GOBIN)/peg

$(GOBIN)/peg:
	$(GO) install github.com/pointlander/peg@latest
