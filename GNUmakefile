default: fmt lint install generate

# Provider name and version
PROVIDER_NAME=kaizen
VERSION=0.0.1
BINARY_NAME=terraform-provider-$(PROVIDER_NAME)_v$(VERSION)

# Determine OS-specific binary extension
ifeq ($(OS),Windows_NT)
	BINARY_EXT=.exe
else
	BINARY_EXT=
endif

build:
	go build -v ./...

install: build
	go build -o $(GOPATH)/bin/$(BINARY_NAME)$(BINARY_EXT) .

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate
