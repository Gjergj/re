SHELL=/bin/bash -eo pipefail

PWD = $(shell pwd)
GO ?= go

dev_run: generate-mocks
	docker compose build
	docker compose up

config:
	go install github.com/vektra/mockery/v2@latest


generate-mocks:
	 mockery --all --dir=pkg