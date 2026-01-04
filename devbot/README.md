# devbot

Fast parallel development workspace tools written in Go.

Part of [slash-commands](../README.md) - Claude Code slash commands for multi-repo workspaces.

## Installation

```bash
make install
# Or from parent: /install-devbot
```

## Commands

### status - Parallel Git Status

```bash
devbot status              # Show dirty repos (clean count summarized)
devbot status --all        # Show all repos
devbot status --dirty      # Only dirty repos
devbot status <repo>       # Single repo details
```

~0.03s for 12 repos.

### diff - Git Diff Summary

```bash
devbot diff <repo>         # Staged/unstaged files with line counts
devbot diff <repo> --full  # Include full diff content
```

Shows branch, staged files, unstaged files with +/- counts.

### branch - Branch and Tracking Info

```bash
devbot branch <repo>       # Branch, tracking, ahead/behind, commits to push
```

### remote - Git Remote Info

```bash
devbot remote <repo>       # Remote URLs and GitHub identifiers
```

### find-repo - Find Repo by GitHub ID

```bash
devbot find-repo owner/repo                              # By identifier
devbot find-repo https://github.com/owner/repo/pull/123  # By URL
```

### check - Auto-Detected Quality Checks

```bash
devbot check <repo>              # Run all checks (lint, typecheck, build, test)
devbot check <repo> --only=lint  # Run specific checks (comma-separated)
devbot check <repo> --fix        # Auto-fix where possible
```

Auto-detects stack (go, ts, nextjs, python, rust) and runs appropriate commands:
- Lint and typecheck run in parallel
- Build and test run sequentially
- Exits with code 1 on first failure

### run - Parallel Command Execution

```bash
devbot run -- git pull             # Pull all repos in parallel
devbot run -- npm install          # Install deps in all repos
devbot run -f myapp -- make build  # Filter to repos matching "myapp"
devbot run -q -- git fetch         # Quiet mode (suppress empty output)
```

### todos - Parallel TODO/FIXME Scanning

```bash
devbot todos                # All TODOs across workspace
devbot todos --type FIXME   # Filter by marker type (TODO, FIXME, HACK, XXX, BUG)
devbot todos --count        # Show counts only
devbot todos <repo>         # Single repo
```

Scans: `.go`, `.ts`, `.tsx`, `.js`, `.jsx`, `.py`, `.md`, `.yaml`, `.yml` files.

### stats - Code Metrics and Complexity

```bash
devbot stats <path>         # Analyze file or directory
devbot stats <path> -l go   # Filter by language
```

Reports: files, lines (code/comments/blank), functions, complexity flags.

### detect - Project Stack Detection

```bash
devbot detect               # Current directory
devbot detect <path>        # Specific path
```

Outputs detected stacks (e.g., `go`, `ts`, `nextjs`, `python`, `rust`).

### deps - Dependency Analysis

```bash
devbot deps               # Show shared dependencies (2+ repos)
devbot deps --all         # Show all dependencies by usage
devbot deps --count       # Show counts only
devbot deps <repo>        # Analyze single repo
```

### tree - Gitignore-Aware Tree

```bash
devbot tree               # Current directory
devbot tree <path>        # Specific path
devbot tree -d 5          # Depth limit (default: 3)
devbot tree --hidden      # Show hidden files
```

### config - Config File Discovery

```bash
devbot config              # All config files by type
devbot config --type go    # Filter by type (node, go, python, infra, iac, ci, config)
devbot config --has node   # Show only repos with this config type
devbot config <repo>       # Single repo
```

### make - Makefile Target Analysis

```bash
devbot make                # All targets grouped by category
devbot make --targets      # Show all targets across all repos
devbot make <repo>         # Single repo
```

### worktrees - Git Worktree Discovery

```bash
devbot worktrees           # All worktrees across repos
devbot worktrees <repo>    # Single repo
```

## Architecture

```
devbot/
├── cmd/devbot/main.go     # CLI entry point (cobra)
├── internal/
│   ├── workspace/         # Repo discovery and parallel git status
│   ├── runner/            # Parallel command execution
│   ├── branch/            # Branch and tracking info
│   ├── check/             # Auto-detected quality checks
│   ├── config/            # Config file discovery
│   ├── deps/              # Dependency analysis
│   ├── detect/            # Project stack detection
│   ├── diff/              # Git diff summary
│   ├── makefile/          # Makefile target parsing
│   ├── output/            # Terminal rendering
│   ├── remote/            # Git remote and GitHub ID parsing
│   ├── stats/             # Code metrics and complexity
│   ├── todos/             # Parallel TODO/FIXME scanning
│   ├── tree/              # Gitignore-aware directory tree
│   └── worktrees/         # Git worktree discovery
├── testdata/              # Test fixtures
├── Makefile               # Build targets
├── go.mod
└── go.sum
```

## Development

```bash
make build       # Build binary
make test        # Run tests
make test-race   # Run with race detector
make test-cover  # Run with coverage
make lint        # Run golangci-lint
make ci          # Run all checks
```

## Codebase Metrics

- **Files:** 31 Go source files
- **Lines:** 7,564 total (5,956 code, 412 comments, 1,196 blank)
- **Functions:** 247 (average 9 lines)
- **Test coverage:** 10 packages with tests
