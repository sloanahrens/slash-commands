# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Slash Commands is the **meta-tooling** for the mono-claude workspace. It provides:
- **Slash commands** (`.md` files at repo root) - Claude Code skills for workspace operations
- **devbot CLI** (`devbot/`) - Go binary for fast parallel repo operations

## Development Commands

```bash
# devbot CLI (from devbot/ directory)
cd devbot
make build          # Build binary
make test           # Run tests with verbose output
make test-race      # Tests with race detector
make test-cover     # Tests with coverage (outputs coverage.out)
make lint           # golangci-lint (requires: brew install golangci-lint)
make ci             # Full CI: fmt, vet, test-race, test-cover
make install        # Install to $GOPATH/bin

# Single test
go test ./internal/workspace/... -v -run TestStatusParsing

# Quick iteration
go build -o devbot ./cmd/devbot && ./devbot status
```

## Architecture

```
slash-commands/
├── *.md                    # Slash commands (Claude Code skills)
├── _shared-repo-logic.md   # Common patterns imported by other commands
├── config.yaml             # Workspace config (gitignored)
├── config.yaml.example     # Template
├── docs/                   # Global config templates (→ ~/.claude/)
└── devbot/
    ├── cmd/devbot/         # CLI entry (Cobra)
    └── internal/           # Command implementations
        ├── workspace/      # Repo discovery, parallel git status
        ├── check/          # Stack-aware quality checks
        ├── detect/         # Stack detection (go, ts, nextjs, python)
        └── ...             # One package per command
```

## Key Patterns

### Slash Commands
- Each `.md` file is a skill that Claude Code can invoke
- `_shared-repo-logic.md` contains patterns included by other commands via `{{include}}`
- Commands use `devbot` for git operations instead of raw git commands

### devbot Internal Organization
- One package per command under `internal/`
- `internal/workspace/` handles config loading and repo discovery
- `internal/output/` provides consistent terminal rendering
- Tests use `testdata/` fixtures where needed

## Testing Changes

After modifying devbot:
1. `make test` in `devbot/` directory
2. Verify with actual repo: `./devbot status`, `./devbot check <repo>`
3. If changing slash commands, test by invoking the command in Claude Code
