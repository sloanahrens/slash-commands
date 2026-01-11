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

## New Machine Setup

After cloning this repo on a new machine, run these steps to set up the full environment:

### 1. Install devbot

```bash
cd devbot
make install
```

### 2. Create notes directory structure

```bash
mkdir -p ~/.claude/notes/hindsight ~/.claude/notes/sessions
```

### 3. Copy hookify rules

```bash
cp docs/config/hookify.*.local.md ~/.claude/
```

### 4. Symlink slash commands

```bash
# Create commands directory if needed
mkdir -p ~/.claude/commands

# Symlink all slash commands
for f in *.md; do
  [[ "$f" != "README.md" && "$f" != "_"* ]] && ln -sf "$(pwd)/$f" ~/.claude/commands/
done

# Also symlink config.yaml
ln -sf "$(pwd)/config.yaml" ~/.claude/commands/
```

### 5. Copy root CLAUDE.md

```bash
cp docs/root-claude/CLAUDE.md ~/.claude/CLAUDE.md
# Edit as needed for your specific workspace path
```

### 6. Update config.yaml

```bash
cp config.yaml.example config.yaml
# Edit to match your workspace structure
```

Or run `/setup-workspace` which automates most of the above.

## docs/ Directory Contents

| Path | Purpose |
|------|---------|
| `docs/config/` | Hookify rules to copy to `~/.claude/` |
| `docs/root-claude/` | Example global CLAUDE.md |
| `docs/patterns/` | Versioned knowledge patterns (agent scaffolding) |
| `docs/templates/` | Subagent templates for `/improve` command |
| `docs/plans/` | Design documents |

## Agent Scaffolding Commands

These commands implement Confucius-inspired agent scaffolding for cross-session learning:

| Command | Purpose |
|---------|---------|
| `/prime <repo>` | Load relevant patterns and notes before work |
| `/capture-hindsight` | Save failure/discovery as a note |
| `/promote-pattern` | Graduate useful note to versioned pattern |
| `/age-notes` | Review and clean up old notes |
| `/improve <repo> <task>` | Meta-agent loop with parallel subagents |
