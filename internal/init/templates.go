package init

import "fmt"

func GetTemplate(templateType string, projectName string) string {
	switch templateType {
	case "go":
		return goTemplate(projectName)
	case "node":
		return nodeTemplate(projectName)
	case "python":
		return pythonTemplate(projectName)
	case "rust":
		return rustTemplate(projectName)
	default:
		return genericTemplate(projectName)
	}
}

func goTemplate(name string) string {
	return fmt.Sprintf(`var PROJECT = %s

task build:
    desc: Build the Go binary
    cache: true
    inputs:
        **/*.go
        go.mod
        go.sum
    outputs:
        bin/${PROJECT}
    run:
        go build -o bin/${PROJECT} ./cmd/${PROJECT}

task test:
    desc: Run all tests
    run:
        go test ./... -v

task lint:
    desc: Run linter
    run:
        golangci-lint run

task fmt:
    desc: Format code
    run:
        go fmt ./...

task dev:
    desc: Watch and rebuild
    watch: **/*.go
    ignore:
        vendor/**
        **/*_test.go
    run:
        go run ./cmd/${PROJECT}

task clean:
    desc: Remove build artifacts
    run:
        rm -rf bin/

profile dev:
    env:
        GO_ENV = development

profile prod:
    env:
        GO_ENV = production
`, name)
}

func nodeTemplate(name string) string {
	return fmt.Sprintf(`var PROJECT = %s

task install:
    desc: Install dependencies
    run:
        npm install

task build:
    desc: Build the project
    deps: install
    cache: true
    inputs:
        src/**/*
        package.json
    outputs:
        dist/**/*
    run:
        npm run build

task test:
    desc: Run tests
    deps: install
    run:
        npm test

task lint:
    desc: Run linter
    run:
        npm run lint

task dev:
    desc: Start development server
    watch: src/**/*
    run:
        npm run dev

task clean:
    desc: Remove build artifacts
    run:
        rm -rf dist/ node_modules/

profile dev:
    env:
        NODE_ENV = development

profile prod:
    env:
        NODE_ENV = production
`, name)
}

func pythonTemplate(name string) string {
	return fmt.Sprintf(`var PROJECT = %s

task install:
    desc: Install dependencies
    run:
        pip install -r requirements.txt

task test:
    desc: Run tests
    run:
        pytest

task lint:
    desc: Run linter
    run:
        ruff check .

task fmt:
    desc: Format code
    run:
        black .

task dev:
    desc: Run development server
    watch: **/*.py
    run:
        python main.py

task clean:
    desc: Remove cache files
    run:
        find . -type d -name __pycache__ -exec rm -rf {} +
        rm -rf .pytest_cache/

profile dev:
    env:
        PYTHON_ENV = development

profile prod:
    env:
        PYTHON_ENV = production
`, name)
}

func rustTemplate(name string) string {
	return fmt.Sprintf(`var PROJECT = %s

task build:
    desc: Build the project
    cache: true
    inputs:
        src/**/*.rs
        Cargo.toml
    outputs:
        target/release/${PROJECT}
    run:
        cargo build --release

task test:
    desc: Run tests
    run:
        cargo test

task lint:
    desc: Run clippy
    run:
        cargo clippy

task fmt:
    desc: Format code
    run:
        cargo fmt

task dev:
    desc: Watch and rebuild
    watch: src/**/*.rs
    run:
        cargo run

task clean:
    desc: Remove build artifacts
    run:
        cargo clean

profile dev:
    env:
        RUST_ENV = development

profile prod:
    env:
        RUST_ENV = production
`, name)
}

func genericTemplate(name string) string {
	return fmt.Sprintf(`var PROJECT = %s

task build:
    desc: Build the project
    run:
        echo "Add build commands here"

task test:
    desc: Run tests
    run:
        echo "Add test commands here"

task clean:
    desc: Clean build artifacts
    run:
        echo "Add clean commands here"

profile dev:
    env:
        ENV = development

profile prod:
    env:
        ENV = production
`, name)
}
