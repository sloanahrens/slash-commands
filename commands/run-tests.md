---
description: Run quality checks and tests for a repository
---

# Run Tests Command

Run quality checks and test suite for a repository.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, and commit rules.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Running tests for: <repo-name>"

### Step 2: Load Context

Per `_shared-repo-logic.md` → "Context Loading":
1. Read `~/.claude/CLAUDE.md` (global settings)
2. Read `<repo-path>/CLAUDE.md` (repo-specific guidance)

### Step 3: Run Quality Checks

Use `devbot check` for fast, auto-detected quality checks:

```bash
devbot check <repo-name>

# Or with prerequisite validation first:
devbot check <repo-name> --prereq
```

The `--prereq` flag validates tools, dependencies, and environment variables before running checks. This catches missing deps/env issues early instead of producing confusing test failures.

**Note:** If tests fail with dependency or environment errors, suggest running `devbot prereq <repo>` to diagnose.

This auto-detects the project stack (go, ts, nextjs, python, rust) and runs:
- **lint** and **typecheck** in parallel
- **build** and **test** sequentially

The command maps stacks to appropriate tools:
| Stack | Lint | Typecheck | Build | Test |
|-------|------|-----------|-------|------|
| nextjs/ts | npm run lint | npm run typecheck | npm run build | npm test |
| go | golangci-lint | - | go build | go test |
| python | ruff check | mypy | - | pytest |

**Override with config.yaml**: If repo has a `commands` block, use those instead.

**Subdirectory projects**: If stack not detected at root, check `work_dir` setting or common subdirs (go-api/, nextapp/, packages/*).

**Manual testing in monorepos**: For targeted testing in specific subdirectories, use `devbot exec`:

```bash
# Run tests in specific monorepo subproject
devbot exec fractals-nextjs/go-api go test ./...
devbot exec fractals-nextjs/nextapp npm test

# Uses work_dir from config automatically
devbot exec fractals-nextjs npm test  # runs in nextapp/
```

### Step 4: Report Results

```
| Check      | Status  | Details                    |
|------------|---------|----------------------------|
| Lint       | PASS    | No warnings or errors      |
| TypeScript | PASS    | No type errors             |
| Build      | PASS    | Production build successful|
| Tests      | PASS    | X passing, Y skipped       |
|------------|---------|----------------------------|
| TOTAL      | PASS    | All quality gates passed   |
```

### Step 5: Code Review (Optional)

**If all checks pass** and `--review` flag is passed (or user requests it):

Invoke `pr-review-toolkit:code-reviewer` to review recent changes:

```
"Launch code-reviewer agent to review unstaged changes for code quality"
```

This provides:
- Project guideline (CLAUDE.md) compliance
- Bug detection with confidence scoring
- Code quality issues

Report findings with severity levels and specific file:line references.

**Skip if**: No unstaged changes, or checks failed.

---

## Options

Parse flags from `$ARGUMENTS`:

| Flag | Effect |
|------|--------|
| `--only=<checks>` | Run only specified checks (comma-separated) |
| `--fix` | Auto-fix issues where possible (lint, format) |
| `--watch` | Run tests in watch mode (if supported) |
| `--review` | Run code-reviewer agent after tests pass |

**Note:** Use `devbot check <repo> --prereq` to validate prerequisites before checks.

**Check names:** `lint`, `typecheck`, `build`, `test`

Examples:
```bash
/run-tests pulumi --only=lint,typecheck   # Skip build and test
/run-tests my-app --only=test             # Just run tests
/run-tests frontend --fix                 # Auto-fix lint issues
```

---

## Error Handling

- If a check fails, analyze the error output

### Check for Known Patterns

Before debugging, search for relevant patterns/insights:

```bash
# Check for patterns tagged with testing or this repo
grep -l "tags:.*testing\|repos:.*<repo-name>" ~/.claude/patterns/*.md 2>/dev/null
grep -l "tags:.*testing\|repos:.*<repo-name>" ~/.claude/notes/insights/*.md 2>/dev/null
```

If a matching pattern exists, apply its solution first.

### Debug Process

- If `--fix` was passed, attempt auto-fix and re-run
- **If tests fail and cause is unclear**, invoke `superpowers:systematic-debugging` to investigate:
  - Gather evidence before hypothesizing
  - Form testable hypotheses
  - Verify fix actually resolves the issue
- Verify all fixes with a re-run
- If unable to fix automatically, report the issue with diagnosis

### Capture Learning

After resolving a non-trivial test failure, suggest:

```
Tests passing. If this failure was tricky:
  /capture-insight — Save this solution for future sessions
```

---

## Examples

```bash
/run-tests                          # Interactive selection
/run-tests pulumi                   # Fuzzy match → my-infra-pulumi
/run-tests my-app --only=test       # Just run tests
/run-tests frontend --fix           # Auto-fix lint issues
/run-tests cli --review             # Run tests then code review
```
