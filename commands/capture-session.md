---
description: Capture or update a session summary for today's work
---

# Capture Session

End-of-session routine: sync Beads, infer and log key decisions, show summary.

**Arguments**: `$ARGUMENTS` - Optional: repo name. If omitted, uses current context or asks.

---

## Purpose

Clean session wrap-up. Automatically captures what happened without prompting:
1. Infers decisions from Beads, commits, and conversation
2. Saves to decisions log
3. Syncs Beads state to git
4. Shows session summary

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

#### 4A.1: Gather Context for Inference

Collect data to infer decisions from:

**Completed Beads today:**

```bash
cd /path/to/repo
bd list --status closed --since today 2>/dev/null
```

**In-progress Beads:**

```bash
bd list --status in_progress 2>/dev/null
```

**Today's commits:**

```bash
git log --since="midnight" --oneline 2>/dev/null
```

**Conversation context:**
- Review current conversation for design decisions, trade-offs discussed, constraints discovered, lessons learned
- Look for "decided to", "chose", "because", "instead of", "blocked by", "discovered that"

#### 4A.2: Infer and Log Decisions (Automatic)

**Analyze the gathered context to extract decisions:**

Look for:
- Architecture choices (why approach A over B)
- Constraints discovered (API limitations, library issues)
- Trade-offs made (performance vs. simplicity)
- Workarounds applied (and why they were necessary)
- Future considerations (things to watch out for)

**Do NOT log:**
- Simple task completions ("added button")
- Routine commits
- Obvious implementations

**If decisions found**, append to `<repo-path>/.claude/decisions.md`:

```markdown
[YYYY-MM-DD] **<brief title>**
<concise explanation of decision and rationale>
```

Create the file with header if it doesn't exist:

```markdown
# Decisions Log

Project-level decisions and context.

---

```

#### 4A.3: Sync and Push Beads

```bash
cd /path/to/repo
bd sync
git push origin beads-sync 2>/dev/null
```

This exports to JSONL, commits to beads-sync branch, and pushes to remote.

#### 4A.4: Update Linear Issues (if configured)

Check if repo has `linear_projects` configured in `~/.claude/config.yaml`.

**If not configured → skip this step.**

**If configured:**

##### 4A.4.1: Gather Active Beads

Get beads with activity today:
```bash
bd list --status=closed --since today
bd list --status=in_progress
```

##### 4A.4.2: Trace Beads to Plans

For each active bead, parse its description for plan file references.
Look for patterns like:
- "See docs/plans/<filename>.md"
- "refs: <filename>.md"
- Any path ending in `-plan.md` or `-design.md`

Build map: `plan-file → [list of beads]`

##### 4A.4.3: Match Plans to Linear Issues

For each plan file found, search for matching Linear issues:

1. Read `linear_projects` and `plan_paths` from config
2. For each Linear project, fetch open issues:
   ```
   mcp__plugin_linear_linear__list_issues(project: "<project-name>")
   ```
3. Match plan to issues using same logic as `/prime-context` Step 4:
   - Keyword match: 2+ keywords from issue title in plan filename
   - Content match: Plan contains Linear issue URL or ID

Build map: `Linear issue → [list of beads with activity]`

##### 4A.4.4: Check for Duplicates (Idempotency)

For each matched Linear issue with bead activity:

1. Fetch recent comments:
   ```
   mcp__plugin_linear_linear__list_comments(issueId: "<issue-id>")
   ```

2. Parse last 3-5 comments for already-reported work:
   - Look for "Completed <task>" patterns
   - Look for "Started <task>" patterns
   - Ignore comments older than 7 days

3. Filter proposed update:
   - Skip tasks already mentioned as "Completed"
   - Skip tasks still "in progress" that were already mentioned as "Started"
   - Include tasks that changed status (was "Started", now "Completed")

##### 4A.4.5: Post Progress Comments

For each Linear issue with NEW progress to report:

1. Format comment as bullet list:
   ```markdown
   Progress update:
   • Completed <task title>
   • Completed <task title>
   • Started <task title>
   ```

2. If more than 5 items, summarize:
   ```markdown
   Progress update:
   • Completed 4 tasks: <task1>, <task2>, <task3>, <task4>
   • Started <task title>
   ```

3. Post comment:
   ```
   mcp__plugin_linear_linear__create_comment(issueId: "<issue-id>", body: "<comment>")
   ```

4. Track results for summary:
   - Posted: issue ID + count of items
   - Skipped: issue ID + reason (no new progress)

**If no Linear issues matched or no new progress:** Skip silently.

#### 4A.5: Show Summary

**After saving**, display:

```
═══════════════════════════════════════════════════════════════════════
                         SESSION SUMMARY
═══════════════════════════════════════════════════════════════════════

Repo: <repo-name>

Completed today:
  ✅ <bead-id>: <title>
  ✅ <bead-id>: <title>

Still in progress:
  🔄 <bead-id>: <title>

Decisions logged:
  • <brief title>: <one-line summary>

Beads: ✓ synced and pushed

Linear updates:
  📝 XYZ-15: Posted progress (2 completed, 1 started)
  ⏭️ XYZ-18: No new progress to report

───────────────────────────────────────────────────────────────────────
Next session: /prime-context <repo-name>
```

**If no decisions inferred:**

```
Decisions logged: (none — routine implementation)
```

**If no Linear integration configured:**

```
Linear updates: (not configured — add linear_projects to config.yaml)
```

**If Linear configured but no updates:**

```
Linear updates: (no matched issues with new progress)
```

---

### Step 4B: Legacy Session Notes (fallback)

Use this path if `.beads/` doesn't exist.

#### 4B.1: Check for Existing Note

```bash
ls /path/to/repo/.claude/sessions/$(date +%Y-%m-%d).md 2>/dev/null
```

If exists, read it for context before updating.

#### 4B.2: Gather Session Context

**Infer from conversation** (do not prompt user):
- Main tasks completed
- Features added/bugs fixed
- Key decisions made
- Unfinished work
- Blockers encountered

**Check git for context:**

```bash
git log --since="midnight" --oneline 2>/dev/null
```

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

#### 4B.6: Show Summary and Suggest Beads

```
═══════════════════════════════════════════════════════════════════════
                         SESSION SUMMARY
═══════════════════════════════════════════════════════════════════════

Repo: <repo-name>
Note: .claude/sessions/YYYY-MM-DD.md

Accomplished:
  • <task 1>
  • <task 2>

Next steps:
  • <unfinished item>

───────────────────────────────────────────────────────────────────────
💡 Consider Beads for better work tracking: bd init
Next session: /prime-context <repo-name>
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

## What Counts as a Decision

**Log these:**
- Why approach A was chosen over B
- Constraints or limitations discovered
- Workarounds and why they're needed
- Architecture/design trade-offs
- Things that will affect future work

**Don't log:**
- "Added feature X" (that's a completion, not a decision)
- "Fixed bug Y" (unless there's a notable why)
- Routine implementation details

---

## Examples

```bash
/capture-session                    # Ask for repo
/capture-session fractals-nextjs   # Capture for specific repo
/capture-session devops-pulumi-ts  # Beads workflow
```

---

## Related Commands

- `/prime-context <repo>` — Load context before starting work
- `/execute-plan <repo>` — Resume plan execution with Beads
