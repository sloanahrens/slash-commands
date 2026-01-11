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

### Step 2: Search Versioned Patterns

Search `docs/patterns/*.md` for patterns that apply:

```bash
# Patterns tagged for this repo or "all"
grep -l "repos:.*<repo-name>\|repos:.*all" /path/to/slash-commands/docs/patterns/*.md 2>/dev/null
```

For each match, extract:
- Filename (without path)
- First heading (after frontmatter)
- Tags from frontmatter

### Step 3: Search Local Notes

Search `~/.claude/notes/` for repo-specific notes:

```bash
# Hindsight notes for this repo
grep -l "repos:.*<repo-name>" ~/.claude/notes/hindsight/*.md 2>/dev/null

# Session notes for this repo
grep -l "repos:.*<repo-name>" ~/.claude/notes/sessions/*.md 2>/dev/null
```

Also get recent notes (last 7 days) regardless of repo:

```bash
find ~/.claude/notes/hindsight -name "*.md" -mtime -7 2>/dev/null
```

### Step 4: Display Results

Group by type, most relevant first:

```
Priming context for: atap-automation2
=====================================

## Patterns (versioned)
Proven knowledge that applies to this repo:

- **bash-execution.md** — Running commands in repository directories
  Tags: bash, devbot, cd, exec

- **hookify-rules.md** — Hookify blocked commands and workarounds
  Tags: hookify, bash, safety, blocked

## Hindsight Notes (local)
Recent failure captures for this repo:

- **2026-01-11-atap-timeout.md** — ATAP session timeout recovery
  Tags: timeout, recovery, atap

## Session Notes (local)
Recent session summaries:

- **2026-01-10-zapier-integration.md** — Zapier webhook integration work
  Tags: zapier, webhooks, integration

## Recent (all repos, last 7 days)
- **2026-01-11-cd-compound-blocked.md** — Hookify blocks cd && compounds
```

### Step 5: Load Key Patterns

Automatically read and display the content of the **most relevant pattern** (first match). This ensures critical context is immediately available.

If there are hindsight notes tagged `status: active` for this repo, mention them explicitly:

```
⚠️ Active hindsight: 2026-01-11-atap-timeout.md
   ATAP sessions timeout after ~20 minutes. Save progress frequently.
```

---

## Output Format

```
Priming context for: <repo-name>
=====================================

## Patterns (versioned)
- **<filename>** — <first-heading>
  Tags: <tags>

## Hindsight Notes (local)
- **<filename>** — <first-heading>
  Tags: <tags>

## Session Notes (local)
- **<filename>** — <first-heading>
  Tags: <tags>

## Recent (all repos, last 7 days)
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

No patterns or notes found for this repo.

Tip: After encountering issues, run /capture-hindsight to start building
knowledge for future sessions.
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
grep -l "tags:.*<tag>" /path/to/slash-commands/docs/patterns/*.md 2>/dev/null
```

2. Search local notes for tag:
```bash
grep -l "tags:.*<tag>" ~/.claude/notes/hindsight/*.md 2>/dev/null
grep -l "tags:.*<tag>" ~/.claude/notes/sessions/*.md 2>/dev/null
```

3. Display all matches regardless of repo, grouped by type

**Common tags:**
- `hookify` — Hookify rules and workarounds
- `bash` — Bash command patterns
- `devbot` — devbot CLI usage
- `git` — Git operations
- `testing` — Test-related patterns

---

## Pattern Promotion Suggestions

When displaying notes, check if any hindsight note has been active for 14+ days and hasn't been promoted:

```
💡 Promotion candidate: 2026-01-05-monorepo-subdir.md
   Active for 16 days, referenced multiple times.
   Run /promote-pattern to make it permanent.
```

---

## Examples

```bash
/prime atap-automation2      # Prime for ATAP work
/prime slash-commands        # Prime for slash-commands work
/prime --tag=hookify         # Find all hookify-related notes
/prime --all                 # Show everything
```

---

## Related Commands

- `/capture-hindsight` — Create a hindsight note after encountering issues
- `/promote-pattern` — Promote a local note to a versioned pattern
