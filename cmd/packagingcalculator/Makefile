SHELL=/bin/bash -eo pipefail

PWD = $(shell pwd)
GO ?= go

build:
	mkdir -p bin
	$(GO) build -o bin .

test:
	$(GO) test -v ./... -race