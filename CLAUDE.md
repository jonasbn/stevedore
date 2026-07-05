# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with
code in this repository.

## Project Overview

**stevedore** is a Go CLI tool that analyzes a directory against a
`.dockerignore` file, reporting which files will be included in or excluded
from a Docker image. Named after dock workers who load/unload cargo.

## Commands

```bash
# Build all packages (CI / compilation check)
go build -v ./...

# Build the local executable used below
go build -o stevedore .

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

The entire application logic lives in **`main.go`** (~410 lines) — there is no
package split. Tests are in `main_test.go`.

**Config struct** centralizes all options: color flags, output filters
(`--included`, `--excluded`), path display (`--fullpath`), verbosity, and
ignore file path. Configuration is layered:

1. Global: `$XDG_CONFIG_HOME/stevedore/config.json`
   (fallback: `~/.config/stevedore/config.json`)
2. Local: `.stevedore.json` in the working directory
3. CLI flags (highest priority)

**Core flow:**

- Loads `.dockerignore` (or `--ignorefile` target) via `go-gitignore`
- Optionally loads `.stevedoreignore` to filter the output listing itself
- Walks the directory with `filepath.Walk()`
- Matches each path against ignore patterns, accounting for directory
  trailing-slash semantics
- Color-codes output: included files vs. excluded files

**Exit codes:**

- `0` — success
- `1` — ignore file not found or unreadable
- `2` — directory read error

## Testing Patterns

The three test suites in `main_test.go`:

- **TestArguments** — enumerates individual CLI flag cases against the
  project's own directory; these cases expect successful execution (exit
  code 0).
- **TestFails** — creates temp directories with permission-restricted ignore
  files to trigger exit code 1.
- **TestConfig** — exercises `realMain()` in the `tests/ok` fixture context
  after setting up `tests/ok/.dockerignore`; it does not currently create or
  validate loading a temporary `.stevedore.json`.

Tests call `main()` directly via `os.Args` manipulation and capture behavior
through exit codes. When adding new flags, consider adding or updating a
focused case in `TestArguments`.

## Key Dependencies

- `github.com/fatih/color` — terminal color output
- `github.com/sabhiram/go-gitignore` — `.gitignore`/`.dockerignore` pattern
  matching

## CI / GitHub Actions

Workflows live under `.github/workflows/`:

- **build** — `go build`/`go test` on push and PR.
- **codeql.yml** — CodeQL analysis (go, ruby). The `init`, `autobuild`, and
  `analyze` steps must all be pinned to the *same* `codeql-action` version.
  Dependabot bumps these as separate PRs, one per sub-action; merging just
  one alone leaves the others behind and breaks CI with `Loaded a
  configuration file for version X, but running version Y`. Bump `init`,
  `autobuild`, and `analyze` together (consolidate the Dependabot PRs into
  one branch/PR) rather than merging them individually.
- **scorecard.yml** — OSSF Scorecard; uses `codeql-action/upload-sarif`,
  which is independent of the init/autobuild/analyze version-sync
  constraint above and can be bumped on its own.
- **Spellcheck** — `rojopolis/spellcheck-github-actions`, config at
  `.github/spellcheck.yml`, checks all `**/*.md` files (including headings)
  against `.github/spellchecker-wordlist.txt`. Adding a doc with a
  project-specific term or an all-caps heading (e.g. `TODO`) that isn't a
  dictionary word requires adding it to the wordlist or the job fails.
- **Markdownlint** — lints all Markdown files.
- **dependency-review** — runs on PRs that change dependencies.
- Copilot's automated PR reviewer and the Coveralls coverage bot both
  comment on PRs automatically; check for these before assuming a human
  reviewer commented.

## Documentation

- `docs/TODO.md` — checklist mirroring the open GitHub enhancement issues.
  Keep it in sync (add/check off entries) when issues are opened or closed,
  if asked to.

## Pre-commit Hooks

The repo uses `.pre-commit-config.yaml` with gitleaks (secret scanning),
golangci-lint, end-of-file-fixer, and trailing-whitespace checks. Run
`pre-commit run --all-files` before submitting changes if hooks are installed.
