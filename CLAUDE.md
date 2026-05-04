# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**stevedore** is a Go CLI tool that analyzes a directory against a `.dockerignore` file, reporting which files will be included in or excluded from a Docker image. Named after dock workers who load/unload cargo.

## Commands

```bash
# Build
go build -v ./...

# Test (all)
go test -v ./...

# Run a single test suite
go test -v -run TestArguments
go test -v -run TestFails
go test -v -run TestConfig

# Test with coverage
go test -cover ./...

# Lint
golangci-lint run

# Run the tool locally
./stevedore .
./stevedore --verbose .
cat .dockerignore | ./stevedore --stdin
```

## Architecture

The entire application logic lives in **`main.go`** (~410 lines) — there is no package split. Tests are in `main_test.go`.

**Config struct** centralizes all options: color flags, output filters (`--included`, `--excluded`), path display (`--fullpath`), verbosity, and ignore file path. Configuration is layered:

1. Global: `$XDG_CONFIG_HOME/stevedore/config.json` (fallback: `~/.config/stevedore/config.json`)
2. Local: `.stevedore.json` in the working directory
3. CLI flags (highest priority)

**Core flow:**
- Loads `.dockerignore` (or `--ignorefile` target) via `go-gitignore`
- Optionally loads `.stevedoreignore` to filter the output listing itself
- Walks the directory with `filepath.Walk()`
- Matches each path against ignore patterns, accounting for directory trailing-slash semantics
- Color-codes output: included files vs. excluded files

**Exit codes:**
- `0` — success
- `1` — ignore file not found or unreadable
- `2` — directory read error

## Testing Patterns

The three test suites in `main_test.go`:

- **TestArguments** — exercises every CLI flag combination against the project's own directory; all expect exit code 0.
- **TestFails** — creates temp directories with permission-restricted ignore files to trigger exit code 1.
- **TestConfig** — validates config file loading with a temporary `.stevedore.json`.

Tests call `main()` directly via `os.Args` manipulation and capture behavior through exit codes. When adding new flags, add a corresponding case to `TestArguments`.

## Key Dependencies

- `github.com/fatih/color` — terminal color output
- `github.com/sabhiram/go-gitignore` — `.gitignore`/`.dockerignore` pattern matching

## Pre-commit Hooks

The repo uses `.pre-commit-config.yaml` with gitleaks (secret scanning), golangci-lint, end-of-file-fixer, and trailing-whitespace checks. Run `pre-commit run --all-files` before submitting changes if hooks are installed.
