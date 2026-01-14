---
description: Capture an insight to the repo's insights file (auto or manual)
---

# Capture Insight

Append an insight to the repo's accumulated insights file.

**Arguments**: `$ARGUMENTS` - Optional: `<repo> <insight description>` for manual capture

---

## Purpose

Capture learnings, patterns, and discoveries as they happen. Unlike the old hindsight system (manual, failure-focused), insights are:
- **Accumulated** per repo in a single file
- **Auto-captured** when Claude generates `★ Insight` blocks
- **Learning-focused** (not just failures)

---

## File Structure

One file per repo, appending insights over time:

```
~/.claude/notes/insights/
├── slash-commands.md      # Insights for this repo
├── my-app.md              # Insights for my-app
└── all.md                 # Cross-repo insights
```

---

## Auto-Capture (Primary Method)

When generating an `★ Insight` block during conversation, Claude should also append it to the appropriate insights file:

```markdown
## 2026-01-13 14:30 — Knowledge System Design

The three-tier knowledge system (patterns → insights → sessions) provides
clear promotion paths. Insights accumulate per-repo, patterns are promoted
when proven across sessions.

Tags: architecture, memory, knowledge
```

**Auto-capture triggers:**
- Any `★ Insight` block generated during work
- After discovering a non-obvious pattern
- When a workaround is found for a blocked command
- After resolving a tricky bug

---

## Manual Capture

For explicit capture without an insight block:

```bash
/capture-insight slash-commands hookify blocks cd && patterns
```

This appends to `~/.claude/notes/insights/slash-commands.md`.

---

## Process (Manual)

### Step 1: Determine Repo

If `$ARGUMENTS` starts with a known repo name, use it.
Otherwise, infer from current context or use "all" for cross-repo insights.

### Step 2: Read Existing File

```bash
# Check if insights file exists
cat ~/.claude/notes/insights/<repo>.md 2>/dev/null
```

### Step 3: Append Insight

Format:
```markdown
## YYYY-MM-DD HH:MM — <Brief Title>

<Insight content - what was learned>

Tags: <relevant-tags>
```

### Step 4: Write File

Append to `~/.claude/notes/insights/<repo>.md`

If file doesn't exist, create with header:
```markdown
# Insights: <repo-name>

Accumulated learnings for this repository.

---

## YYYY-MM-DD HH:MM — <First Insight Title>
...
```

### Step 5: Confirm

```
✓ Insight captured: ~/.claude/notes/insights/slash-commands.md

  Added: "Knowledge System Design"
  Total insights in file: 5

  This will appear when you run /prime slash-commands.
```

---

## Auto-Capture Implementation

When Claude generates an insight block like:

```
★ Insight ─────────────────────────────────────
The notes system now has clearer flow...
─────────────────────────────────────────────────
```

Claude should ALSO silently append to the insights file:

1. Determine repo from current context
2. Format as dated entry
3. Append to `~/.claude/notes/insights/<repo>.md`
4. No confirmation needed (silent capture)

---

## Viewing Insights

```bash
# View all insights for a repo
cat ~/.claude/notes/insights/slash-commands.md

# Search across all insights
grep -r "hookify" ~/.claude/notes/insights/
```

Or use `/prime <repo>` which surfaces relevant insights.

---

## Promoting to Patterns

When an insight proves valuable across multiple sessions:

1. Run `/promote-pattern`
2. Select from insights file
3. Generalize and move to `~/.claude/patterns/`

---

## Examples

```bash
/capture-insight                                    # Auto from context
/capture-insight slash-commands devbot exec pattern # Manual with repo
/capture-insight all rate limiting strategies       # Cross-repo insight
```

---

## Related Commands

- `/prime <repo>` — Load insights before starting work
- `/promote-pattern` — Graduate useful insights to patterns
- `/capture-session` — Capture session summary (different purpose)
