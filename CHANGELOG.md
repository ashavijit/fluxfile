# Changelog

All notable changes to the Flux project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.1.1] - 2025-12-04

### Changed
- **Command Logging**: Improved command execution logs to include timestamps and better color formatting for improved visibility.
- **Build System**: Fixed build command in documentation and scripts to correctly build the entire `cmd/flux` package.

## [2.1.0] - 2025-12-03

### Added
- **Comprehensive Testing Suite**: Added `FluxFile.test_features` and verified all CLI commands (basic, feature-specific, lockfile management).
- **Refactored Parser Logic**: Implemented tagged switch statements in `internal/parser` for improved performance and readability.
- **Lint Fixes**: Resolved unused code warnings and optimized switch statements across `internal/lexer`, `internal/parser`, and `internal/executor`.

### Changed
- **Codebase Optimization**: Refactored `parseRetries`, `parseDesc`, `parseWatch`, and other parser functions to use tagged switches.
- **Cleaned Up Code**: Removed unused methods (`expectNewline`, `peekChar`) and constants (`colorMagenta`).

## [2.0.0] - 2025-12-02

### Added

#### Lockfile v2.0 - Major Enhancement
- **Comprehensive Metadata Tracking**
  - System information (OS, architecture, Go version)
  - Environment context (hostname, username)
  - FluxFile integrity hash to detect FluxFile modifications
  - Per-task update timestamps

- **Advanced Change Detection**
  - Task configuration hashing (dependencies, environment variables, conditions)
  - Command hashing for run directives
  - File metadata tracking (size, modification time, content hash)
  - Detailed diff output showing exact changes

#### New CLI Commands
- `--lock-update --task <name>` - Update specific tasks without regenerating entire lockfile
- `--lock-diff` - Show detailed differences between lock and current state
- `--lock-clean` - Remove stale tasks that no longer exist in FluxFile
- `--json` - Machine-readable JSON output for all lock commands

#### Documentation
- **FLUX_CLI_REFERENCE.md** - Comprehensive command reference with:
  - All Flux commands documented A-Z
  - Real tested outputs from actual command execution
  - FluxFile syntax examples for all features
  - Lock file v2.0 format specification
  - Workflows and best practices
  - CI/CD integration examples
  - Troubleshooting guide
  - Tips and tricks section

### Changed

- **Lockfile Format Version Bump**: v1.0 → v2.0
  - `TaskLock` structure enhanced with `config_hash`, `command_hash`, and `last_updated`
  - File tracking upgraded from simple hash strings to `FileInfo` objects
  - Added top-level `Metadata` and `FluxFileHash` fields

- **Enhanced Output Formatting**
  - Lock generation shows detailed statistics (total inputs/outputs tracked)
  - Verification displays file size changes in mismatch reports
  - Update commands show abbreviated hashes for quick reference
  - Better error messages with actionable suggestions

- **Improved Lock Verification**
  - Size information included in hash mismatch reports
  - Helpful suggestions to run `--lock-diff` for details
  - JSON output support for programmatic parsing

### Fixed

- Updated `.gitignore` to exclude `COMMANDS.md` (internal documentation)

### Breaking Changes

⚠️ **Lockfile Format**: v1.0 lockfiles are not compatible with v2.0. To migrate:
1. Delete old `FluxFile.lock`
2. Run `flux --lock` to generate new v2.0 lockfile
3. Commit the new lockfile to version control

---

## [1.0.0-beta] - 2025-12-01

### Added

#### Interactive TUI Mode
- Interactive terminal UI with real-time task status
- Visual task execution monitoring
- Clean, stable display without scrolling
- In-place updates using ANSI cursor controls
- Status indicators: ⏸ pending, ✓ completed, ✗ failed
- Execution time tracking per task

#### Lock File v1.0 (Initial Implementation)
- Generate dependency lock files with `--lock` flag
- Verify lock files with `--check-lock` flag
- Track input and output file hashes
- Detect file changes for incremental builds

#### Enhanced Task Display
- `flux show` command for formatted task list display
- Box-drawing characters for professional UI
- Dependency count indicators
- Task descriptions displayed
- Total task count summary

### Fixed
- TUI color constant conflicts resolved
- Removed duplicate color definitions in `main.go`
- TUI scrolling issues - now updates in-place
- Clean, stable TUI display

### Documentation
- Added competitive feature roadmap
- Added commands reference documentation
- Enhanced UI documentation

---

## [1.0.0] - 2025-11-30

### Added

#### Core Features - Phase 1 & 2

**Phase 1: Essential Features**
- Task dependency resolution with cycle detection
- Task result caching based on file hashes
- Enhanced caching with input/output tracking
- File watching with ignore patterns for automatic re-execution
- Conditional task execution based on environment
- Parallel task execution for dependencies
- Task descriptions for better documentation

**Phase 2: Advanced Features**
- Matrix builds for multi-platform compilation
- Docker container execution support
- Remote execution over SSH
- Variable expansion with shell command execution
- Profile support for environment-specific configuration
- Include directive for modular FluxFiles

**Security & Resilience**
- Secret management for sensitive environment variables
- Preconditions (file existence, command availability, env var checks)
- Retry logic with configurable attempts and delays
- Timeout support for long-running tasks

#### FluxFile Parser
- Complete parsing infrastructure for FluxFile syntax
- AST (Abstract Syntax Tree) generation
- Support for all task directives:
  - `desc` - Task descriptions
  - `deps` - Task dependencies
  - `parallel` - Parallel execution
  - `if` - Conditional execution
  - `env` - Environment variables
  - `run` - Command execution
  - `watch` - File watching patterns
  - `ignore` - Ignore patterns for watch
  - `cache` - Caching control
  - `inputs` - Input file patterns
  - `outputs` - Output file patterns
  - `matrix` - Matrix build dimensions
  - `docker` - Docker execution
  - `remote` - Remote SSH execution
  - `secrets` - Secret management
  - `pre` - Preconditions
  - `retries` - Retry configuration
  - `retry_delay` - Delay between retries
  - `timeout` - Task timeout

#### Task Executor
- Working executor with all Phase 1 & 2 features
- Conditional evaluation with proper operator parsing
- Spaced operator support (`==`, `!=`, `>`, `<`, `>=`, `<=`)
- Environment variable expansion
- Shell command execution
- Dependency graph resolution
- Parallel execution support
- Caching mechanism

### CLI Interface
- `-t <task>` - Execute specific task
- `-p <profile>` - Apply environment profile
- `-l` - List all available tasks
- `-w` - Watch mode for automatic re-execution
- `--no-cache` - Disable task caching
- `-f <path>` - Specify custom FluxFile path
- `-v` - Show version information

### Installation
- Install script for Linux/macOS
- Install script for Windows (PowerShell)
- Makefile for building from source

### Documentation
- README.md with features, installation, and usage
- Quick start guide
- Syntax reference
- CLI usage documentation
- Examples for common use cases

### Build & Release
- GitHub Actions workflow for releases
- Automated builds for multiple platforms
- Version tagging (v1.0.0)

---

## Initial Commits

### Added
- Basic project structure
- Go module initialization
- Core AST types
- Parser foundation
- Executor framework
- Line ending normalization

---

## Version Comparison

### v2.0.0 vs v1.0.0

| Feature | v1.0.0 | v2.0.0 |
|---------|--------|--------|
| **Lockfile Version** | 1.0 | 2.0 |
| **Metadata Tracking** | None | Full (OS, arch, Go version, user, hostname) |
| **File Information** | Hash only | Hash + size + modification time |
| **Task Config Tracking** | ❌ | ✅ (config_hash + command_hash) |
| **FluxFile Integrity** | ❌ | ✅ (fluxfile_hash) |
| **Selective Updates** | ❌ | ✅ (`--lock-update`) |
| **Detailed Diff** | ❌ | ✅ (`--lock-diff`) |
| **Cleanup Utility** | ❌ | ✅ (`--lock-clean`) |
| **JSON Output** | ❌ | ✅ (`--json`) |
| **Per-Task Timestamps** | ❌ | ✅ (`last_updated`) |
| **Interactive TUI** | ✅ | ✅ |
| **Enhanced Display** | ✅ | ✅ |

---

## Migration Guides

### Upgrading to v2.0.0 from v1.x

1. **Backup your current lockfile** (optional):
   ```bash
   cp FluxFile.lock FluxFile.lock.v1.backup
   ```

2. **Delete the old lockfile**:
   ```bash
   rm FluxFile.lock
   ```

3. **Generate new v2.0 lockfile**:
   ```bash
   flux --lock
   ```

4. **Verify the new lockfile**:
   ```bash
   flux --check-lock
   ```

5. **Commit to version control**:
   ```bash
   git add FluxFile.lock
   git commit -m "chore: upgrade to lockfile v2.0"
   ```

---

## Contributors

- [@ashavijit](https://github.com/ashavijit) - Creator and maintainer

---

## Links

- [Repository](https://github.com/ashavijit/fluxfile)
- [Issues](https://github.com/ashavijit/fluxfile/issues)
- [Releases](https://github.com/ashavijit/fluxfile/releases)

---

## Notes

### Semantic Versioning

- **MAJOR** version (X.0.0): Incompatible API changes or breaking changes
- **MINOR** version (0.X.0): New features in a backward-compatible manner
- **PATCH** version (0.0.X): Backward-compatible bug fixes

### Release Types

- **Alpha**: Early testing, unstable
- **Beta**: Feature complete, testing phase
- **RC** (Release Candidate): Final testing before stable release
- **Stable**: Production-ready release

---

**Keep a Changelog** format ensures:
- Readable by humans and machines
- Chronological order (newest first)
- Clear categorization (Added, Changed, Deprecated, Removed, Fixed, Security)
- Release dates in ISO format (YYYY-MM-DD)
- Semantic versioning references

---

*This changelog is maintained by the Flux team and updated with each release.*
