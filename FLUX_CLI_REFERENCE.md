# Flux Command Line Reference

> Complete documentation of all Flux CLI commands with real output examples

**Version:** 1.0.0  
**Lock Format:** 2.0  
**Last Updated:** 2025-12-02

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Basic Commands](#basic-commands)
3. [Task Execution](#task-execution)
4. [Lock File Management](#lock-file-management)
5. [Advanced Features](#advanced-features)
6. [Profiles](#profiles)
7. [Watch Mode](#watch-mode)
8. [Command Reference Table](#command-reference-table)
9. [Exit Codes](#exit-codes)
10. [Examples & Workflows](#examples--workflows)

---

## Getting Started

### Installation

Flux is installed and available in your PATH. Verify with:

```bash
flux -v
```

### Quick Start

```bash
# List all available tasks
flux -l

# Run a task
flux <task-name>

# Watch for changes
flux -w <task-name>

# Generate lock file
flux --lock
```

---

## Basic Commands

### Show Version

Display the current Flux version.

**Command:**
```bash
flux -v
```

**Output:**
```
Flux version 1.0.0
```

---

### Show Help

Display command-line usage and available flags.

**Command:**
```bash
flux --help
```

**Available Flags:**
```
  -check-lock
        Verify lock file
  -f string
        Path to FluxFile
  -json
        Output in JSON format
  -l    List all tasks
  -lock
        Generate dependency lock file
  -lock-clean
        Remove stale tasks from lock file
  -lock-diff
        Show detailed diff between lock and current state
  -lock-update
        Update specific task in lock file
  -no-cache
        Disable caching
  -p string
        Profile to apply
  -show
        Show all tasks with enhanced UI
  -t string
        Task to execute
  -task string
        Task name for --lock-update
  -tui
        Run interactive TUI mode
  -v    Show version
  -w    Watch mode
```

---

### List Tasks (Simple)

List all available tasks in the FluxFile.

**Command:**
```bash
flux -l
```

**Output:**
```
Available tasks:
  - build
  - fmt
  - lint
  - dev
  - build-all
  - docker-build
  - deploy
  - test
  - clean
  - docker-shell
  - docker-logs
  - docker-clean
```

---

### Show Tasks (Enhanced UI)

Display tasks with an enhanced formatted interface.

**Command:**
```bash
flux show
```

or

```bash
flux -show
``​`

**Output:**
```
╔════════════════════════════════════════════════════════════════╗
║                        AVAILABLE TASKS                         ║
╚════════════════════════════════════════════════════════════════╝

  TASK             DESCRIPTION
  ─────────────────────────────────────────────────────────────────────
  build            (no description) [→ 2 deps]
  fmt              (no description)
  lint             (no description)
  dev              (no description)
  build-all        (no description)
  docker-build     (no description)
  deploy           (no description) [→ 1 deps]
  test             (no description)
  clean            (no description)
  docker-shell     (no description)
  docker-logs      (no description)
  docker-clean     (no description)

  Total: 12 tasks

  Run a task: flux <task>
```

**Note:** Tasks with dependencies show `[→ N deps]`

---

### List Tasks with Custom FluxFile

Use a different FluxFile than the default.

**Command:**
```bash
flux -f FluxFile.test2 -l
```

**Output:**
```
Available tasks:
  build                Build the project
  test                 Run tests
  lint                 Run linter
  clean                Clean build artifacts
```

**Note:** This FluxFile has task descriptions defined.

---

## Task Execution

### Run a Task (Basic)

Execute a single task.

**Command:**
```bash
flux fmt
```

**Output:**
```
[→] Running task: fmt
  $ go fmt ./...

  Completed in 874.9772ms

[✓] Task fmt completed
```

---

### Run a Task with -t Flag

Alternative syntax using the `-t` flag.

**Command:**
```bash
flux -t fmt
```

**Output:** Same as above

---

### Run a Task (With Dependencies)

When a task has dependencies, they run first.

**Command:**
```bash
flux build
```

**Output:**
```
[→] Running task: fmt
  $ go fmt ./...
  Completed in 800ms
[✓] Task fmt completed

[→] Running task: lint
  $ golangci-lint run
  [Error if linter not installed]

[→] Running task: build
  $ go build -o dist/flux ./cmd/flux
  Completed in 2.5s
[✓] Task build completed
```

---

### Run Task with No Cache

Disable caching to force re-execution.

**Command:**
```bash
flux --no-cache build
```

**Note:** Forces task to run even if inputs haven't changed.

---

### Task Execution Errors

When a task fails, Flux shows the error clearly.

**Example - Missing Command:**
```bash
flux lint
```

**Output:**
```
[→] Running task: lint
  $ golangci-lint run
  golangci-lint : The term 'golangci-lint' is not recognized...
  
[✗] Task lint failed: command failed: exit status 1
[ERROR] command failed: exit status 1
```

---

### No Task Specified Error

Running flux without arguments shows an error.

**Command:**
```bash
flux
```

**Output:**
```
[ERROR] No task specified. Use -t <task> or provide task name as argument
```

---

## Lock File Management

### Generate Lock File

Create a new lock file capturing current state.

**Command:**
```bash
flux --lock
```

**Output:**
```
[✓] Lock file generated: FluxFile.lock (v2.0)
    Generated: 2025-12-02 16:42:10
    OS/Arch: windows/amd64
    Tasks locked: 0
    Total: 0 inputs, 0 outputs tracked
```

**Note:** v2.0 format includes comprehensive metadata.

---

### Generate Lock with JSON Output

Output lock file in machine-readable JSON format.

**Command:**
```bash
flux --lock --json
```

**Output:**
```json
{
    "version":  "2.0",
    "generated":  "2025-12-02T16:43:13.9431649+05:30",
    "metadata":  {
        "fluxfile_path":  "FluxFile",
        "hostname":  "Avijit",
        "user":  "aviji",
        "go_version":  "go1.24.1",
        "os":  "windows",
        "arch":  "amd64"
    },
    "fluxfile_hash":  "a1aa4309a9d6bc1b6ea9e8e1392d880b27a29e7b9e0b3c6e11b4a973ee7e148a",
    "tasks":  {}
}
```

---

### Verify Lock File

Check if lock file matches current state.

**Command:**
```bash
flux --check-lock
```

**Output (Valid):**
```
[✓] Lock file verified - all files match
```

**Output (Changes Detected):**
```
[⚠] Lock file verification failed - 2 task(s) changed:

  Task: build
    - input src/main.go: hash mismatch (size: 1024 -> 1156)
    - output dist/app: missing

Run 'flux --lock-diff' for detailed differences
```

---

### Show Detailed Diff

View detailed differences between lock and current state.

**Command:**
```bash
flux --lock-diff
```

**Output (No Differences):**
```
[✓] No differences detected
```

**Output (With Differences):**
```
[!] Found differences in 2 task(s):

Task: build
  [~] Task configuration changed
  [~] Run commands changed

Task: test
  [~] Task configuration changed
  [~] Run commands changed
```

**Symbols:**
- `[~]` - Modified
- `[-]` - Missing/Deleted
- `[+]` - New/Added

---

### Update Specific Task in Lock

Update only a specific task without regenerating entire lock.

**Command:**
```bash
flux --lock-update --task build
```

**Output:**
```
[✓] Updated task 'build' in lock file
    Inputs: 5 files
    Outputs: 1 files
    Config hash: c8b25d10130e
    Command hash: b207bb14fb20
```

**Error (Task Not Found):**
```
[ERROR] task 'unknown' not found in FluxFile
```

**Error (Missing Task Name):**
```
[ERROR] Task name required for --lock-update
```

---

### Clean Stale Tasks from Lock

Remove tasks from lock that no longer exist in FluxFile.

**Command:**
```bash
flux --lock-clean
```

**Output (Nothing to Clean):**
```
[✓] No stale tasks to clean
```

**Output (Tasks Removed):**
```
[✓] Removed 3 stale task(s) from lock file
```

---

## Advanced Features

### Use Custom FluxFile

Specify a different FluxFile path.

**Command:**
```bash
flux -f /path/to/FluxFile.custom build
```

or

```bash
flux -f FluxFile.prod --lock
```

---

### Interactive TUI Mode

Launch interactive terminal UI (if implemented).

**Command:**
```bash
flux --tui
```

or

```bash
flux -tui -t build
```

**Note:** TUI provides real-time task monitoring and visualization.

---

## Profiles

### List Tasks with Profile

Apply a profile when listing tasks.

**Command:**
```bash
flux -p dev -l
```

**Output:**
```
Available tasks:
  - build
  - fmt
  - lint
  ...
```

**Note:** Environment variables from the profile are applied.

---

### Run Task with Profile

Execute a task with specific profile settings.

**Command:**
```bash
flux -p prod deploy
```

**FluxFile Profile Definition:**
```
profile dev:
    env:
        MODE = dev
        LOG = debug

profile prod:
    env:
        MODE = production
        LOG = error
```

**Effect:** The `MODE` and `LOG` environment variables are set during execution.

---

## Watch Mode

### Watch for File Changes

Automatically re-run task when files change.

**Command:**
```bash
flux -w dev
```

**FluxFile Task:**
```
task dev:
    watch: **/*.go
    ignore:
        vendor/**
        .git/**
    run:
        go run ./cmd/flux
```

**Output:**
```
[→] Starting watch mode for task: dev
[→] Running task: dev
  $ go run ./cmd/flux
  
[✓] Task dev completed

Watching for changes in: **/*.go
```

**Note:** Task re-runs automatically when matching files change.

---

## Command Reference Table

| Command | Flags | Description | Output |
|---------|-------|-------------|--------|
| `flux -v` | | Show version | Version number |
| `flux --help` | | Display help | Flag documentation |
| `flux -l` | | List tasks | Simple task list |
| `flux show` | `-show` | Enhanced task list | Formatted UI |
| `flux <task>` | | Run task | Execution output |
| `flux -t <task>` | | Run task (alt syntax) | Execution output |
| `flux -f <path> <task>` | `-f` | Use custom FluxFile | Execution output |
| `flux -p <profile> <task>` | `-p` | Apply profile | Execution output |
| `flux -w <task>` | `-w` | Watch mode | Continuous monitoring |
| `flux --no-cache <task>` | `--no-cache` | Disable caching | Forced execution |
| `flux --lock` | `--lock` | Generate lock | Lock file created |
| `flux --check-lock` | `--check-lock` | Verify lock | Validation result |
| `flux --lock-diff` | `--lock-diff` | Show differences | Detailed diff |
| `flux --lock-update --task <name>` | `--lock-update`, `--task` | Update task in lock | Update confirmation |
| `flux --lock-clean` | `--lock-clean` | Clean stale tasks | Cleanup result |
| `flux --json` | `--json` | JSON output | Machine-readable |
| `flux --tui` | `--tui` | Interactive TUI | Terminal UI |

---

## Exit Codes

| Code | Meaning | Example |
|------|---------|---------|
| `0` | Success | Task completed successfully |
| `1` | Error | Task failed, command not found, invalid arguments |

**Examples:**

```bash
# Success
flux fmt
echo $?  # 0

# Failure
flux nonexistent-task
echo $?  # 1

# Command failure
flux lint  # If golangci-lint not installed
echo $?  # 1
```

---

## Examples & Workflows

### Workflow 1: Basic Development

```bash
# 1. List available tasks
flux -l

# 2. Format code
flux fmt

# 3. Run tests
flux test

# 4. Build
flux build
```

---

### Workflow 2: Working with Lock Files

```bash
# 1. Generate initial lock
flux --lock

# 2. Make code changes
# ... edit files ...

# 3. Check what changed
flux --lock-diff

# 4. Verify lock
flux --check-lock  # Will show changes

# 5. Update specific task
flux --lock-update --task build

# 6. Verify again
flux --check-lock  # Should pass
```

---

### Workflow 3: Multiple Environments

```bash
# Development
flux -p dev build
flux -p dev test

# Production
flux -p prod deploy

# Different FluxFile
flux -f FluxFile.prod --lock
flux -f FluxFile.prod -l
```

---

### Workflow 4: CI/CD Integration

```bash
#!/bin/bash
# CI Pipeline

# Verify lock is up to date
if ! flux --check-lock; then
    echo "Lock file out of date"
    flux --lock-diff
    exit 1
fi

# Run tests
flux test || exit 1

# Build
flux build || exit 1

echo "CI passed ✓"
```

---

### Workflow 5: Lock Maintenance

```bash
# After removing tasks from FluxFile
flux --lock-clean

# After major refactoring
flux --lock

# Export lock for analysis
flux --lock --json > lock-snapshot.json

# Verify specific FluxFile
flux -f FluxFile.prod --check-lock
```

---

## FluxFile Format Examples

### Basic Task
```
task build:
    run:
        go build -o app ./cmd
```

### Task with Dependencies
```
task build:
    deps: fmt, lint
    run:
        go build -o app ./cmd
```

### Task with Description
```
task build:
    desc: Build the application
    run:
        go build -o app ./cmd
```

### Task with Inputs/Outputs
```
task build:
    inputs:
        src/**/*.go
        go.mod
    outputs:
        dist/app
    run:
        go build -o dist/app ./cmd
```

### Task with Environment
```
task build:
    env:
        GO111MODULE = on
        CGO_ENABLED = 0
    run:
        go build -o app ./cmd
```

### Task with Watch
```
task dev:
    watch: **/*.go
    ignore:
        vendor/**
        .git/**
    run:
        go run ./cmd
```

### Matrix Build
```
task build-all:
    matrix:
        os: linux, darwin, windows
        arch: amd64, arm64
    run:
        GOOS=${os} GOARCH=${arch} go build -o dist/app-${os}-${arch}
```

### Docker Task
```
task docker-build:
    docker: true
    run:
        docker build -t myapp:latest .
```

### Remote Execution
```
task deploy:
    remote: user@prod.server.com
    run:
        docker-compose pull
        docker-compose up -d
```

### Profile Definition
```
profile dev:
    env:
        MODE = development
        DEBUG = true
        LOG_LEVEL = debug

profile prod:
    env:
        MODE = production
        DEBUG = false
        LOG_LEVEL = error
```

### Variables
```
var PROJECT = myapp
var VERSION = $(shell "git describe --tags")

task build:
    run:
        go build -ldflags "-X main.version=${VERSION}" -o ${PROJECT}
```

---

## Lock File Format (v2.0)

### Structure

```json
{
  "version": "2.0",
  "generated": "2025-12-02T16:42:10+05:30",
  "metadata": {
    "fluxfile_path": "FluxFile",
    "hostname": "machine-name",
    "user": "username",
    "go_version": "go1.24.1",
    "os": "windows",
    "arch": "amd64"
  },
  "fluxfile_hash": "sha256-hash-of-fluxfile",
  "tasks": {
    "build": {
      "config_hash": "sha256-of-task-config",
      "command_hash": "sha256-of-commands",
      "inputs": {
        "src/main.go": {
          "hash": "sha256-of-file",
          "size": 1024,
          "mod_time": "2025-12-02T15:00:00Z"
        }
      },
      "outputs": {
        "dist/app": {
          "hash": "sha256-of-file",
          "size": 4096,
          "mod_time": "2025-12-02T15:05:00Z"
        }
      },
      "hash": "combined-hash",
      "last_updated": "2025-12-02T16:42:10+05:30"
    }
  }
}
```

### Key Features

- **Metadata**: System info, environment context
- **FluxFile Hash**: Detects FluxFile modifications
- **Config Hash**: Tracks task definition changes
- **Command Hash**: Tracks run command changes
- **File Info**: Hash + size + modification time
- **Per-Task Timestamps**: When each task was last updated

---

## Best Practices

### 1. Always Commit Lock Files
```bash
git add FluxFile.lock
git commit -m "Update lock file"
```

### 2. Verify in CI
```bash
flux --check-lock || exit 1
```

### 3. Use Profiles for Environments
```bash
flux -p dev test
flux -p prod deploy
```

### 4. Use Selective Updates
```bash
# Faster than full regeneration
flux --lock-update --task build
```

### 5. Regular Lock Maintenance
```bash
# Clean after refactoring
flux --lock-clean

# Full regeneration periodically
flux --lock
```

### 6. Watch Mode for Development
```bash
flux -w dev  # Auto-reload on changes
```

---

## Troubleshooting

### Task Not Found
```
[ERROR] No task specified...
```
**Solution:** Check task name with `flux -l`

### Command Not Recognized
```
golangci-lint : The term 'golangci-lint' is not recognized...
```
**Solution:** Install required tool or update PATH

### Lock Verification Failed
```
[⚠] Lock file verification failed...
```
**Solution:** Run `flux --lock-diff` to see changes, then `flux --lock-update` or `flux --lock`

### No FluxFile Found
```
[ERROR] FluxFile not found
```
**Solution:** Create FluxFile in current directory or use `-f` to specify path

---

## Tips & Tricks

### Quick Task Execution
```bash
# Short form
flux build

# Long form (same result)
flux -t build
```

### Combine Flags
```bash
flux -p prod -f FluxFile.prod deploy
```

### JSON Output for Scripts
```bash
LOCK_DATA=$(flux --lock --json)
echo $LOCK_DATA | jq '.metadata'
```

### Chain Tasks
```
task ci:
    deps: fmt, lint, test, build
    run:
        echo "All checks passed"
```

### Debug with Watch
```bash
flux -w -p dev test  # Watch + profile
```

---

## Changelog

### Version 2.0 (Lock Format)
- ✅ Enhanced metadata tracking
- ✅ Configuration and command hashing
- ✅ File timestamps and sizes
- ✅ Selective task updates
- ✅ Detailed diff output
- ✅ Stale task cleanup
- ✅ JSON output support

### Version 1.0
- Initial release with basic features

---

## License

MIT License

---

**Flux CLI** - Modern task runner and build automation tool  
For more information, visit the project repository or run `flux --help`
