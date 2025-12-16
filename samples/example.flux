# FluxFile - Modern Task Runner
# https://github.com/ashavijit/fluxfile
# Sample file for GitHub Linguist

var PROJECT = myapp
var VERSION = $(shell "git describe --tags")
var MODE = development

# Build task with caching
task build:
    desc: Build the application
    deps: fmt, lint
    cache: true
    inputs:
        src/**/*.go
        go.mod
        go.sum
    outputs:
        dist/${PROJECT}
    env:
        CGO_ENABLED = 0
        GOOS = linux
    run:
        go build -ldflags="-X main.version=${VERSION}" -o dist/${PROJECT} ./cmd

# Formatting
task fmt:
    desc: Format Go code
    run:
        go fmt ./...

# Linting
task lint:
    desc: Run linter
    run:
        golangci-lint run

# Testing
task test:
    desc: Run tests with coverage
    deps: build
    run:
        go test -v -cover ./...

# Development with watch mode
task dev:
    desc: Watch and rebuild on changes
    watch: **/*.go
    ignore:
        vendor/**
        **/*_test.go
        .git/**
    run:
        go run ./cmd

# Docker build
task docker-build:
    desc: Build Docker image
    docker: true
    run:
        docker build -t ${PROJECT}:${VERSION} .

# Remote deployment
task deploy:
    desc: Deploy to production
    if: MODE == production
    remote: deploy@prod.example.com
    deps: build, test
    run:
        docker-compose pull
        docker-compose up -d

# Matrix build for cross-compilation
task build-all:
    desc: Cross-compile for multiple platforms
    matrix:
        os: linux, darwin, windows
        arch: amd64, arm64
    run:
        GOOS=${os} GOARCH=${arch} go build -o dist/${PROJECT}-${os}-${arch}

# Parallel CI tasks
task ci:
    desc: Run all CI checks
    parallel: true
    deps: test, lint, build

# Clean up
task clean:
    desc: Remove build artifacts
    run:
        rm -rf dist/
        rm -rf .flux/

# Development profile
profile dev:
    env:
        MODE = development
        LOG_LEVEL = debug
        PORT = 3000

# Production profile
profile prod:
    env:
        MODE = production
        LOG_LEVEL = error
        PORT = 80

# Include external tasks
include "plugins/docker.flux"
