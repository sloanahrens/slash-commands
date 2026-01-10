# devbot prereq Command Design

**Date:** 2026-01-09
**Status:** Approved

## Purpose

Validate prerequisites before starting work on a repo. Catches missing tools, dependencies, and environment variables upfront instead of failing mid-task.

## Usage

```bash
devbot prereq <repo>[/subdir]

# Examples
devbot prereq atap-automation2      # Uses work_dir from config (nextapp)
devbot prereq mango/go-api          # Explicit subdir
devbot prereq slash-commands/devbot # Subdir within repo
```

Uses same directory resolution as `devbot exec`.

## Output Format

Compact table matching devbot style:

```
  ~/code/mono-claude/atap-automation2/nextapp (ts, nextjs)
──────────────────────────────────────────────────────────
  ✓ node       v22.1.0
  ✓ npm        10.7.0
  ✓ deps       node_modules present
  ✗ env        GOOGLE_CLOUD_PROJECT missing

  1 issue(s) found
```

**Exit codes:**
- 0: All checks pass
- 1: Any check failed

## Checks Performed

### 1. Tool Checks

Per detected stack, verify required binaries exist:

| Stack | Required Tools |
|-------|----------------|
| go | `go` |
| ts/js/nextjs | `node`, `npm` |
| python | `python3`, `pip3` |
| rust | `cargo` |

Reports tool version if found, "not found" if missing.

### 2. Dependency Checks

Per stack, verify deps are installed:

| Stack | Check |
|-------|-------|
| ts/js/nextjs | `node_modules` directory exists |
| go | `go.sum` file exists |
| python | `venv` or `.venv` directory exists |
| rust | `Cargo.lock` file exists |

### 3. Environment Checks

**Convention:**
- `.env.local.example` must exist in work directory (warn if missing)
- `.env.local` should have all vars defined in example

**Logic:**
1. If `.env.local.example` missing → Warn "`.env.local.example` missing"
2. Parse example file for var names (lines matching `^[A-Z][A-Z0-9_]*=`)
3. For each var, check if present in `.env.local` OR in actual environment
4. Report missing vars by name

## Package Structure

```
internal/prereq/
├── prereq.go      # Main Run() function, orchestrates checks
├── prereq_test.go # Tests
├── tools.go       # Tool binary checking
└── env.go         # Env file comparison
```

## Core Types

```go
type Status int
const (
    Pass Status = iota
    Fail
    Warn
)

type Check struct {
    Name   string // "node", "deps", "env"
    Status Status
    Detail string // "v22.1.0" or "GOOGLE_CLOUD_PROJECT missing"
}

type Result struct {
    Path   string
    Stack  []string
    Checks []Check
}

func Run(repoPath string, workDir string) (*Result, error)
```

## Implementation Notes

- Reuses `detect.ProjectStack()` for stack detection
- Reuses `exec` package's directory resolution logic
- Tool versions via `<tool> --version`, parse first line
- Env parsing: simple line-by-line, regex `^[A-Z][A-Z0-9_]*=`

## Estimated Size

~200 lines across the package.
