---
description: Review and age old notes, marking stale ones for cleanup
---

# Age Notes

Review notes that haven't been referenced recently and mark them as stale or archive them.

**Arguments**: `$ARGUMENTS` - Optional: `--days=N` to set threshold (default: 30)

---

## Purpose

Prevent note accumulation by periodically reviewing old notes. Notes that haven't proven useful should be marked stale or deleted to keep the system focused.

---

## Process

### Step 1: Find Old Notes

Find notes older than threshold (default 30 days):

```bash
# Hindsight notes older than 30 days
find ~/.claude/notes/hindsight -name "*.md" -mtime +30 2>/dev/null

# Session notes older than 30 days
find ~/.claude/notes/sessions -name "*.md" -mtime +30 2>/dev/null
```

### Step 2: Check Status

For each old note, read the frontmatter and check:
- `status: active` — Candidate for review
- `status: promoted` — Already graduated, can be deleted
- `status: stale` — Already marked stale

### Step 3: Analyze Usefulness

For `status: active` notes, consider:
- Was it created from this date format? Extract date from filename
- Has similar content been captured in patterns?
- Is the issue/pattern still relevant?

### Step 4: Present Options

For each candidate note:

```
Reviewing: 2025-12-15-api-retry-logic.md (27 days old)
─────────────────────────────────────────────────────
Tags: api, retry, error-handling
Repos: mango

## Summary
Discovered need for exponential backoff in API client...

Options:
1. Keep active — Still relevant
2. Mark stale — Hide from /prime, keep for reference
3. Promote to pattern — Graduate to docs/patterns/
4. Delete — No longer useful
```

### Step 5: Apply Changes

Based on user selection:

**Mark stale:**
```yaml
# Update frontmatter
status: stale
stale_date: 2026-01-11
```

**Delete:**
```bash
rm ~/.claude/notes/hindsight/<filename>
```

**Promote:**
Run `/promote-pattern <filename>` workflow

---

## Batch Mode

With `--batch` flag, automatically mark notes as stale if:
- Older than 30 days (or `--days=N`)
- Status is `active`
- Not referenced in any pattern

```bash
/age-notes --batch --days=45
```

Output:
```
Batch aging notes older than 45 days...

Marked stale:
- 2025-11-20-docker-build-cache.md
- 2025-11-25-env-var-loading.md

Skipped (already stale): 2
Skipped (promoted): 1

Total: 2 notes marked stale
```

---

## Cleanup Mode

With `--cleanup` flag, delete all stale notes:

```bash
/age-notes --cleanup
```

Output:
```
Cleaning up stale notes...

Will delete:
- 2025-10-01-old-issue.md (stale for 40 days)
- 2025-10-15-resolved-bug.md (stale for 25 days)

Proceed? [y/N]
```

---

## Options

| Flag | Effect |
|------|--------|
| `--days=N` | Set age threshold (default: 30) |
| `--batch` | Auto-mark old active notes as stale |
| `--cleanup` | Delete all stale notes |
| `--dry-run` | Show what would happen without making changes |

---

## Examples

```bash
/age-notes                    # Interactive review of old notes
/age-notes --days=14          # Review notes older than 14 days
/age-notes --batch            # Auto-mark old notes as stale
/age-notes --cleanup          # Delete all stale notes
/age-notes --batch --dry-run  # Preview batch aging
```

---

## Related Commands

- `/prime` — Load notes (hides stale by default, use `--stale` to include)
- `/capture-hindsight` — Create new hindsight notes
- `/promote-pattern` — Graduate useful notes to patterns
