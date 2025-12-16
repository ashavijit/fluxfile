<p align="center">
  <h1 align="center">⚡ FluxFile</h1>
  <p align="center">
    <strong>Modern task runner and build automation tool with a clean, minimal syntax.</strong>
  </p>
</p>

<p align="center">
  <a href="https://github.com/ashavijit/fluxfile/actions/workflows/ci.yaml"><img src="https://github.com/ashavijit/fluxfile/actions/workflows/ci.yaml/badge.svg" alt="CI"></a>
  <a href="https://goreportcard.com/report/github.com/ashavijit/fluxfile"><img src="https://goreportcard.com/badge/github.com/ashavijit/fluxfile" alt="Go Report Card"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/github/go-mod/go-version/ashavijit/fluxfile" alt="Go Version"></a>
</p>

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🚀 **Task Runner** | Clean, indentation-based DSL for defining tasks |
| 🔗 **Dependencies** | Automatic dependency resolution with cycle detection |
| ⚡ **Parallel Execution** | Run tasks concurrently for faster builds |
| 💾 **Smart Caching** | Input/output tracking for incremental builds |
| 👀 **File Watching** | Auto-rerun tasks when files change |
| 🐳 **Docker Support** | Run tasks inside containers |
| 🌐 **Remote Execution** | Deploy and run tasks over SSH |
| 📊 **Matrix Builds** | Cross-compile for multiple platforms |
| 📝 **Profiles** | Environment-specific configurations |
| 📈 **Execution Reports** | Timing metrics and performance insights |

---

## 🚀 Quick Install

```bash
# Linux / macOS
curl -fsSL https://raw.githubusercontent.com/ashavijit/fluxfile/main/scripts/install.sh | sh

# Windows (PowerShell)
iwr -useb https://raw.githubusercontent.com/ashavijit/fluxfile/main/scripts/install.ps1 | iex

# From Source
git clone https://github.com/ashavijit/fluxfile && cd fluxfile && make install
```

---

## 📖 Getting Started

### Basic Example

Create a `FluxFile` in your project root:

```yaml
var PROJECT = myapp

task build:
    desc: Build the binary
    run:
        go build -o bin/${PROJECT} ./cmd

task test:
    desc: Run tests
    deps: build
    run:
        go test ./... -v

task dev:
    desc: Watch and rebuild
    watch: **/*.go
    run:
        go run ./cmd
```

Run tasks:

```bash
flux build          # Run build task
flux -t test        # Run test task
flux -w dev         # Watch mode
flux -l             # List all tasks
flux --report test  # Show timing report
```

---

### Realistic Workflow: Full-Stack JS + Python

A monorepo with a React frontend, Python API, and shared tooling:

```yaml
var ENV = development

# Frontend (React/Node.js)
task frontend:install:
    desc: Install frontend dependencies
    inputs:
        frontend/package.json
        frontend/package-lock.json
    outputs:
        frontend/node_modules
    cache: true
    run:
        cd frontend && npm ci

task frontend:build:
    desc: Build React app
    deps: frontend:install
    inputs:
        frontend/src/**/*
        frontend/public/**/*
    outputs:
        frontend/dist
    cache: true
    run:
        cd frontend && npm run build

task frontend:dev:
    desc: Start frontend dev server
    deps: frontend:install
    watch: frontend/src/**/*
    run:
        cd frontend && npm run dev

# Backend (Python/FastAPI)
task backend:install:
    desc: Install Python dependencies
    inputs:
        backend/requirements.txt
    outputs:
        backend/.venv
    cache: true
    run:
        cd backend && python -m venv .venv
        cd backend && .venv/bin/pip install -r requirements.txt

task backend:dev:
    desc: Start Python API server
    deps: backend:install
    watch: backend/**/*.py
    run:
        cd backend && .venv/bin/uvicorn main:app --reload

task backend:test:
    desc: Run Python tests
    deps: backend:install
    run:
        cd backend && .venv/bin/pytest -v

# Full Stack
task dev:
    desc: Run full stack in parallel
    parallel: true
    deps: frontend:dev, backend:dev

task test:
    desc: Run all tests
    parallel: true
    deps: frontend:test, backend:test

task build:
    desc: Production build
    deps: frontend:build, backend:install

# Deployment
task deploy:
    desc: Deploy to production
    if: ENV == production
    deps: build, test
    remote: deploy@prod.example.com
    run:
        docker-compose pull
        docker-compose up -d

profile dev:
    env:
        ENV = development
        DEBUG = true

profile prod:
    env:
        ENV = production
        DEBUG = false
```

---

## ⚔️ FluxFile vs Other Tools

| Feature | FluxFile | Make | Taskfile | just | Mage |
|---------|----------|------|----------|------|------|
| **Syntax** | Clean YAML-like DSL | Tab-based, cryptic | YAML | Simple custom | Go code |
| **Learning Curve** | ⭐ Low | 🔴 High | ⭐ Low | ⭐ Low | 🟡 Medium |
| **Smart Caching** | ✅ Built-in tracking | Manual timestamps | ⚠️ Basic | ❌ No | ❌ No |
| **Watch Mode** | ✅ Native | ❌ External tools | ✅ Native | ❌ No | ❌ No |
| **Parallel Tasks** | ✅ Native | ❌ Manual with `-j` | ✅ Native | ❌ No | ⚠️ Manual |
| **Matrix Builds** | ✅ Yes | ❌ No | ❌ No | ❌ No | ❌ No |
| **Docker Support** | ✅ Native | ❌ No | ❌ No | ❌ No | ❌ No |
| **Remote Execution** | ✅ SSH built-in | ❌ No | ❌ No | ❌ No | ❌ No |
| **Profiles/Envs** | ✅ First-class | ❌ Manual | ✅ Native | ❌ Limited | ⚠️ Manual |
| **Cross-Platform** | ✅ Yes | ⚠️ Varies | ✅ Yes | ✅ Yes | ✅ Yes |
| **Dependencies** | ✅ Cycle detection | ✅ Basic | ✅ Basic | ✅ Basic | ⚠️ Manual |

### When to Use What

| Use Case | Best Tool |
|----------|-----------|
| **Modern projects needing caching, watch, parallel** | ✅ **FluxFile** |
| **Legacy C/C++ projects with existing Makefiles** | Make |
| **Simple scripts without caching needs** | just or Taskfile |
| **Pure Go projects wanting Go-based tasks** | Mage |
| **Docker-based builds with remote deployment** | ✅ **FluxFile** |
| **Cross-compilation matrices** | ✅ **FluxFile** |

---

## 🛠️ Task DSL Reference

### Complete Syntax

```yaml
task name:
    desc: string           # Task description
    deps: task1, task2     # Dependencies (run before this task)
    parallel: true|false   # Run dependencies in parallel
    if: VAR == value       # Conditional execution

    env:                   # Environment variables
        KEY = value
        KEY2 = ${VAR}

    run:                   # Commands to execute
        command1
        command2 ${VAR}

    # Caching & Incremental Builds
    cache: true|false      # Enable caching
    inputs:                # Files that trigger rebuild (glob patterns)
        src/**/*.go
        go.mod
    outputs:               # Build outputs (for cache validation)
        dist/binary
        build/**/*

    # Watch Mode
    watch: **/*.go         # Glob pattern to watch
    ignore:                # Patterns to ignore in watch mode
        vendor/**
        **/*_test.go
        .git/**

    # Execution Environment
    docker: true           # Run in Docker container
    remote: user@host      # Run via SSH

    # Matrix Builds
    matrix:
        os: linux, darwin, windows
        arch: amd64, arm64
```

### Variables

```yaml
# Static variable
var PROJECT = myapp

# Shell command output
var VERSION = $(shell "git describe --tags")

# Environment variable reference
var HOME_DIR = ${HOME}

# Usage in tasks
task build:
    run:
        echo "Building ${PROJECT} version ${VERSION}"
```

### Glob Patterns

| Pattern | Matches |
|---------|---------|
| `*.go` | All `.go` files in current directory |
| `**/*.go` | All `.go` files recursively |
| `src/**/*` | Everything under `src/` |
| `{*.go,*.mod}` | Files with `.go` or `.mod` extension |
| `!vendor/**` | Exclude vendor directory (in ignore) |

### Profiles

```yaml
profile dev:
    env:
        MODE = development
        LOG_LEVEL = debug
        PORT = 3000

profile prod:
    env:
        MODE = production
        LOG_LEVEL = error
        PORT = 80
```

Apply with: `flux -p dev build` or `flux -p prod deploy`

---

## 📂 Templates

### Go Project

```yaml
var PROJECT = $(shell "basename $(pwd)")
var VERSION = $(shell "git describe --tags --always")

task build:
    desc: Build Go binary
    deps: fmt, vet
    cache: true
    inputs:
        **/*.go
        go.mod
        go.sum
    outputs:
        bin/${PROJECT}
    run:
        go build -ldflags="-X main.version=${VERSION}" -o bin/${PROJECT} .

task fmt:
    desc: Format code
    run:
        go fmt ./...

task vet:
    desc: Run go vet
    run:
        go vet ./...

task test:
    desc: Run tests with coverage
    run:
        go test -v -cover ./...

task dev:
    desc: Watch and run
    watch: **/*.go
    ignore: vendor/**
    run:
        go run .

task build-all:
    desc: Cross-compile
    matrix:
        os: linux, darwin, windows
        arch: amd64, arm64
    run:
        GOOS=${os} GOARCH=${arch} go build -o dist/${PROJECT}-${os}-${arch}
```

### Node.js Project

```yaml
task install:
    desc: Install dependencies
    cache: true
    inputs:
        package.json
        package-lock.json
    outputs:
        node_modules
    run:
        npm ci

task build:
    desc: Build for production
    deps: install
    cache: true
    inputs:
        src/**/*
        tsconfig.json
    outputs:
        dist
    run:
        npm run build

task dev:
    desc: Start dev server
    deps: install
    watch: src/**/*
    run:
        npm run dev

task test:
    desc: Run tests
    deps: install
    run:
        npm test

task lint:
    desc: Lint code
    deps: install
    run:
        npm run lint
```

### Python Project

```yaml
task venv:
    desc: Create virtual environment
    cache: true
    inputs:
        requirements.txt
    outputs:
        .venv
    run:
        python -m venv .venv
        .venv/bin/pip install -r requirements.txt

task dev:
    desc: Run development server
    deps: venv
    watch: **/*.py
    ignore: .venv/**
    run:
        .venv/bin/uvicorn main:app --reload

task test:
    desc: Run pytest
    deps: venv
    run:
        .venv/bin/pytest -v

task lint:
    desc: Run linters
    deps: venv
    run:
        .venv/bin/ruff check .
        .venv/bin/mypy .

task format:
    desc: Format code
    deps: venv
    run:
        .venv/bin/black .
        .venv/bin/isort .
```

### Rust Project

```yaml
var PROJECT = $(shell "basename $(pwd)")

task build:
    desc: Build release binary
    cache: true
    inputs:
        src/**/*
        Cargo.toml
        Cargo.lock
    outputs:
        target/release/${PROJECT}
    run:
        cargo build --release

task dev:
    desc: Watch and run
    watch: src/**/*
    run:
        cargo run

task test:
    desc: Run tests
    run:
        cargo test

task check:
    desc: Check code
    run:
        cargo check
        cargo clippy -- -D warnings

task fmt:
    desc: Format code
    run:
        cargo fmt
```

---

## 💻 CLI Reference

```
flux [options] <task>

Options:
  -t string      Task to execute
  -p string      Profile to apply
  -l             List all tasks
  -w             Watch mode
  --no-cache     Disable caching
  -f string      Path to FluxFile
  -v             Show version
  --init         Initialize new FluxFile
  --template     Template (go, node, python, rust)
  --report       Show execution report
  --graph        Show dependency graph
  --dry-run      Simulate execution
  --lock         Generate lock file
  --check-lock   Verify lock file
  --lock-diff    Show lock differences
  --json         Output in JSON format
  --tui          Interactive TUI mode

Commands:
  flux init      Create FluxFile from project type
  flux logs      Open execution logs in browser
```

---

## 📊 Performance

| Component | Time | Allocations |
|-----------|------|-------------|
| Lexer | 7.6µs | 47 allocs/op |
| Parser | 10.3µs | 34 allocs/op |
| Graph Build | ~1µs | minimal |
| Executor Create | 1.1µs | 112 B/op |
| Cache Hash | <1µs | minimal |

> Run benchmarks: `cd benchmark && go test -bench Benchmark -benchmem`

---

## 🤝 Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

```bash
git clone https://github.com/ashavijit/fluxfile
cd fluxfile
make test
make build
```

---

## 📄 License

MIT © [Avijit Sen](https://github.com/ashavijit)
