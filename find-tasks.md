---
description: Find high-priority tasks for a repository
---

# Find Tasks

Analyze a repository and suggest 3-5 high-priority tasks.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

---

## Process

1. **Resolve repo** per `_shared-repo-logic.md`
2. **Check implementation plans** in `docs/plans/` (first priority)
3. **Scan for TODOs** with `devbot todos <repo>`
4. **Check complexity** with `devbot stats <path>`
5. **Review recent commits** with `devbot log <repo>`
6. **Present 3-5 prioritized tasks**

---

## What to Check

### Implementation Plans (First Priority)

Check `docs/plans/*.md` for incomplete work:
- If plan complete → delete the file
- If plan incomplete → extract remaining tasks as high priority

### TODO/FIXME Comments

Use `devbot todos <repo>` to find markers (TODO, FIXME, HACK, XXX, BUG).

### Complexity Hotspots

From `devbot stats`:
- Large files (>500 lines) → suggest splitting
- Long functions (>50 lines) → suggest refactoring

### Test Coverage Gaps

Look for untested code paths or modules with no test files.

---

## Priority Levels

| Priority | Criteria |
|----------|----------|
| High | From plans, blocks other work, critical gaps |
| Medium | Improves coverage, adds features, reduces complexity |
| Low | Nice-to-have, documentation, minor fixes |

---

## Output

For each task, include:
- Brief description (imperative mood)
- Location (file:line if applicable)
- Impact (why it matters)
- Suggested starting point

---

## Examples

```bash
/find-tasks mango        # Tasks for mango
/find-tasks fractals     # Tasks for fractals-nextjs
/find-tasks slash        # Tasks for slash-commands
```
