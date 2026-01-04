# devbot

Fast parallel development workspace tools. Part of [slash-commands](../README.md).

## Installation

```bash
make install                    # Or run /setup-workspace in Claude Code
```

## Commands

### NAME Commands (take repo name)

#### path - Get Repository Path
```bash
devbot path <repo>              # Returns full path including work_dir
```

#### status - Parallel Git Status
```bash
devbot status                   # Dirty repos (~0.03s for 12 repos)
devbot status --all             # All repos
devbot status <repo>            # Single repo
```

#### diff - Git Diff Summary
```bash
devbot diff <repo>              # Staged/unstaged with line counts
devbot diff <repo> --full       # Include diff content
```

#### branch - Branch Info
```bash
devbot branch <repo>            # Branch, tracking, ahead/behind
```

#### remote - Remote Info
```bash
devbot remote <repo>            # Remote URLs and GitHub identifiers
```

#### find-repo - Find by GitHub ID
```bash
devbot find-repo owner/repo
devbot find-repo https://github.com/owner/repo/pull/123
```

#### check - Quality Checks
```bash
devbot check <repo>             # lint, typecheck, build, test
devbot check <repo> --only=lint # Specific checks
devbot check <repo> --fix       # Auto-fix
```

Auto-detects stack (go, ts, nextjs, python, rust).

#### last-commit - Commit Recency
```bash
devbot last-commit <repo>       # When was repo last committed
devbot last-commit <repo> FILE  # When was specific file last committed
```

#### todos - TODO/FIXME Scanning
```bash
devbot todos                    # All TODOs
devbot todos --type FIXME       # Filter by marker
devbot todos --count            # Counts only
devbot todos <repo>             # Single repo
```

#### config - Config File Discovery
```bash
devbot config                   # All config files
devbot config --type go         # Filter by type
devbot config <repo>            # Single repo
```

#### make - Makefile Targets
```bash
devbot make                     # All targets
devbot make <repo>              # Single repo
```

#### worktrees - Git Worktrees
```bash
devbot worktrees                # All worktrees
devbot worktrees <repo>         # Single repo
```

#### deps - Dependency Analysis
```bash
devbot deps                     # Shared dependencies (2+ repos)
devbot deps --all               # All by usage
devbot deps <repo>              # Single repo
```

#### run - Parallel Command Execution
```bash
devbot run -- git pull          # Run in all repos
devbot run -f myapp -- make     # Filter repos
devbot run -q -- git fetch      # Quiet mode
```

### PATH Commands (take filesystem path)

#### tree - Gitignore-Aware Tree
```bash
devbot tree <path>              # Directory tree
devbot tree -d 5                # Depth limit
```

#### stats - Code Metrics
```bash
devbot stats <path>             # File/dir analysis
devbot stats <path> -l go       # Filter by language
```

#### detect - Stack Detection
```bash
devbot detect <path>            # Outputs: go, ts, nextjs, etc.
```

## Development

```bash
make build       # Build binary
make test        # Run tests
make ci          # Full CI: fmt, vet, test, lint, build
make install     # Install to PATH
```

## Architecture

```
devbot/
├── cmd/devbot/main.go     # CLI entry (cobra)
├── internal/
│   ├── workspace/         # Repo discovery, parallel git status
│   ├── branch/            # Branch and tracking
│   ├── check/             # Quality checks
│   ├── config/            # Config discovery
│   ├── deps/              # Dependency analysis
│   ├── detect/            # Stack detection
│   ├── diff/              # Git diff
│   ├── lastcommit/        # Commit recency
│   ├── makefile/          # Makefile parsing
│   ├── output/            # Terminal rendering
│   ├── remote/            # Git remote parsing
│   ├── runner/            # Parallel execution
│   ├── stats/             # Code metrics
│   ├── todos/             # TODO scanning
│   ├── tree/              # Directory tree
│   └── worktrees/         # Worktree discovery
└── Makefile
```
