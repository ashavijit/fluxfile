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
  <a href="https://go.dev/"><img src="https://img.shields.io/github/go-mod-go-version/ashavijit/fluxfile" alt="Go Version"></a>
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

## 🚀 Quick Install

```bash
# Linux / macOS
curl -fsSL https://raw.githubusercontent.com/ashavijit/fluxfile/main/scripts/install.sh | sh

# Windows (PowerShell)
iwr -useb https://raw.githubusercontent.com/ashavijit/fluxfile/main/scripts/install.ps1 | iex

# From Source
git clone https://github.com/ashavijit/fluxfile && cd fluxfile && make install
```

## 📖 Quick Start

Create a `FluxFile` in your project:

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

## 📊 Performance Benchmarks

| Component | Time | Allocations |
|-----------|------|-------------|
| Lexer | 7.6µs | 47 allocs/op |
| Parser | 10.3µs | 34 allocs/op |
| Graph Build | ~1µs | minimal |
| Executor Create | 1.1µs | 112 B/op |
| Cache Hash | <1µs | minimal |
| Large File (100x) | 450ms | 1.8MB |

> Run benchmarks: `cd benchmark && go test -bench Benchmark -benchmem`

## 🛠️ Syntax Reference

### Variables
```yaml
var PROJECT = myapp
var VERSION = $(shell "git describe --tags")
```

### Tasks
```yaml
task name:
    desc: Task description
    deps: dep1, dep2
    parallel: true
    if: MODE == prod
    env:
        KEY = value
    run:
        command1
        command2
    cache: true
    inputs:
        src/**/*.go
    outputs:
        dist/binary
```

### Matrix Builds
```yaml
task cross-compile:
    matrix:
        os: linux, darwin, windows
        arch: amd64, arm64
    run:
        GOOS=${os} GOARCH=${arch} go build -o dist/app-${os}-${arch}
```

### Docker Execution
```yaml
task docker-test:
    docker: true
    run:
        npm install
        npm test
```

### Remote Deployment
```yaml
task deploy:
    remote: deploy@prod.example.com
    run:
        docker-compose pull
        docker-compose up -d
```

### Profiles
```yaml
profile dev:
    env:
        MODE = development

profile prod:
    env:
        MODE = production
```

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

Commands:
  flux init      Create FluxFile from project type
  flux logs      Open execution logs in browser
```

## 🤝 Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

```bash
git clone https://github.com/ashavijit/fluxfile
cd fluxfile
make test
make build
```

## 📄 License

MIT © [Avijit Sen](https://github.com/ashavijit)
