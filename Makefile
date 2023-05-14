#!/usr/bin/make

.PHONY: build
build:
	go build -a -o ./gitformer

.PHONY: test
test:
	go test -v ./cmd/
	go test -v ./pkg/playbook/

.PHONY: lint
lint:
	go vet -v ./cmd/
	go vet -v ./pkg/playbook/