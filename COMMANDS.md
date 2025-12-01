# Fluxfile Commands Reference

## Basic Commands

```bash
# List all tasks
flux -l

# Run a task
flux <task-name>
flux -t <task-name>

# Run with profile
flux -p <profile> <task-name>

# Watch mode
flux -w <task-name>

# Disable cache
flux --no-cache <task-name>

# Use custom FluxFile
flux -f <path> <task-name>

# Show version
flux -v
```

## FluxFile Syntax

### Variables
```
var PROJECT = myapp
var VERSION = $(shell "git describe --tags")
```

### Task Directives

```
task name:
    desc: Task description
    deps: dep1, dep2
    parallel: true
    if: MODE == prod
    profile_task: prod
    secrets:
        API_KEY
        DB_PASSWORD
    pre:
        file: config.yaml
        command: docker
        env: PATH
    retries: 3
    retry_delay: 2s
    timeout: 30s
    cache: true
    inputs:
        src/**/*.go
        go.mod
    outputs:
        dist/binary
    env:
        KEY = value
    run:
        echo "Hello"
        go build -o app
    watch: **/*.go
    ignore:
        node_modules/**
        .git/**
    matrix:
        os: linux, darwin, windows
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

profile prod:
    env:
        MODE = production
        LOG_LEVEL = error
```

### Includes

```
include "common.flux"
include "docker.flux"
```

## Operators

- `==` - Equal
- `!=` - Not equal
- `>` - Greater than
- `<` - Less than
- `>=` - Greater than or equal
- `<=` - Less than or equal

## Examples

```
# Conditional task
task deploy:
    if: MODE == prod
    run:
        kubectl apply -f k8s/

# Parallel deps
task ci:
    parallel: true
    deps: test, lint, build

# Cached build
task build:
    cache: true
    inputs:
        src/**/*.go
    outputs:
        dist/app
    run:
        go build -o dist/app

# Retry on failure
task flaky-test:
    retries: 3
    retry_delay: 5s
    run:
        npm test

# With timeout
task deploy:
    timeout: 5m
    run:
        ./deploy.sh
```
