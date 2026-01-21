---
description: Capture or update a session summary for today's work
---

# Capture Session

End-of-session routine: sync Beads and optionally log key decisions.

**Arguments**: `$ARGUMENTS` - Optional: repo name. If omitted, uses current context or asks.

---

## Purpose

Clean session wrap-up. With Beads, work items are tracked as you go, so this command focuses on:
1. Syncing Beads state to git
2. Prompting for any key decisions to log
3. Showing session summary

---

## Process

### Step 1: Determine Repo

If `$ARGUMENTS` provided, use that. Otherwise:
1. Check current working directory for repo context
2. Ask user which repo this session was for

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

#### 4A.1: Show Session Summary

```bash
cd /path/to/repo
bd list --status closed --since today 2>/dev/null
bd list --status in_progress 2>/dev/null
```

Display:
```
Session summary for: <repo-name>
================================

Completed today:
[list of closed issues]

Still in progress:
[list of in_progress issues]
```

#### 4A.2: Prompt for Decisions

Use AskUserQuestion:

```
Any key decisions to log?

Options:
- Yes, log a decision
- No, just sync
```

**If yes:**
- Ask: "What decision or learning should be recorded?"
- Append to `<repo-path>/.claude/decisions.md`:

```markdown
[YYYY-MM-DD] **<brief title>**
<user's input>
```

Create the file if it doesn't exist (with header).

#### 4A.3: Sync Beads

```bash
cd /path/to/repo
bd sync
```

#### 4A.4: Confirm

```
✓ Session captured for: <repo-name>

  Beads: synced
  Decisions: [logged / no new decisions]

  Next session: /prime <repo-name>
```

---

### Step 4B: Legacy Session Notes (fallback)

Use this path if `.beads/` doesn't exist.

#### 4B.1: Check for Existing Note

```bash
ls /path/to/repo/.claude/sessions/$(date +%Y-%m-%d).md 2>/dev/null
```

If exists, read it for context before updating.

#### 4B.2: Gather Session Summary

Ask user (or infer from conversation):

**What was accomplished?**
- Main tasks completed
- Features added/bugs fixed
- Key decisions made

**What's next?**
- Unfinished work
- Blockers encountered
- Recommended next steps

#### 4B.3: Find Previous Session Notes

```bash
ls -t /path/to/repo/.claude/sessions/*.md 2>/dev/null | head -5
```

#### 4B.4: Generate/Update Note

```markdown
---
type: session
repo: <repo-name>
date: YYYY-MM-DD
updated: YYYY-MM-DD HH:MM
status: active
tags: <relevant-tags>
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

## Related
- Previous session: [YYYY-MM-DD.md](filename) — <brief description>
```

#### 4B.5: Write Note

Write to `<repo-path>/.claude/sessions/YYYY-MM-DD.md`

Ensure `.claude/sessions/` directory exists.

#### 4B.6: Suggest Beads

```
💡 Consider initializing Beads for better work tracking:
   cd /path/to/repo && bd init
```

---

## Decisions Log Format

When logging decisions, append to `<repo>/.claude/decisions.md`:

```markdown
# Decisions Log

Project-level decisions and context.

---

[2026-01-20] **Chose JWT over sessions**
Client's mobile app can't handle cookies. Refresh token rotation every 7 days.

[2026-01-21] **Fixed tar vulnerability**
Used pnpm override. Transitive dep from npm internals.
```

**Rules:**
- Date + bold title + brief explanation (1-3 sentences)
- Append only — never edit old entries
- Only log decisions, learnings, constraints — not task completions

---

## Examples

```bash
/capture-session                    # Ask for repo
/capture-session fractals-nextjs   # Capture for specific repo
/capture-session devops-pulumi-ts  # Beads workflow
```

---

## Related Commands

- `/prime <repo>` — Load context before starting work
