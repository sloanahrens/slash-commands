---
description: Load relevant notes and patterns for a repo before starting work
---

# Prime

Search and display relevant patterns and notes to prime context before working on a repo.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Purpose

Surface distilled knowledge from previous sessions so Claude doesn't operate from a blank slate. This implements the "context priming" phase of the Confucius-inspired agent scaffolding.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Priming context for: <repo-name>"

### Step 1.5: Confirm Global CLAUDE.md

Claude Code automatically loads `~/.claude/CLAUDE.md` at session start, but `/prime` should explicitly acknowledge this and surface key reminders:

1. Confirm the file exists and was loaded
2. Surface relevant sections based on repo context:
   - **Bash patterns** — Always relevant (devbot exec, no compound commands)
   - **Tool selection guide** — When repo work involves bash commands
   - **Knowledge capture workflow** — When starting new work

Output format:
```
📋 Global CLAUDE.md loaded
   Key reminders for this session:
   - Use `devbot exec <repo> <cmd>` not `cd && cmd`
   - No Claude/Anthropic attribution in commits
   - Insights auto-captured as you work
```

If `~/.claude/CLAUDE.md` doesn't exist, warn:
```
⚠️ No global CLAUDE.md found at ~/.claude/CLAUDE.md
   Consider running /setup-workspace to initialize.
```

### Step 2: Search Versioned Patterns

Search `~/.claude/patterns/*.md` for patterns that apply:

```bash
# Patterns tagged for this repo or "all"
grep -l "repos:.*<repo-name>\|repos:.*all" ~/.claude/patterns/*.md 2>/dev/null
```

For each match, extract:
- Filename (without path)
- First heading (after frontmatter)
- Tags from frontmatter

### Step 3: Search Local Notes

Search `~/.claude/notes/` for repo-specific notes:

```bash
# Insights file for this repo (one file per repo, accumulated)
cat ~/.claude/notes/insights/<repo-name>.md 2>/dev/null

# Session notes for this repo
grep -l "repos:.*<repo-name>" ~/.claude/notes/sessions/*.md 2>/dev/null
```

Also check the cross-repo insights file:

```bash
cat ~/.claude/notes/insights/all.md 2>/dev/null
```

### Step 4: Display Results

Group by type, most relevant first:

```
Priming context for: fractals-nextjs
=====================================

## Patterns (versioned)
Proven knowledge that applies to this repo:

- **bash-execution.md** — Running commands in repository directories
  Tags: bash, devbot, cd, exec

- **hookify-rules.md** — Hookify blocked commands and workarounds
  Tags: hookify, bash, safety, blocked

## Insights (accumulated)
Learnings captured for this repo:

  [From ~/.claude/notes/insights/fractals-nextjs.md]
  - 2026-01-11 14:30 — Session timeout handling
  - 2026-01-10 09:15 — Form validation patterns
  - 2026-01-08 16:45 — Zapier webhook retry logic

## Session Notes (local)
Recent session summaries:

- **2026-01-10-fractals-nextjs.md** — Zapier webhook integration work
```

### Step 5: Load Key Patterns

Automatically read and display the content of the **most relevant pattern** (first match). This ensures critical context is immediately available.

If there's an insights file for this repo, show a summary of recent entries:

```
📝 Recent insights for fractals-nextjs:
   - Session timeout handling (2026-01-11)
   - Form validation patterns (2026-01-10)
```

---

## Output Format

```
Priming context for: <repo-name>
=====================================

## Patterns (versioned)
- **<filename>** — <first-heading>
  Tags: <tags>

## Insights (accumulated)
  [From ~/.claude/notes/insights/<repo>.md]
  - <date> — <title>
  - <date> — <title>

## Session Notes (local)
- **<filename>** — <first-heading>

---

[Auto-loaded pattern: bash-execution.md]

# Running commands in repository directories
...
```

---

## No Results

If no patterns or notes match:

```
Priming context for: <repo-name>
=====================================

No patterns or insights found for this repo.

Tip: Insights are auto-captured as you work. Run /capture-insight manually
to add specific learnings.
```

---

## Options

| Flag | Effect |
|------|--------|
| `--all` | Show all patterns and notes, not just repo-specific |
| `--tag=<tag>` | Filter by specific tag across all notes |
| `--verbose` | Show full content of all matched notes |
| `--stale` | Include notes marked as stale (normally hidden) |

---

## Tag-Based Search

When using `--tag=<tag>`:

1. Search patterns for tag:
```bash
grep -l "tags:.*<tag>" ~/.claude/patterns/*.md 2>/dev/null
```

2. Search insights for tag:
```bash
grep -l "Tags:.*<tag>" ~/.claude/notes/insights/*.md 2>/dev/null
```

3. Search session notes for tag:
```bash
grep -l "tags:.*<tag>" ~/.claude/notes/sessions/*.md 2>/dev/null
```

4. Display all matches regardless of repo, grouped by type

**Common tags:**
- `hookify` — Hookify rules and workarounds
- `bash` — Bash command patterns
- `devbot` — devbot CLI usage
- `git` — Git operations
- `testing` — Test-related patterns

---

## Pattern Promotion Suggestions

When displaying insights, check if any insight has been referenced multiple times or is particularly valuable:

```
💡 Promotion candidate: "devbot exec patterns" from insights/slash-commands.md
   Referenced in 3 sessions.
   Run /promote-pattern to make it permanent.
```

---

## Examples

```bash
/prime fractals-nextjs      # Prime for fractals work
/prime slash-commands        # Prime for slash-commands work
/prime --tag=hookify         # Find all hookify-related notes
/prime --all                 # Show everything
```

---

## Related Commands

- `/capture-insight` — Manually capture an insight (usually auto-captured)
- `/promote-pattern` — Promote an insight to a versioned pattern
