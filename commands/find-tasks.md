---
description: Find high-priority tasks for a repository
---

# Find Tasks

Analyze a repository and suggest 3-5 high-priority tasks.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

---

## Process

1. **Resolve repo** per `_shared-repo-logic.md`
2. **Check session notes** for unfinished "Next Steps" (highest priority)
3. **Check insights** for blockers or known issues
4. **Check implementation plans** in `<repo>/docs/`
5. **Scan for TODOs** with `devbot todos <repo>`
6. **Check complexity** with `devbot stats <path>`
7. **Review recent commits** with `devbot log <repo>`
8. **Present 3-5 prioritized tasks**

---

## What to Check

### Session Notes (Highest Priority)

Check recent session notes for this repo:

```bash
# Get repo path
devbot path <repo-name>
# Output: /path/to/repo

# Find session notes for this repo
ls -t /path/to/repo/.claude/sessions/*.md 2>/dev/null | head -3
```

For each note, extract unchecked items from "Next Steps":
```markdown
## Next Steps
- [ ] Unfinished task 1   ← Extract these
- [x] Completed task      ← Ignore
- [ ] Unfinished task 2   ← Extract these
```

Session notes represent **explicit continuity** from prior work — prioritize these over discovered TODOs.

### Implementation Plans (Second Priority)

Check `<repo>/docs/*.md` for incomplete plans:
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
/find-tasks fractals-nextjs        # Tasks for fractals-nextjs
/find-tasks fractals     # Tasks for fractals-nextjs
/find-tasks slash        # Tasks for slash-commands
```
