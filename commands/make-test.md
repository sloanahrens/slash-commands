---
description: Test Makefile targets for a repository
---

# Make Test

Run Makefile targets for a repository. Skips clean by default for faster incremental builds.

**Arguments**: `$ARGUMENTS` - Repo name, optionally with flags. See `_shared-repo-logic.md`.

---

## Process

1. **Resolve repo** per `_shared-repo-logic.md`
2. **Parse Makefile** with `devbot make <repo>`
3. **Execute targets** in build order (skip clean by default)
4. **Report results** with timing summary

---

## Execution Order

Run targets in this order:

1. **Setup** - `install`, `setup`, `init`
2. **Infrastructure** - `docker-build`, `db-up`, `migrate`
3. **Build** - `build`, `compile`, `dist`
4. **Verify** - `lint`, `typecheck`, `format`, `check`, `test`
5. **Smoke test** - `dev`, `run`, `start` (5s timeout, then kill)

**Clean targets skipped by default.** Pass `clean` argument to include them.

---

## Special Handling

| Target Type | Behavior |
|-------------|----------|
| Blocking (`dev`, `start`) | 5s timeout to verify startup |
| Docker-dependent | Check Docker running first, prompt if not |
| Long-running (`build`, `test`) | 5 min timeout |
| Format (`format`, `fmt`) | Run but don't fail pipeline |

---

## Error Handling

On failure:
- Capture full error output
- Check common issues: missing deps, Docker not running, port in use
- Run `devbot prereq <repo>` to diagnose environment
- Offer to fix and retry

---

## Options

| Flag | Effect |
|------|--------|
| `clean` | Include clean targets |
| `--dry-run` | Analyze only, don't execute |
| `--quick` | Only lint/typecheck/format |
| `--skip-docker` | Skip Docker targets |

---

## Examples

```bash
/make-test mango              # Run targets (skip clean)
/make-test mango clean        # Include clean targets
/make-test fractals --quick   # Only lint/typecheck
/make-test slash --dry-run    # Analyze without running
```
