# Adding FluxFile to GitHub Linguist

This guide explains how to submit FluxFile for official GitHub syntax highlighting recognition.

## Overview

GitHub uses [github/linguist](https://github.com/github/linguist) for language detection and syntax highlighting. To get FluxFile recognized like Makefile, you need to submit a PR with:

1. Language entry in `languages.yml`
2. TextMate grammar file
3. Sample files

---

## Step 1: Fork Linguist Repository

```bash
# Fork https://github.com/github/linguist on GitHub, then:
git clone https://github.com/YOUR_USERNAME/linguist.git
cd linguist
git checkout -b add-fluxfile-language
```

---

## Step 2: Add Language Entry

Edit `lib/linguist/languages.yml` and add (alphabetically under "F"):

```yaml
FluxFile:
  type: data
  color: "#6366f1"
  extensions:
  - ".flux"
  filenames:
  - FluxFile
  - Fluxfile
  - fluxfile
  tm_scope: source.fluxfile
  ace_mode: yaml
  codemirror_mode: yaml
  codemirror_mime_type: text/x-yaml
  language_id: 1090
```

### Field Explanations

| Field | Value | Description |
|-------|-------|-------------|
| `type` | `data` | Language type (programming/data/markup/prose) |
| `color` | `#6366f1` | Color shown in GitHub language stats (indigo) |
| `extensions` | `.flux` | File extensions to recognize |
| `filenames` | `FluxFile` | Exact filenames to recognize |
| `tm_scope` | `source.fluxfile` | TextMate grammar scope |
| `ace_mode` | `yaml` | Fallback editor mode |
| `language_id` | `1090` | Unique ID (check unused IDs in existing file) |

---

## Step 3: Add TextMate Grammar

Create directory and copy grammar:

```bash
mkdir -p vendor/grammars/fluxfile
```

Copy the `fluxfile.tmLanguage.json` file to `vendor/grammars/fluxfile/fluxfile.tmLanguage.json`.

The grammar file is already created at: `grammar/fluxfile.tmLanguage.json`

---

## Step 4: Add Sample Files

Create sample directory:

```bash
mkdir -p samples/FluxFile
```

Create `samples/FluxFile/example.flux`:

```yaml
# FluxFile - Modern Task Runner
# https://github.com/ashavijit/fluxfile

var PROJECT = myapp
var VERSION = $(shell "git describe --tags")
var MODE = development

task build:
    desc: Build the application
    deps: fmt, lint
    cache: true
    inputs:
        src/**/*.go
        go.mod
    outputs:
        dist/${PROJECT}
    env:
        CGO_ENABLED = 0
    run:
        go build -ldflags="-X main.version=${VERSION}" -o dist/${PROJECT} ./cmd

task test:
    desc: Run tests with coverage
    deps: build
    run:
        go test -v -cover ./...

task dev:
    desc: Watch and rebuild on changes
    watch: **/*.go
    ignore:
        vendor/**
        **/*_test.go
    run:
        go run ./cmd

task deploy:
    desc: Deploy to production
    if: MODE == production
    remote: deploy@prod.example.com
    deps: build, test
    run:
        docker-compose pull
        docker-compose up -d

task build-all:
    desc: Cross-compile for multiple platforms
    matrix:
        os: linux, darwin, windows
        arch: amd64, arm64
    run:
        GOOS=${os} GOARCH=${arch} go build -o dist/${PROJECT}-${os}-${arch}

profile dev:
    env:
        MODE = development
        LOG_LEVEL = debug

profile prod:
    env:
        MODE = production
        LOG_LEVEL = error

include "plugins/docker.flux"
```

---

## Step 5: Update Vendor Index

Add grammar to `vendor/grammars/fluxfile/`:

```bash
# Create a basic package.json for the grammar
cat > vendor/grammars/fluxfile/package.json << 'EOF'
{
  "name": "fluxfile-grammar",
  "version": "1.0.0",
  "description": "TextMate grammar for FluxFile",
  "repository": {
    "type": "git",
    "url": "https://github.com/ashavijit/fluxfile.git"
  },
  "license": "MIT"
}
EOF
```

---

## Step 6: Run Tests

```bash
# Install dependencies
bundle install

# Run Linguist tests
bundle exec rake test

# Verify language detection
bundle exec bin/github-linguist samples/FluxFile/example.flux
```

Expected output:
```
samples/FluxFile/example.flux: 100.00% FluxFile
```

---

## Step 7: Submit Pull Request

```bash
git add .
git commit -m "Add FluxFile language support

FluxFile is a modern task runner and build automation tool with a clean,
YAML-like DSL for defining tasks, dependencies, and workflows.

- Homepage: https://github.com/ashavijit/fluxfile
- Grammar: source.fluxfile
- Extensions: .flux
- Filenames: FluxFile, Fluxfile, fluxfile"

git push origin add-fluxfile-language
```

Then create a PR on GitHub with:

**Title:** `Add FluxFile language support`

**Description:**
```markdown
## Language Information

**FluxFile** is a modern task runner and build automation tool with a clean, minimal syntax.

- **Homepage**: https://github.com/ashavijit/fluxfile
- **Documentation**: https://github.com/ashavijit/fluxfile#readme
- **Repository Stars**: [current stars]
- **Usage**: Task automation, build systems, CI/CD pipelines

## Features

- Clean YAML-like DSL for defining tasks
- Dependency resolution with cycle detection
- Smart caching with input/output tracking
- File watching for auto-rebuild
- Matrix builds for cross-compilation
- Docker and remote SSH execution

## Checklist

- [x] Added entry to `lib/linguist/languages.yml`
- [x] Added TextMate grammar to `vendor/grammars/fluxfile/`
- [x] Added sample file to `samples/FluxFile/`
- [x] Tests pass locally
```

---

## Timeline

- **Review**: 1-4 weeks (Linguist team is busy)
- **Merge**: After approval
- **Deploy**: Next Linguist release (usually within a month)

---

## Temporary Workaround

Until your PR is merged, use `.gitattributes` in your repos:

```
FluxFile linguist-language=YAML
*.flux linguist-language=YAML
```

This makes GitHub render FluxFile with YAML highlighting immediately.

---

## Resources

- [Linguist Contributing Guide](https://github.com/github/linguist/blob/master/CONTRIBUTING.md)
- [Adding a Language](https://github.com/github/linguist/blob/master/CONTRIBUTING.md#adding-a-language)
- [TextMate Grammar Guide](https://macromates.com/manual/en/language_grammars)
