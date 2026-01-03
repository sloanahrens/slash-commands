# CLAUDE.md

This file provides guidance to Claude Code when working with devbot.

## Project Overview

**devbot** is a Go CLI for accelerating development operations through parallelization.

**Stack**: Go 1.23
**Purpose**: Fast, parallel workspace operations across ~/code/
**Location**: `~/code/slash-commands/devbot/` (part of slash-commands repo)

## Quick Reference

```bash
# Build and install (from this directory)
make install

# Or use the slash command
/install-devbot

# Test
make test
```

## Commands

### status - Parallel Git Status (0.03-0.05s for 12 repos)

```bash
devbot status           # Show dirty repos (clean count summarized)
devbot status --all     # Show all repos
devbot status --dirty   # Only dirty repos
devbot status <repo>    # Single repo details
```

### run - Parallel Command Execution

```bash
devbot run -- git pull              # Pull all repos in parallel
devbot run -- npm install           # Install deps in all repos
devbot run -f mango -- make build   # Filter to repos matching "mango"
devbot run -q -- git fetch          # Quiet mode (suppress empty output)
```

### deps - Dependency Analysis

```bash
devbot deps             # Show shared dependencies (2+ repos)
devbot deps --all       # Show all dependencies by usage
devbot deps --count     # Show dependency counts per repo
devbot deps <repo>      # Analyze single repo
```

### tree - Gitignore-Aware Tree

```bash
devbot tree                 # Current directory
devbot tree ~/code/mango    # Specific path
devbot tree -d 5            # Depth limit (default: 3)
devbot tree --hidden        # Include hidden files
```

### detect - Project Stack Detection

```bash
devbot detect               # Current directory
devbot detect ~/code/mango  # Specific path
# Output: Detected: go, ts, nextjs
```

### todos - Parallel TODO/FIXME Scanning

```bash
devbot todos                    # All TODOs across workspace
devbot todos --type FIXME       # Filter by marker type
devbot todos --limit 20         # Limit results
devbot todos <repo>             # Single repo
```

Scans for: TODO, FIXME, HACK, XXX, BUG in .go, .ts, .tsx, .js, .jsx, .py, .md, .yaml, .yml files.

### config - Config File Discovery

```bash
devbot config                   # All config files by type
devbot config --type go         # Filter by config type
devbot config <repo>            # Single repo
```

Detects: package.json, go.mod, Dockerfile, Makefile, CLAUDE.md, .env, pyproject.toml, etc.

### make - Makefile Target Analysis

```bash
devbot make                     # All targets grouped by category
devbot make --category test     # Filter by category
devbot make <repo>              # Single repo
```

Categories: setup, dev, database, test, build, clean, other

### worktrees - Git Worktree Discovery

```bash
devbot worktrees                # All worktrees across repos
devbot worktrees <repo>         # Single repo
```

Scans: .trees/, worktrees/, .worktrees/ directories

## Architecture

```
cmd/devbot/main.go       # CLI entry point (cobra)
internal/
├── workspace/           # Repo discovery and parallel git status
├── runner/              # Parallel command execution
├── deps/                # Dependency analysis (package.json, go.mod)
├── tree/                # Gitignore-aware directory tree
├── detect/              # Project stack detection
├── todos/               # Parallel TODO/FIXME scanning
├── config/              # Config file discovery
├── makefile/            # Makefile target parsing
├── worktrees/           # Git worktree discovery
└── output/              # Terminal rendering
```

## Design Decisions

1. **Discovery**: Scans `~/code/` for directories with `.git` (not recursive)
2. **Parallelization**: One goroutine per repo, results collected via channels
3. **Stack detection**: Checks root + common subdirs (go-api/, nextapp/, packages/*, apps/*)
4. **Tree filtering**: Built-in ignores + .gitignore parsing

## Adding New Commands

1. Add command definition in `cmd/devbot/main.go`
2. Create internal package in `internal/<feature>/`
3. Use parallel pattern from `workspace/status.go` or `runner/run.go`

## Testing

```bash
go test ./internal/... -v           # Unit tests
go test ./... -v -race              # Race detection
```
