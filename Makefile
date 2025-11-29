.PHONY: build test clean install

VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o dist/flux ./cmd/flux

build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/flux-linux-amd64 ./cmd/flux
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/flux-linux-arm64 ./cmd/flux
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/flux-darwin-amd64 ./cmd/flux
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/flux-darwin-arm64 ./cmd/flux
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/flux-windows-amd64.exe ./cmd/flux

test:
	go test ./... -v -race -coverprofile=coverage.out

clean:
	rm -rf dist/
	rm -f coverage.out

install: build
	cp dist/flux /usr/local/bin/flux

fmt:
	go fmt ./...

lint:
	golangci-lint run

mod:
	go mod tidy
	go mod download
