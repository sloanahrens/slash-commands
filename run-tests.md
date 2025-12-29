---
description: Run quality checks and tests (for specified repo, or prompts for selection)
---

# Run Tests Command

Run quality checks and test suite for a repository.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, and commit rules.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Running tests for: <repo-name>"

### Step 2: Detect Capabilities

Read `<repo>/package.json` and check for scripts:

| Script | Capability |
|--------|------------|
| `test` | Run tests |
| `lint` | Run linting |
| `type-check` or `typecheck` | Type checking |
| `build` | Production build |

Check for overrides in `.env.local`:
- `<PREFIX>_TEST_CMD` → custom test command
- `<PREFIX>_WORK_DIR` → subdirectory to run from (e.g., `nextapp`)

### Step 3: Run Quality Checks

Run in order (if available):

```bash
cd <repo-path> && npm run lint
cd <repo-path> && npm run type-check
cd <repo-path> && npm run build
cd <repo-path> && npm test
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
