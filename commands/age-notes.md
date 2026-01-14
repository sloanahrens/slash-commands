---
description: Review and age old notes, pruning stale entries
---

# Age Notes

Review notes that haven't been referenced recently and clean up stale entries.

**Arguments**: `$ARGUMENTS` - Optional: `--days=N` to set threshold (default: 30)

---

## Purpose

Prevent note accumulation by periodically reviewing old notes. Keep the knowledge system focused and relevant.

---

## What Gets Aged

| Type | Structure | Aging Strategy |
|------|-----------|----------------|
| **Insights** | One file per repo | Prune old entries within file |
| **Sessions** | One file per day/repo | Delete old files |

---

## Process

### Step 1: Review Session Notes

Find session notes older than threshold (default 30 days):

```bash
# Session notes older than 30 days
find ~/.claude/notes/sessions -name "*.md" -mtime +30 2>/dev/null
```

For each old session note:
- Check if "Next Steps" items were completed
- Check if insights were captured from the session
- Offer to delete or archive

### Step 2: Review Insight Files

For each repo's insight file:

```bash
ls ~/.claude/notes/insights/*.md 2>/dev/null
```

Parse the file and identify entries older than threshold. Present options:

```
Reviewing: ~/.claude/notes/insights/my-app.md
─────────────────────────────────────────────────────

5 insights total, 2 older than 30 days:

Old entries:
  - 2025-12-01 — API rate limiting discovery
  - 2025-12-10 — Database connection pooling

Options:
1. Keep all — Still relevant
2. Prune old entries — Remove from file
3. Promote to pattern — Move best insights to patterns/
```

### Step 3: Apply Changes

**For session notes:**
```bash
rm ~/.claude/notes/sessions/<filename>
```

**For insight entries:**
Edit the insights file to remove old entries, preserving recent ones.

**For promotion:**
Run `/promote-pattern` workflow with the selected insight.

---

## Batch Mode

With `--batch` flag:

```bash
/age-notes --batch --days=45
```

- Delete session notes older than threshold
- Do NOT auto-prune insights (require manual review)

Output:
```
Batch aging notes older than 45 days...

Deleted sessions:
- 2025-11-20-my-app.md
- 2025-11-25-api-server.md

Insights requiring review:
- insights/my-app.md has 3 entries older than 45 days
- insights/api-server.md has 1 entry older than 45 days

Run /age-notes without --batch to review insights interactively.
```

---

## Cleanup Mode

With `--cleanup` flag, delete old session notes:

```bash
/age-notes --cleanup
```

Output:
```
Cleaning up old session notes...

Will delete:
- 2025-10-01-my-app.md (73 days old)
- 2025-10-15-api-server.md (58 days old)

Proceed? [y/N]
```

---

## Options

| Flag | Effect |
|------|--------|
| `--days=N` | Set age threshold (default: 30) |
| `--batch` | Auto-delete old session notes |
| `--cleanup` | Delete all old session notes |
| `--dry-run` | Show what would happen without changes |

---

## Examples

```bash
/age-notes                    # Interactive review
/age-notes --days=14          # Review notes older than 14 days
/age-notes --batch            # Auto-clean old sessions
/age-notes --cleanup          # Delete old session notes
/age-notes --batch --dry-run  # Preview batch aging
```

---

## Related Commands

- `/prime` — Load notes before starting work
- `/capture-insight` — Capture learnings (usually auto)
- `/promote-pattern` — Graduate insights to patterns
