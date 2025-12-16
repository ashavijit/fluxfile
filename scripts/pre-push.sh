#!/bin/bash

echo "Running pre-push checks..."

echo "Checking code formatting..."
gofmt_output=$(gofmt -l .)
if [ -n "$gofmt_output" ]; then
    echo "The following files are not formatted:"
    echo "$gofmt_output"
    echo "Please run 'go fmt ./...' to fix them."
    exit 1
fi

echo "Running Linter..."
if ! golangci-lint run; then
    echo "Linting failed. Please fix the errors above."
    exit 1
fi

echo "Running Tests..."
if ! go test ./... -race -v; then
    echo "Tests failed."
    exit 1
fi

echo "All checks passed!"
exit 0
