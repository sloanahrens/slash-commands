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

Run commands in order based on language (skip if command not available):

**TypeScript/JavaScript** (package.json):
```bash
npm run lint
npx tsc --noEmit      # or npm run type-check
npm run build
npm test
```

**Go** (go.mod):
```bash
golangci-lint run     # if available
go build ./...
go test ./...
```

**Python** (pyproject.toml):
```bash
ruff check .          # or flake8
mypy .                # if available
pytest
```

**Rust** (Cargo.toml):
```bash
cargo clippy
cargo build
cargo test
```

If `commands` block exists in config, use those instead.

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

## Error Handling

- If a check fails, analyze the error output
- Apply fixes directly to source files
- Verify fixes with a re-run
- If unable to fix automatically, report the issue

---

## Examples

```bash
/run-tests              # Interactive selection
/run-tests pulumi       # Fuzzy match → devops-gcp-pulumi
/run-tests atap         # Fuzzy match → atap-automation2
```
