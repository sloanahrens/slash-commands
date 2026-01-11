---
description: Update documentation for a repository
---

# Update Documentation

Update CLAUDE.md and README.md for a repository to reflect current state.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

---

## Process

1. **Resolve repo** per `_shared-repo-logic.md`
2. **Gather current state** - Use `devbot tree`, `devbot config`, `devbot stats`
3. **Read existing docs** - CLAUDE.md, README.md, docs/ folder
4. **Update files** to reflect current state
5. **Verify consistency** between docs and actual code

---

## What to Check

- Commands documented actually exist
- Patterns match actual code behavior
- No obsolete instructions or removed features
- Build/test commands are current
- Warnings/gotchas are still relevant

---

## File Guidelines

| File | Purpose | Limit |
|------|---------|-------|
| `CLAUDE.md` | Claude Code reference | 100-200 lines |
| `README.md` | Human quick start | Under 100 lines |
| `docs/*.md` | Detailed docs, plans | As needed |

**CLAUDE.md priorities:** Commands, patterns, warnings/gotchas. Keep concise.

**README.md priorities:** Brief overview, quick start, link to detailed docs.

**docs/ folder:** Active plans and reference docs only. Delete completed plans.

---

## Output

Report:
- Files updated with before/after line counts
- Inconsistencies found and fixed
- Suggestions for improvement (if any)

---

## Examples

```bash
/update-docs mango        # Update mango docs
/update-docs fractals     # Update fractals-nextjs docs
/update-docs slash        # Update slash-commands docs
```
