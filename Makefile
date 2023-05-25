#!/usr/bin/make

CURRENT_GIT_HASH := $(shell git rev-parse --verify HEAD)
DSN := ${SENTRY_DSN}

.PHONY: build
build:
	go build -a -o ./dist/gitformer \
	-ldflags="-X github.com/peachpielabs/gitformer/cmd/gitformer.version=$(CURRENT_GIT_HASH) \
			  -X main.version=$(CURRENT_GIT_HASH) \
			  -X main.buildTime=$(shell date +%Y-%m-%dT%H:%M:%S%z) \
			  -X main.dsn=$(DSN) \
			  -X main.environment=development"

.PHONY: test
test:
	go test -v ./cmd/gitformer/
	go test -v ./pkg/playbook/

.PHON: lint
lint:
	go vet -v ./cmd/*/
	go vet -v ./pkg/playbook/