.PHONY: all fmt build test

GO ?= go
GOFMT ?= gofmt
GOFMT_FLAGS = -w -l -s

all: fmt build

fmt:
	find . -name '*.go' | xargs -r $(GOFMT) $(GOFMT_FLAGS)
	$(GO) mod tidy || true

build:
	$(GO) get -v ./...

test:
	$(GO) test -v ./...
