# FluxFile

[![CI](https://github.com/ashavijit/fluxfile/actions/workflows/ci.yaml/badge.svg)](https://github.com/ashavijit/fluxfile/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ashavijit/fluxfile)](https://goreportcard.com/report/github.com/ashavijit/fluxfile)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ashavijit/fluxfile)](https://go.dev/)

Modern task runner and build automation tool with a clean, minimal syntax.

## Features

- **Task Descriptions** - Add descriptions to tasks for better documentation
- Clean, indentation-based DSL
- Dependency graph resolution with cycle detection
- Task result caching based on file hashes
- **Enhanced Caching** - Input/output tracking for incremental builds
- File watching for automatic re-execution with ignore patterns
- **Conditional Execution** - Run tasks based on environment conditions
- **Parallel Task Execution** - Run dependencies concurrently
- Matrix builds for multi-platform compilation
- Docker container execution
- Remote execution over SSH
- Variable expansion and shell command execution
- Profile support for environment-specific configuration
- Include directive for modular FluxFiles
- **Project Scaffolding** - `flux init` with templates for Go, Node, Python, Rust
- **Execution Reports** - `--report` flag for timing metrics
- **HTML Log Viewer** - `flux logs` opens execution history in browser

## Installation

### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/ashavijit/fluxfile/main/scripts/install.sh | sh
```

### Windows

```powershell
iwr -useb https://raw.githubusercontent.com/ashavijit/fluxfile/main/scripts/install.ps1 | iex
```

### From Source

```bash
git clone https://github.com/ashavijit/fluxfile
cd fluxfile
make install
```

## Quick Start

Create a `FluxFile` in your project root:

```
task build:
    run:
        go build -o app ./cmd

task test:
    deps: build
    run:
        go test ./...

task dev:
    watch: **/*.go
    run:
        go run ./cmd
```

Run tasks:

```bash
flux build
flux -t test
flux -w dev
```

## Syntax Reference

### Variables

```
var PROJECT = myapp
var VERSION = $(shell "git describe --tags")
```

### Tasks

```
task name:
    desc: Task description shown in help
    deps: dep1, dep2
    parallel: true
    if: MODE == prod
    env:
        KEY = value
    run:
        command1
        command2
    watch: **/*.go
    ignore:
        node_modules/**
        .git/**
    cache: true
    inputs:
        src/**/*.go
    outputs:
        dist/binary
    matrix:
        os: linux, darwin
        arch: amd64, arm64
    docker: true
    remote: user@host
```

### Profiles

```
profile dev:
    env:
        MODE = development
        DEBUG = true
```

### Include

```
include "tasks/docker.flux"
```

## CLI Usage

```
flux [options] <task>

Options:
  -t string     Task to execute
  -p string     Profile to apply
  -l            List all tasks
  -w            Watch mode
  -no-cache     Disable caching
  -f string     Path to FluxFile
  -v            Show version
  --init        Initialize new FluxFile
  --template    Template for init (go, node, python, rust, generic)
  --report      Show execution timing report
  --report-json Save report as JSON file

Commands:
  flux init              Create FluxFile from detected project type
  flux logs              Open execution logs in browser
```

## Examples

### Matrix Build

```
task cross-compile:
    matrix:
        os: linux, darwin, windows
        arch: amd64, arm64
    run:
        GOOS=${os} GOARCH=${arch} go build -o dist/app-${os}-${arch}
```

### Docker Build

```
task docker-test:
    docker: true
    run:
        npm install
        npm test
```

### Remote Deployment

```
task deploy:
    remote: deploy@prod.example.com
    run:
        docker-compose pull
        docker-compose up -d
```

## License

MIT
