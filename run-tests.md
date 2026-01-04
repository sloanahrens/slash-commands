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
```

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
- If `--fix` was passed, attempt auto-fix and re-run
- **If tests fail and cause is unclear**, invoke `superpowers:systematic-debugging` to investigate:
  - Gather evidence before hypothesizing
  - Form testable hypotheses
  - Verify fix actually resolves the issue
- Verify all fixes with a re-run
- If unable to fix automatically, report the issue with diagnosis

---

## Examples

```bash
/run-tests                          # Interactive selection
/run-tests pulumi                   # Fuzzy match → my-infra-pulumi
/run-tests my-app --only=test       # Just run tests
/run-tests frontend --fix           # Auto-fix lint issues
/run-tests cli --review             # Run tests then code review
```
