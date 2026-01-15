---
description: Capture or update a session summary for today's work
---

# Capture Session

Create or update a session summary note for today's work on a repo.

**Arguments**: `$ARGUMENTS` - Optional: repo name. If omitted, uses current context or asks.

---

## Purpose

Summarize what was accomplished in a session. Session notes track progress and decisions for continuity across sessions.

**Idempotent**: Running multiple times on the same day updates the existing note.

---

## Process

### Step 1: Determine Repo

If `$ARGUMENTS` provided, use that. Otherwise:
1. Check current working directory for repo context
2. Ask user which repo this session was for

### Step 2: Check for Existing Note

Session notes live in the repo's `.claude/sessions/` directory:
```
<repo-path>/.claude/sessions/YYYY-MM-DD.md
```

```bash
# Get repo path
devbot path <repo-name>
# Output: /path/to/repo

# Check if today's session note exists
ls /path/to/repo/.claude/sessions/$(date +%Y-%m-%d).md 2>/dev/null
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

### Step 3.5: Find Previous Session Notes

Search for existing session notes for this repo:

```bash
# Previous sessions are in the same repo directory
ls -t /path/to/repo/.claude/sessions/*.md 2>/dev/null | head -5
```

If previous sessions exist, note their filenames for the "Related" section.

### Step 4: Generate/Update Note

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
- <other related notes or docs>
```

**IMPORTANT:** The "Related" section MUST include links to previous session notes if they exist. This creates a navigable history chain.

### Step 5: Write Note

Write to `<repo-path>/.claude/sessions/YYYY-MM-DD.md`

If file exists, **replace it** with updated content (keeps the same filename).

**Note:** Ensure `.claude/sessions/` directory exists in the repo. Create if needed.

### Step 6: Confirm

```
✓ Session captured: /path/to/fractals-nextjs/.claude/sessions/2026-01-11.md

  Repo: fractals-nextjs
  Status: updated (previous version replaced)
  Linked to: 2026-01-10.md (previous session)

  This note will appear when you run /prime fractals-nextjs.
```

### Step 7: Check for CLAUDE.md Updates

Review the session for tooling discoveries or pattern changes that should be reflected in the global `~/.claude/CLAUDE.md`:

**Check for:**
- New devbot commands used or discovered
- New bash patterns or workarounds
- Hookify rule encounters and solutions
- New slash commands created or modified
- Workflow changes that affect multiple repos

**If tooling changes detected**, suggest updates:

```
💡 CLAUDE.md update suggestions:

   The session involved tooling changes that may warrant updating ~/.claude/CLAUDE.md:

   - New devbot command: `devbot prereq` used for dependency checking
     → Consider adding to "devbot CLI" section

   - New pattern discovered: `npm run --prefix` for package.json scripts
     → Consider adding to "Bash Patterns" alternatives table

   Update CLAUDE.md now? [y/N]
```

**Keep CLAUDE.md general:**
- Only suggest changes that apply to ALL repos or the tooling itself
- Repo-specific guidance belongs in repo's CLAUDE.md
- Focus on: commands, tools, workflows, critical rules

If no tooling changes: skip this step silently.

---

## Quick Capture

From conversation context:
```bash
/capture-session fractals-nextjs
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
2026-01-11-morning.md
2026-01-11-evening.md
```

---

## Examples

```bash
/capture-session                    # Ask for repo, summarize
/capture-session fractals-nextjs   # Quick capture for specific repo
/capture-session slash-commands     # Update today's slash-commands session
```

---

## Related Commands

- `/prime <repo>` — Load session notes before starting work
