---
description: Capture or update a session summary for today's work
---

# Capture Session

Create or update a session summary note for today's work on a repo.

**Arguments**: `$ARGUMENTS` - Optional: repo name. If omitted, uses current context or asks.

---

## Purpose

Summarize what was accomplished in a session. Unlike hindsight notes (failure captures), session notes track progress and decisions for continuity across sessions.

**Idempotent**: Running multiple times on the same day updates the existing note.

---

## Process

### Step 1: Determine Repo

If `$ARGUMENTS` provided, use that. Otherwise:
1. Check current working directory for repo context
2. Ask user which repo this session was for

### Step 2: Check for Existing Note

Session notes use date-based filenames:
```
~/.claude/notes/sessions/YYYY-MM-DD-<repo-slug>.md
```

```bash
# Check if today's session note exists
ls ~/.claude/notes/sessions/$(date +%Y-%m-%d)-<repo-slug>.md 2>/dev/null
```

If exists, read it for context before updating.

### Step 3: Gather Session Summary

Ask user (or infer from conversation):

**What was accomplished?**
- Main tasks completed
- Features added/bugs fixed
- Key decisions made

**What's next?**
- Unfinished work
- Blockers encountered
- Recommended next steps

### Step 4: Generate/Update Note

```markdown
---
type: session
repo: <repo-name>
date: YYYY-MM-DD
updated: YYYY-MM-DD HH:MM
status: active
---

# Session: <repo-name> - YYYY-MM-DD

## Accomplished
- <task 1>
- <task 2>

## Key Decisions
- <decision and rationale>

## Next Steps
- [ ] <unfinished item>
- [ ] <follow-up task>

## Notes
<Any additional context for future sessions>
```

### Step 5: Write Note

```bash
mkdir -p ~/.claude/notes/sessions
```

Write to `~/.claude/notes/sessions/YYYY-MM-DD-<repo-slug>.md`

If file exists, **replace it** with updated content (keeps the same filename).

### Step 6: Confirm

```
✓ Session captured: ~/.claude/notes/sessions/2026-01-11-atap-automation2.md

  Repo: atap-automation2
  Status: updated (previous version replaced)

  This note will appear when you run /prime atap-automation2.
```

---

## Quick Capture

From conversation context:
```bash
/capture-session atap-automation2
```

Claude should:
1. Summarize work from current conversation
2. Extract any TODOs or next steps mentioned
3. Note key decisions made
4. Write/update the session file

---

## Multiple Sessions Per Day

If you work on the same repo multiple times in a day, each `/capture-session` **replaces** the previous note for that day. The `updated` timestamp tracks when.

For truly separate sessions on the same day, add a suffix:
```
2026-01-11-atap-automation2-morning.md
2026-01-11-atap-automation2-evening.md
```

---

## Output Format

```
Capturing session for: atap-automation2

---
# Session: atap-automation2 - 2026-01-11

## Accomplished
- Restructured ~/.claude directory (patterns/, templates/ to root)
- Fixed session hooks JSON schema
- Updated all command paths

## Key Decisions
- Moved patterns/templates to root for simpler paths
- Merged SETUP.md into README.md

## Next Steps
- [ ] Test /prime command with new paths
- [ ] Review hindsight promotion workflow

## Notes
Session hooks use systemMessage for Stop events, not hookSpecificOutput.
---

Save to ~/.claude/notes/sessions/2026-01-11-slash-commands.md? [Y/n]
```

---

## Examples

```bash
/capture-session                    # Ask for repo, summarize
/capture-session atap-automation2   # Quick capture for specific repo
/capture-session slash-commands     # Update today's slash-commands session
```

---

## Related Commands

- `/capture-hindsight` — Capture failures and lessons learned
- `/prime <repo>` — Load session notes before starting work
- `/age-notes` — Review old session notes for cleanup
