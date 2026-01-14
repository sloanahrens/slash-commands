---
description: Load the most recent session note for a repo before starting work
---

# Prime

Load the most recent session note for a repo to get context from previous work.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Purpose

Surface context from previous sessions so Claude doesn't operate from a blank slate. Session notes link to previous sessions, forming a chain of context.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Priming context for: <repo-name>"

### Step 2: Confirm Global CLAUDE.md

Claude Code automatically loads `~/.claude/CLAUDE.md` at session start. Acknowledge this:

```
📋 Global CLAUDE.md loaded
   Key reminders:
   - Use `devbot exec <repo> <cmd>` not `cd && cmd`
   - No Claude/Anthropic attribution in commits
```

If `~/.claude/CLAUDE.md` doesn't exist, warn:
```
⚠️ No global CLAUDE.md found at ~/.claude/CLAUDE.md
   Consider running /setup-workspace to initialize.
```

### Step 3: Load Project Context (if exists)

Check for project context file:

```bash
# Get repo path
devbot path <repo-name>
# Output: /path/to/repo

# Check for project context
ls /path/to/repo/.claude/project-context.md 2>/dev/null
```

If exists, read and summarize key info:
- External links (Linear, Notion, etc.)
- Key stakeholders
- Important decisions

### Step 4: Load Most Recent Session Note

Find the most recent session note:

```bash
# Find most recent session note in repo
ls -t /path/to/repo/.claude/sessions/*.md 2>/dev/null | head -1
```

Read and display the full content. Session notes contain:
- What was accomplished
- Key decisions made
- Action items and next steps
- Links to previous related sessions

If no session note exists:
```
📝 No session notes found for <repo-name>
   Session notes are created via /capture-session after working.
```

---

## Output Format

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

## Following the Chain

Session notes include a "Related" section linking to previous sessions. If more context is needed:

1. Check the "Related" links in the loaded session
2. Read previous sessions as needed
3. The chain provides full history without searching

---

## Options

| Flag | Effect |
|------|--------|
| `--verbose` | Also load previous linked sessions |

---

## Examples

```bash
/prime fractals-nextjs      # Prime for fractals work
/prime hanscom-plaid        # Prime for hanscom work
/prime --verbose            # Load session + linked previous sessions
```

---

## Related Commands

- `/capture-session` — Save session progress and decisions
