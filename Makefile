.PHONY: build test lint fmt check clean

BINARY_NAME=pulumi-resource-flexera
PROVIDER_PATH=./cmd/pulumi-resource-flexera

build:
	go build -o $(BINARY_NAME) $(PROVIDER_PATH)

test:
	go test -v -cover ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...
	gofmt -s -w .

check: fmt lint test

clean:
	rm -f $(BINARY_NAME)
	go clean

install: build
	cp $(BINARY_NAME) $${GOPATH}/bin/

run-example:
	cd examples && pulumi up

.DEFAULT_GOAL := build