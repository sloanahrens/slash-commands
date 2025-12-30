---
description: Test Makefile targets for a repository
---

# Make Test Command

Parse and interactively test Makefile targets for a repository.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, and commit rules.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Testing Makefile for: <repo-name>"

### Step 2: Find and Parse Makefile

1. Check for Makefile in repo root (or `work_dir` if configured)
2. If no Makefile found, report error and exit
3. Parse Makefile to extract:
   - All target names (lines matching `^target-name:`)
   - `.PHONY` declarations
   - Target descriptions (comments above targets)
   - Target dependencies

### Step 3: Analyze Targets

Categorize targets by type:

| Category | Detection | Examples |
|----------|-----------|----------|
| **Setup** | Contains `setup`, `install`, `init` | `make setup`, `make install` |
| **Development** | Contains `dev`, `run`, `start`, `serve` | `make dev`, `make run` |
| **Database** | Contains `db`, `migrate`, `docker` | `make db-up`, `make migrate` |
| **Testing** | Contains `test`, `lint`, `check` | `make test`, `make lint` |
| **Build** | Contains `build`, `compile`, `dist` | `make build`, `make clean` |
| **Other** | Everything else | `make help`, custom targets |

Present analysis:

```
Makefile Analysis for <repo-name>
=================================

Found X targets in Y categories:

Setup (run first):
  - setup: First time setup (install, DB, migrations)
  - install: Install dependencies only

Development:
  - dev: Start API server on port 4000
  - dev-app: Start test app on port 5173

Database:
  - db-up: Start PostgreSQL container
  - db-down: Stop PostgreSQL container

Testing:
  - test: Run unit tests
  - lint: Run linter

Build:
  - build: Build all packages
  - clean: Remove build artifacts
```

### Step 4: Interactive Testing

Ask user which targets to test:

```
Which targets would you like to test?

1. All targets (in logical order)
2. Setup targets only
3. Test/Lint targets only
4. Select specific targets
5. Skip testing (just analyze)

Enter choice:
```

**For option 1 (All targets):**
- Run in logical order: setup → db → build → test → dev (skip dev, it blocks)
- Skip targets that would block (dev, run, start, serve, watch)
- Ask before running destructive targets (clean, reset, down)

**For option 4 (Select specific):**
- Show numbered list of all targets
- Allow comma-separated selection: "1,3,5" or "test,lint,build"

### Step 5: Execute Tests

For each selected target:

1. **Announce**: "Testing: make <target>"
2. **Check dependencies**: Warn if target depends on others not yet run
3. **Execute**: Run `make <target>` with timeout (default 60s, 5s for quick targets)
4. **Capture output**: Show abbreviated output (last 20 lines on success, full on failure)
5. **Report result**: PASS / FAIL / TIMEOUT / SKIPPED

**Special handling:**

| Target Type | Behavior |
|-------------|----------|
| Blocking (`dev`, `run`, `start`) | Skip with note, or run with 5s timeout to verify it starts |
| Destructive (`clean`, `reset`, `down`) | Ask confirmation before running |
| Docker-dependent (`db-up`, `docker-*`) | Check if Docker is running first |
| Long-running (`build`, `test`) | Use extended timeout (5 min) |

### Step 6: Report Results

```
Makefile Test Results for <repo-name>
=====================================

| Target     | Status  | Time   | Notes                    |
|------------|---------|--------|--------------------------|
| install    | PASS    | 12.3s  | Dependencies installed   |
| db-up      | PASS    | 3.2s   | PostgreSQL ready         |
| migrate    | PASS    | 1.1s   | Migrations applied       |
| lint       | PASS    | 2.4s   | No issues found          |
| test       | PASS    | 8.7s   | 36 tests passed          |
| build      | PASS    | 5.2s   | Build successful         |
| dev        | SKIPPED | -      | Blocking target          |
|------------|---------|--------|--------------------------|
| TOTAL      | 6/7     | 32.9s  | 1 skipped                |

Issues Found:
  (none)

Recommendations:
  - All tested targets working correctly
  - Consider adding: make check (runs lint + typecheck + test)
```

---

## Options

Parse flags from `$ARGUMENTS`:

| Flag | Effect |
|------|--------|
| `--dry-run` | Analyze only, don't execute |
| `--all` | Test all targets without prompting |
| `--quick` | Only test quick targets (lint, typecheck) |
| `--force` | Run destructive targets without confirmation |

Examples:
```bash
/make-test my-app --dry-run         # Just analyze Makefile
/make-test infra --all              # Test everything
/make-test frontend --quick         # Only lint/typecheck
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
/make-test                          # Interactive selection
/make-test my-app                   # Fuzzy match to my-nextjs-app
/make-test infra --dry-run          # Analyze without running
/make-test api --quick              # Only quick targets
```
