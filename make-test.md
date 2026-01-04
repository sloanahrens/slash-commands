---
description: Test Makefile targets for a repository
---

# Make Test Command

Run all Makefile targets for a repository in clean-build order.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, and commit rules.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Testing Makefile for: <repo-name>"

### Step 2: Parse Makefile with devbot

Use devbot for instant Makefile analysis:

```bash
devbot make <repo-name>
```

This parses in ~0.01s and returns:
- All targets with `.PHONY` status
- Category classification (setup, dev, database, test, build, clean, other)
- Comments/descriptions from lines above targets

If no Makefile found, report error and exit.

### Step 3: Display Target Summary

Show a brief summary of discovered targets:

```
Makefile for <repo-name>: 13 targets
  clean: clean
  setup: install
  build: build
  test:  lint, typecheck, check
  dev:   dev, start (smoke test)
```

### Step 4: Execute All Targets (Default Behavior)

**Run all targets automatically in clean-build order.** No prompting unless `--interactive` flag is passed.

**Execution order:**

1. **Setup/Install** - Restore dependencies: `install`, `setup`, `init`
2. **Database/Docker** - Infrastructure: `docker-build`, `db-up`, `migrate`
3. **Build** - Compile: `build`, `compile`, `dist`
4. **Test/Lint** - Verify: `lint`, `typecheck`, `format`, `check`, `test`
5. **Dev/Blocking** - Smoke test: `dev`, `run`, `start`, `serve` (with timeout)

**Note:** Clean targets (`clean`, `reset`, `docker-down`, etc.) are **skipped by default** for faster incremental builds. Pass `clean` argument to run them first.

**Special handling:**

| Target Type | Behavior |
|-------------|----------|
| Blocking (`dev`, `run`, `start`, `serve`) | Run with 5s timeout to verify startup, then kill |
| Docker-dependent (`db-up`, `docker-*`) | Check if Docker is running first; if not, stop and ask user to start it, then continue when ready |
| Long-running (`build`, `test`) | Use extended timeout (5 min) |
| Format (`format`, `fmt`) | Run but don't fail pipeline (may modify files) |

**Docker check:** Before running any Docker-dependent target, run `docker info` to verify Docker is running. If it fails:
1. Display: "Docker is not running. Please start Docker Desktop and press Enter to continue."
2. Wait for user confirmation
3. Re-check Docker status before proceeding

### Step 5: Execute and Time Targets

For each target:

1. **Announce**: "Testing: make <target>"
2. **Start timer**: Record start time
3. **Execute**: Run with appropriate timeout
4. **Stop timer**: Record elapsed time
5. **Capture output**: Show last 20 lines on success, full output on failure
6. **Report result**: PASS / FAIL / TIMEOUT / SKIPPED with elapsed time

**Use `time` command** to measure each target:
```bash
time make -C <repo-path> <target>
```

**Timeouts by category:**
- Quick (lint, typecheck, format): 60s
- Build/test: 5 min
- Blocking (dev, start): 5s (smoke test only)
- Docker: 30s

### Step 6: Report Results with Timing Summary

Display a table showing all targets with their execution times, sorted by execution order:

```
Makefile Test Results for <repo-name>
=====================================

| Target     | Status  | Time    | Tests   | Coverage | Notes                    |
|------------|---------|---------|---------|----------|--------------------------|
| clean      | PASS    |   0.3s  |    -    |    -     | Build artifacts removed  |
| install    | PASS    |  12.3s  |    -    |    -     | Dependencies installed   |
| db-up      | PASS    |   3.2s  |    -    |    -     | PostgreSQL ready         |
| build      | PASS    |   5.2s  |    -    |    -     | Build successful         |
| lint       | PASS    |   2.4s  |    -    |    -     | No issues found          |
| typecheck  | PASS    |   1.8s  |    -    |    -     | Types valid              |
| test       | PASS    |   8.7s  |   36    |   78%    | All tests passed         |
| dev        | PASS    |   5.0s  |    -    |    -     | Server started (smoke)   |
|------------|---------|---------|---------|----------|--------------------------|
| TOTAL      | 8/8     |  38.9s  |   36    |   78%    | All targets passed       |

Slowest targets:
  1. install    12.3s  (31%)
  2. test        8.7s  (22%)
  3. build       5.2s  (13%)

Issues Found:
  (none)
```

**Extracting test counts and coverage:**
- Parse test output for patterns like "36 passed", "36 tests", "Tests: 36"
- Parse coverage output for patterns like "78%", "Coverage: 78%", "78% covered"
- Common test frameworks:
  - Jest: "Tests: X passed" and "All files | XX% |"
  - Vitest: "X passed" and "Coverage: XX%"
  - pytest: "X passed" and "TOTAL XX%"
  - Go: "ok" count and "coverage: XX%"

**If test counts or coverage not available:**

After the results table, display:

```
Missing metrics:
  - Test count: Not found in test output
  - Coverage: No coverage report configured

💡 Run `/super <repo-name>` to brainstorm adding test counts and coverage reporting.
```

**Timing details:**
- Times are wall-clock elapsed time in seconds
- Right-align times for easy comparison
- Show top 3 slowest targets with percentage of total time
- If any target exceeds 30s, flag it as potentially worth optimizing

---

## Options

Parse flags from `$ARGUMENTS`:

| Flag | Effect |
|------|--------|
| `--dry-run` | Analyze only, don't execute |
| `--interactive` | Prompt for target selection instead of running all |
| `--quick` | Only test quick targets (lint, typecheck, format) |
| `--skip-clean` | Skip destructive targets, start from install |
| `--skip-docker` | Explicitly skip Docker targets (don't prompt to start Docker) |

Examples:
```bash
/make-test my-app                   # Run all targets in clean-build order
/make-test my-app --dry-run         # Just analyze Makefile
/make-test frontend --quick         # Only lint/typecheck
/make-test api --skip-docker        # Skip Docker targets (no Docker running)
/make-test infra --interactive      # Prompt for target selection
```

---

## Error Handling

- If a target fails, capture full error output
- Analyze common failure patterns:
  - Missing dependencies → suggest `make install` or `pnpm install`
  - Docker not running → suggest starting Docker
  - Port in use → identify blocking process
  - Missing env vars → check for `.env.example`
- Offer to fix issues and retry
- If unable to diagnose, report with full error context

---

## Makefile Improvements

After testing, suggest improvements if applicable:

- Missing `help` target → offer to add one
- No `.PHONY` declarations → suggest adding them
- Undocumented targets → suggest adding comments
- Missing common targets → suggest: `lint`, `typecheck`, `clean`
- Circular dependencies → warn about them
- Targets that could be parallelized → suggest `make -j`

Ask: "Would you like me to apply any of these improvements?"

---

## Examples

```bash
/make-test                          # Select repo, run all targets
/make-test my-app                   # Run all targets for my-nextjs-app
/make-test infra --dry-run          # Analyze without running
/make-test api --quick              # Only lint/typecheck
/make-test frontend --interactive   # Choose which targets to run
```
