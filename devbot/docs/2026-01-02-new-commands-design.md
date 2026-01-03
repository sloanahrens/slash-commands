# Devbot New Commands Design

## Overview

Four new commands to accelerate common development operations beyond git:

| Command | Purpose | Expected Speed |
|---------|---------|----------------|
| `todos` | Parallel TODO/FIXME scanning | 0.05-0.10s for 12 repos |
| `config` | Config file discovery and indexing | 0.02-0.05s |
| `make` | Makefile parsing and target analysis | 0.01-0.03s |
| `worktrees` | Fast worktree discovery with branch info | 0.03-0.05s |

## Command 1: `devbot todos`

Parallel TODO/FIXME scanner across all repos.

### Usage

```bash
devbot todos              # Scan all repos, show grouped by repo
devbot todos --count      # Just show counts per repo
devbot todos <repo>       # Single repo, full detail
devbot todos --type=FIXME # Filter to specific markers
```

### Markers Scanned

`TODO`, `FIXME`, `HACK`, `XXX`, `BUG`

### File Types

`.go`, `.ts`, `.tsx`, `.js`, `.jsx`, `.py`, `.md`

### Exclusions

- Respects `.gitignore` (reuse tree package logic)
- Skips: `node_modules`, `dist`, `build`, `.git`, `vendor`

### Output Format

```
devbot/
  internal/runner/run.go:45    TODO: add timeout support
  cmd/devbot/main.go:112       FIXME: handle empty filter

mango/
  go-api/handlers/query.go:89  TODO: validate input bounds

(12 repos, 47 items, 0.08s)
```

### Implementation

- Package: `internal/todos/`
- Reuse `workspace.Discover()` for repo list
- Reuse gitignore logic from `tree` package
- Parallel scan per repo, parallel file reads within repo

---

## Command 2: `devbot config`

Config file discovery and indexing.

### Usage

```bash
devbot config                    # List all config files across workspace
devbot config --type=package     # Only package.json files
devbot config <repo>             # Config files in one repo
devbot config --has=typescript   # Repos with tsconfig.json
```

### Config Files Tracked

| Type | Files |
|------|-------|
| Node | `package.json`, `tsconfig.json`, `pnpm-workspace.yaml` |
| Go | `go.mod`, `go.sum` |
| Python | `pyproject.toml`, `requirements.txt`, `setup.py` |
| Infra | `Makefile`, `Dockerfile`, `docker-compose.yml` |
| IaC | `Pulumi.yaml`, `serverless.yml` |
| CI | `.github/workflows/*.yml`, `.gitlab-ci.yml` |
| Config | `config.yaml`, `.env.example`, `CLAUDE.md` |

### Output Format

```
devbot/
  go.mod, Makefile, CLAUDE.md

mango/
  go-api/go.mod, nextapp/package.json, Makefile

(12 repos, 45 config files, 0.02s)
```

### Implementation

- Package: `internal/config/`
- Parallel repo scan
- Check known paths (not recursive grep)
- Optional: cache to `~/.cache/devbot/config-index.json`

---

## Command 3: `devbot make`

Makefile parsing and target analysis.

### Usage

```bash
devbot make                 # List repos with Makefiles and target counts
devbot make <repo>          # Parse and categorize targets for one repo
devbot make --targets       # Show all targets across all repos
```

### Target Categories

| Category | Pattern matches |
|----------|-----------------|
| Setup | `setup`, `install`, `init`, `bootstrap` |
| Dev | `dev`, `run`, `start`, `serve`, `watch` |
| Database | `db`, `migrate`, `docker`, `postgres` |
| Test | `test`, `lint`, `check`, `typecheck`, `fmt` |
| Build | `build`, `compile`, `dist`, `release` |
| Clean | `clean`, `reset`, `purge` |

### Output Format (single repo)

```
mango/Makefile - 12 targets

Setup:     setup, install
Dev:       dev, dev-app
Database:  db-up, db-down, migrate
Test:      test, lint, fmt
Build:     build, build-go, build-next
Clean:     clean
```

### Output Format (all repos)

```
Makefiles found:
  devbot         8 targets  (setup, test, build, ...)
  mango         12 targets  (setup, dev, db-up, ...)

(12 repos, 3 with Makefiles, 0.03s)
```

### Implementation

- Package: `internal/makefile/`
- Parse targets: regex `^[a-zA-Z_-]+:`
- Parse `.PHONY` declarations
- Extract comments above targets as descriptions
- Categorize by pattern matching

---

## Command 4: `devbot worktrees`

Fast worktree discovery with branch info.

### Usage

```bash
devbot worktrees              # List all worktrees across repos
devbot worktrees <repo>       # Worktrees for one repo
devbot worktrees --branches   # Group by branch name
```

### What It Discovers

- Scans each repo for `.trees/` directories
- Extracts branch name via `git branch --show-current`
- Gets dirty status (reuses workspace package)

### Output Format

```
my-workspace/
  .trees/feature-new-auth     → feature/new-auth (clean)
  .trees/fix-mcp-timeout      → fix/mcp-timeout (2 modified)

mango/
  .trees/duckdb-optimization  → feature/duckdb-opt (clean)

(12 repos, 3 with worktrees, 5 total, 0.04s)
```

### Implementation

- Package: `internal/worktrees/`
- Parallel scan for `.trees/` directories
- Parallel git calls for branch names
- Reuse `workspace.getRepoStatus()` for dirty status

---

## Architecture

All commands follow existing devbot patterns:

```
cmd/devbot/main.go           # Add new command definitions
internal/
├── todos/todos.go           # TODO scanner
├── config/config.go         # Config file indexer
├── makefile/makefile.go     # Makefile parser
├── worktrees/worktrees.go   # Worktree discovery
├── workspace/               # Existing - repo discovery
├── tree/                    # Existing - gitignore logic (reuse)
└── output/                  # Existing - rendering
```

### Shared Patterns

1. **Parallel execution**: goroutines per repo
2. **Channel collection**: results sent to buffered channel
3. **Sorted output**: alphabetical by repo name
4. **Timing**: elapsed time reported at end
5. **Gitignore**: reuse `tree.shouldIgnore()` logic

---

## Implementation Order

1. `todos` - highest value, reuses existing gitignore logic
2. `config` - straightforward file existence checks
3. `make` - text parsing, medium complexity
4. `worktrees` - git operations, reuses workspace patterns
