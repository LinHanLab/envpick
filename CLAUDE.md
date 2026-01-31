# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

envpick is a CLI tool for managing multiple environment variable configurations with interactive selection via fzf. It allows users to define different environment configurations in a TOML file and switch between them interactively.

## Build and Development Commands

```bash
# Build the binary
make compile          # Creates ./app binary
go build -o app .     # Alternative direct build

# Install to $GOPATH/bin
make install
go install .

# Run quality checks (tests, formatting, linting)
make quality

# Run tests only
go test ./...

# Format code
go fmt ./...

# Run linter
golangci-lint run
```

## Architecture

### Core Components

**Engine (`internal/core/engine.go`)**: Central orchestrator that coordinates config loading, state management, and namespace handling. Each Engine instance operates on a specific namespace.

**Config (`internal/config/config.go`)**: Manages `~/.envpick/config.toml` which stores environment variable definitions. Supports hierarchical namespaces using dot notation (e.g., `db.local`, `db.prod`). Variables starting with `_` are metadata (e.g., `_web_url`).

**State (`internal/config/state.go`)**: Manages `~/.envpick/state.toml` which tracks the currently active configuration per namespace. Handles migration from legacy single-config format.

**Selector (`internal/selector/fzf.go`)**: Wraps fzf for interactive configuration selection.

**Text (`internal/text/text.go`)**: Centralizes all user-facing text including command descriptions, error messages, prompts, and format strings. Uses typed structs (`CommandText`, `ErrorsText`, `MessagesText`, `FormatsText`, `PromptsText`) for better maintainability and potential i18n support.

**Commands (`cmd/`)**: Cobra-based CLI commands (use, env, edit, web) that interact with the Engine. All command text is sourced from the Text package.

### Namespace System

Configurations can be organized hierarchically:
- Default namespace: `dev`, `prod` (no prefix)
- Named namespaces: `db.local`, `db.prod`, `deploy.aws`

The `-n/--namespace` flag filters operations to a specific namespace. State is tracked separately per namespace.

### Key Design Patterns

- Config names are stored in two forms:
  - **Full name**: includes namespace prefix (e.g., `db.local`)
  - **Short name**: without namespace prefix (e.g., `local`)
- State stores short names, Engine handles conversion
- `ParseConfigName()` and `BuildConfigName()` handle name transformations
- All user-facing text is centralized in `internal/text/text.go` for consistency and maintainability
- Text package uses typed structs instead of anonymous structs for better code organization

## Testing

The project uses a comprehensive testing strategy:

### Unit Tests
Table-driven tests for individual components:
- `internal/config/config_test.go`: Config parsing and namespace logic
- `internal/config/state_test.go`: State management and migration
- `internal/core/engine_test.go`: Engine integration tests

### E2E Tests
End-to-end tests that verify complete user workflows:
- `test/e2e/e2e_test.go`: Full feature testing with isolated test environments
- Tests use in-process approach with mocked external dependencies (fzf, browser, editor)
- Each test runs in complete isolation with temporary HOME directories
- See `test/e2e/README.md` for detailed documentation

Run all tests:
```bash
go test ./...
make quality  # Runs tests + formatting + linting
```

## Dependencies

- `github.com/spf13/cobra`: CLI framework
- `github.com/BurntSushi/toml`: TOML parsing
- `fzf`: External dependency for interactive selection (must be in PATH)
