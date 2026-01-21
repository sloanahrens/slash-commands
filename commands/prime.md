---
description: Load the most recent session note for a repo before starting work
---

# Prime

Load context from previous work before starting a session.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Purpose

Surface context from previous sessions so Claude doesn't operate from a blank slate. Uses Beads for structured work tracking when available, with decisions log for narrative context.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Priming context for: <repo-name>"

### Step 2: Get Repo Path

```bash
devbot path <repo-name>
# Output: /path/to/repo
```

### Step 3: Check for Beads

```bash
ls /path/to/repo/.beads/ 2>/dev/null
```

**If `.beads/` exists → Use Beads workflow (Step 4A)**
**If no `.beads/` → Use legacy session notes (Step 4B)**

---

### Step 4A: Beads Workflow (preferred)

#### 4A.1: Run bd ready

```bash
cd /path/to/repo
bd ready
```

This shows unblocked work ready to pick up.

#### 4A.2: Show blocked issues (if any)

```bash
bd blocked 2>/dev/null | head -10
```

#### 4A.3: Load decisions log (if exists)

```bash
tail -30 /path/to/repo/.claude/decisions.md 2>/dev/null
```

Show recent decisions for context.

#### 4A.4: Output format (Beads)

```
Priming context for: <repo-name>
=====================================

## Ready Work
[output from bd ready]

## Blocked (if any)
[output from bd blocked]

## Recent Decisions
[tail of decisions.md if exists]

---
Ready to continue. Use `bd show <id>` for issue details.
```

---

### Step 4B: Legacy Session Notes (fallback)

Use this path if `.beads/` doesn't exist.

#### 4B.1: Confirm Global CLAUDE.md

```
📋 Global CLAUDE.md loaded
   Key reminders:
   - Use `devbot exec <repo> <cmd>` not `cd && cmd`
   - No Claude/Anthropic attribution in commits
```

#### 4B.2: Load Project Context (if exists)

```bash
ls /path/to/repo/.claude/project-context.md 2>/dev/null
```

If exists, read and summarize key info.

#### 4B.3: Load Most Recent Session Note

```bash
ls -t /path/to/repo/.claude/sessions/*.md 2>/dev/null | head -1
```

Read and display the full content.

If no session note exists:
```
📝 No session notes found for <repo-name>
   Consider initializing Beads: cd /path/to/repo && bd init
```

#### 4B.4: Output format (legacy)

```
Priming context for: <repo-name>
=====================================

📋 Global CLAUDE.md loaded

## Project Context (if exists)
[Key info from .claude/project-context.md]

## Most Recent Session
[Full content of most recent session note]

---
Ready to continue where you left off.
```

---

## Beads Quick Reference

When working in a Beads-enabled repo:

| Action | Command |
|--------|---------|
| See ready work | `bd ready` |
| Issue details | `bd show <id>` |
| Start working | `bd update <id> --status in_progress` |
| Create issue | `bd create "Title" --type task` |
| Complete work | `bd close <id>` |

---

## Options

| Flag | Effect |
|------|--------|
| `--verbose` | Show all open issues, not just ready |
| `--full` | Also run `bd prime` for full workflow context |

---

## Examples

```bash
/prime fractals-nextjs      # Prime for fractals work
/prime hanscom-plaid        # Prime for hanscom work
/prime devops-pulumi-ts     # Prime with Beads workflow
```

---

## Related Commands

- `/capture-session` — Save decisions and sync Beads
- `/switch` — Switch context to another repo (calls /prime)
