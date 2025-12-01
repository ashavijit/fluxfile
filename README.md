# FluxFile

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
