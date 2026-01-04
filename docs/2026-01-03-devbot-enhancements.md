# Devbot Enhancements

*Date: 2026-01-03*

## Summary

Added two new devbot commands and updated slash commands to use them for faster execution.

## New Commands

### `devbot diff <repo>`

**Purpose**: Get git diff summary in a single call instead of multiple git commands.

**Output**:
```
slash-commands/
────────────────────────────────────────────────────────────
  Branch:   master
  Staged:   2 files (+45, -12)
  Unstaged: 1 file (+3, -1)

  Staged:
    M  src/config.ts (+32, -8)
    A  src/utils/helper.ts (+13, -4)

  Unstaged:
    M  README.md (+3, -1)

(0.02s)
```

**Speed**: ~0.02s (replaces 3 separate git commands)

**Files added**:
- `devbot/internal/diff/diff.go`

### `devbot check <repo>`

**Purpose**: Auto-detect project stack and run lint/typecheck/build/test.

**Features**:
- **Multi-app support**: Discovers sub-applications in `go-api/`, `nextapp/`, `packages/*`, etc.
- Runs lint and typecheck in **parallel** (per sub-app)
- Runs build and test **sequentially** after
- Supports `--only=lint,test` to run specific checks
- Supports `--fix` to auto-fix lint issues
- Exits with code 1 on first failure

**Example output for multi-app repo**:
```
mango/ (go-api:go | nextapp:ts,nextjs)
────────────────────────────────────────────────────────────

  go-api/
    build        ✓ PASS   6.4s

  nextapp/
    build        ✓ PASS   4.1s

  ────────────────────────────────────────
  Total: PASS (10.5s)
```

**Stack detection**:
| Stack | Lint | Typecheck | Build | Test |
|-------|------|-----------|-------|------|
| nextjs/ts | npm run lint | npm run typecheck | npm run build | npm test |
| go | golangci-lint | - | go build | go test |
| python | ruff check | mypy | - | pytest |
| rust | cargo clippy | - | cargo build | cargo test |

**Files added**:
- `devbot/internal/check/check.go`

## Slash Command Updates

### `/yes-commit`

**Before**: 3 steps with separate commands
```
Step 2: devbot status <repo>
Step 3: git diff --stat
```

**After**: 1 step with devbot diff
```
Step 2: devbot diff <repo>
```

### `/run-tests`

**Before**: Manual language detection and sequential command execution

**After**: Uses `devbot check <repo>` for auto-detected parallel execution

### `_shared-repo-logic.md`

Added `devbot diff` and `devbot check` to the devbot CLI reference table.

### `README.md`

Added documentation for both new commands.

## Performance Impact

| Operation | Before | After | Speedup |
|-----------|--------|-------|---------|
| `/yes-commit` status check | ~0.05s (3 git calls) | ~0.02s (1 devbot call) | 2.5x |
| `/run-tests` (full suite) | ~10s (sequential) | ~6s (parallel lint/typecheck) | ~40% |

## Files Changed

```
devbot/cmd/devbot/main.go          # Added diff and check commands
devbot/internal/diff/diff.go       # New package
devbot/internal/check/check.go     # New package
yes-commit.md                       # Updated to use devbot diff
run-tests.md                        # Updated to use devbot check
_shared-repo-logic.md              # Added new commands to reference
README.md                          # Added documentation
```
