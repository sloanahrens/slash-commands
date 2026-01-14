---
description: Update documentation for a repository
---

# Update Documentation

Update canonical documentation files for a repository to reflect current state.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

---

## Process

1. **Resolve repo** per `_shared-repo-logic.md`
2. **Gather current state** - Use `devbot tree`, `devbot config`, `devbot stats`
3. **Read existing docs** - Check all canonical files (see below)
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

## Canonical Files

Check and update these files if they exist:

| File | Purpose | Limit |
|------|---------|-------|
| `CLAUDE.md` | Claude Code reference | 100-200 lines |
| `README.md` | Human quick start | Under 100 lines |
| `docs/project-context.md` | External resources, stakeholders, decisions | As needed |
| `docs/architecture.md` | Technical architecture | As needed |
| `docs/overview.md` | Project overview (if exists) | As needed |

**CLAUDE.md priorities:** Commands, patterns, warnings/gotchas. Keep concise.

**README.md priorities:** Brief overview, quick start, link to detailed docs.

**project-context.md priorities:** External links (Linear, Notion, Drive), key decisions, stakeholders, open questions. Keep links current.

**docs/ folder:** Active plans and reference docs only. Delete completed plans.

---

## Output

Report:
- Files updated with before/after line counts
- Inconsistencies found and fixed
- Suggestions for improvement (if any)

---

## Auto-Capture Session

After updating docs, automatically run `/capture-session <repo>` to record:
- What documentation was updated
- Any remaining doc tasks as "Next Steps"

This ensures doc updates are tracked for continuity.

---

## Examples

```bash
/update-docs fractals-nextjs        # Update fractals-nextjs docs
/update-docs fractals     # Update fractals-nextjs docs
/update-docs slash        # Update slash-commands docs
```
