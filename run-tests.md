---
description: Run quality checks and tests for a repository
---

# Run Tests Command

Run quality checks and test suite for a repository.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, and commit rules.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Running tests for: <repo-name>"

### Step 2: Detect Language & Commands

Check `config.yaml` for repo-specific settings:
- `work_dir` → subdirectory to run commands from
- `language` → explicit language setting
- `commands` → custom command overrides

If `language` not specified, detect from files (see `_shared-repo-logic.md`).

### Step 3: Run Quality Checks

Run commands in order based on detected language (skip if command not available). See `_shared-repo-logic.md` for:
- Language detection rules
- Default commands per language (lint → typecheck → build → test)

If `commands` block exists in repo's config.yaml, use those instead of defaults.

Skip any command that's not available for the project.

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

---

## Options

Parse flags from `$ARGUMENTS`:

| Flag | Effect |
|------|--------|
| `--only=<checks>` | Run only specified checks (comma-separated) |
| `--fix` | Auto-fix issues where possible (lint, format) |
| `--watch` | Run tests in watch mode (if supported) |

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
```
