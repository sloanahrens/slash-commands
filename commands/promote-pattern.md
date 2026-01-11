---
description: Promote a local hindsight note to a versioned pattern
---

# Promote Pattern

Move a validated hindsight note to the versioned patterns directory, generalizing it for long-term use.

**Arguments**: `$ARGUMENTS` - Optional: filename of note to promote (e.g., "2026-01-11-hookify-cd-blocked.md")

---

## Purpose

Graduate proven knowledge from local notes to versioned patterns. Patterns are committed to git and travel with the slash-commands repo, making them permanently available.

**Promotion criteria:**
- Note has been referenced/useful in 2+ sessions
- Pattern is generalizable (not session-specific)
- Knowledge is timeless (not "today I learned...")

---

## Process

### Step 1: List Candidates

If no `$ARGUMENTS`, show promotion candidates:

```bash
# List hindsight notes
ls -t ~/.claude/notes/hindsight/*.md 2>/dev/null
```

Display with metadata:

```
Hindsight notes available for promotion:
=========================================

1. 2026-01-11-hookify-cd-blocked.md
   Tags: hookify, bash, devbot
   Repos: all
   Status: active

2. 2026-01-10-atap-timeout-recovery.md
   Tags: timeout, atap, recovery
   Repos: atap-automation2
   Status: active

3. 2026-01-09-git-worktree-cleanup.md
   Tags: git, worktree
   Repos: all
   Status: promoted  ← Already promoted

Select a note to promote (1-3), or 'q' to quit:
```

### Step 2: Read Selected Note

Read the full content of the selected note:

```bash
~/.claude/notes/hindsight/<selected-file>
```

### Step 3: Generalize Content

Transform the note from temporal/specific to timeless/general:

**Remove:**
- Date-specific language ("Today I...", "Just now...")
- Session-specific context
- Personal pronouns where possible

**Add:**
- Clear problem statement
- Multiple examples if applicable
- "Related" section linking to other patterns

**Transform frontmatter:**
```yaml
# From hindsight format:
---
type: hindsight
repos: [atap-automation2]
tags: [hookify, bash]
created: 2026-01-11
status: active
---

# To pattern format:
---
tags: [hookify, bash]
repos: [all]  # Generalize if applicable
created: 2026-01-11
updated: 2026-01-11
---
```

### Step 4: Generate Pattern Filename

Create a descriptive slug (not date-prefixed):
- `hookify-cd-blocked.md` → `bash-execution.md` or `hookify-compound-commands.md`
- Choose names that describe the **solution**, not the **problem**

### Step 5: Preview and Confirm

Show the transformed pattern:

```
Promoting: 2026-01-11-hookify-cd-blocked.md
Target: ~/.claude/patterns/bash-execution.md

Preview:
---
tags: [bash, devbot, hookify]
repos: [all]
created: 2026-01-11
updated: 2026-01-11
---

# Running commands in repository directories

## Problem
Need to run a command in a repo directory, but hookify blocks compound commands.

## Solution
Use `devbot exec <repo> <command>` instead of `cd && command`.

## Examples
...

---

Proceed with promotion? [Y/n]
```

### Step 6: Write Pattern

Write to `~/.claude/patterns/<filename>`:

```bash
# Get slash-commands path
devbot path slash-commands
```

Use Write tool to create the pattern file.

### Step 7: Update Original Note

Mark the original hindsight note as promoted:

```yaml
---
type: hindsight
repos: [atap-automation2]
tags: [hookify, bash]
created: 2026-01-11
status: promoted  # ← Updated
promoted_to: ~/.claude/patterns/bash-execution.md  # ← Added
---
```

### Step 8: Offer to Commit

```
✓ Pattern created: ~/.claude/patterns/bash-execution.md
✓ Original note marked as promoted

Commit this pattern? [Y/n]
```

If yes, create commit (following commit rules from `_shared-repo-logic.md`):
- No Claude/Anthropic attribution
- Message: "Add bash-execution pattern"

---

## Output Format

```
Promoting hindsight to pattern...

Source: ~/.claude/notes/hindsight/2026-01-11-hookify-cd-blocked.md
Target: ~/.claude/patterns/bash-execution.md

Changes:
- Generalized from atap-automation2 → all repos
- Removed date-specific language
- Added multiple examples
- Added "Related" section

✓ Pattern created: ~/.claude/patterns/bash-execution.md
✓ Original marked: status: promoted

Would you like to commit this pattern? [Y/n]
```

---

## Merge with Existing

If a similar pattern already exists:

```
⚠️ Similar pattern exists: ~/.claude/patterns/bash-execution.md

Options:
1. Merge new content into existing pattern
2. Create separate pattern with different name
3. Cancel promotion

Select (1-3):
```

If merging, update the existing pattern's `updated` date.

---

## Examples

```bash
/promote-pattern                                    # Interactive selection
/promote-pattern 2026-01-11-hookify-cd-blocked.md  # Promote specific note
```

---

## Related Commands

- `/capture-hindsight` — Create hindsight notes
- `/prime <repo>` — Load patterns before starting work
