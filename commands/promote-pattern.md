---
description: Promote an insight to a versioned pattern
---

# Promote Pattern

Move a validated insight to the versioned patterns directory, generalizing it for long-term use.

**Arguments**: `$ARGUMENTS` - Optional: `<repo>` to select from that repo's insights, or `<repo> <insight-title>` to promote specific insight

---

## Purpose

Graduate proven knowledge from insights to versioned patterns. Patterns are committed to git and travel with the slash-commands repo, making them permanently available.

**Promotion criteria:**
- Insight has been referenced/useful in 2+ sessions
- Pattern is generalizable (not session-specific)
- Knowledge is timeless (not "today I learned...")

---

## Process

### Step 1: List Insight Files

If no `$ARGUMENTS`, show available insight files:

```bash
# List insight files
ls ~/.claude/notes/insights/*.md 2>/dev/null
```

Display:

```
Insight files available:
========================

1. slash-commands.md (8 insights)
2. my-app.md (3 insights)
3. all.md (2 insights)

Select a file to browse (1-3), or 'q' to quit:
```

### Step 2: Browse Insights in File

Parse the selected insight file and show entries:

```
Insights in slash-commands.md:
==============================

1. 2026-01-13 14:30 — Knowledge System Design
   Tags: architecture, memory, knowledge

2. 2026-01-12 09:15 — Hookify compound command patterns
   Tags: hookify, bash, devbot

3. 2026-01-11 16:45 — devbot exec for repo commands
   Tags: devbot, bash, exec

Select an insight to promote (1-3), or 'b' to go back:
```

### Step 3: Read Selected Insight

Extract the full content of the selected insight entry.

### Step 4: Generalize Content

Transform the insight from temporal/specific to timeless/general:

**Remove:**
- Date-specific language ("Today I...", "Just now...")
- Session-specific context
- Personal pronouns where possible

**Add:**
- Clear problem statement
- Multiple examples if applicable
- "Related" section linking to other patterns

**Create pattern frontmatter:**
```yaml
---
tags: [hookify, bash, devbot]
repos: [all]  # Generalize if applicable
created: 2026-01-13
updated: 2026-01-13
---
```

### Step 5: Generate Pattern Filename

Create a descriptive slug (not date-prefixed):
- "Hookify compound command patterns" → `hookify-compound-commands.md`
- Choose names that describe the **solution**, not the **problem**

### Step 6: Preview and Confirm

Show the transformed pattern:

```
Promoting insight: "Hookify compound command patterns"
Source: ~/.claude/notes/insights/slash-commands.md
Target: ~/.claude/patterns/hookify-compound-commands.md

Preview:
---
tags: [hookify, bash, devbot]
repos: [all]
created: 2026-01-13
updated: 2026-01-13
---

# Hookify compound command patterns

## Problem
Hookify blocks compound bash commands (cd && cmd) to prevent...

## Solution
Use `devbot exec <repo> <command>` instead.

## Examples
...

---

Proceed with promotion? [Y/n]
```

### Step 7: Write Pattern

Write to `~/.claude/patterns/<filename>` using the Write tool.

### Step 8: Mark Insight as Promoted

Add a note in the original insights file that this entry was promoted:

```markdown
## 2026-01-12 09:15 — Hookify compound command patterns ✓ PROMOTED

> Promoted to ~/.claude/patterns/hookify-compound-commands.md on 2026-01-13

...original content...
```

### Step 9: Offer to Commit

```
✓ Pattern created: ~/.claude/patterns/hookify-compound-commands.md
✓ Original insight marked as promoted

Commit this pattern? [Y/n]
```

If yes, create commit (following commit rules from `_shared-repo-logic.md`):
- No Claude/Anthropic attribution
- Message: "Add hookify-compound-commands pattern"

---

## Output Format

```
Promoting insight to pattern...

Source: ~/.claude/notes/insights/slash-commands.md
        Entry: "Hookify compound command patterns" (2026-01-12)
Target: ~/.claude/patterns/hookify-compound-commands.md

Changes:
- Generalized from slash-commands → all repos
- Removed date-specific language
- Added multiple examples
- Added "Related" section

✓ Pattern created: ~/.claude/patterns/hookify-compound-commands.md
✓ Original insight marked as promoted

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
/promote-pattern                              # Interactive selection
/promote-pattern slash-commands               # Browse slash-commands insights
/promote-pattern slash-commands hookify       # Promote insight matching "hookify"
```

---

## Related Commands

- `/capture-insight` — Capture learnings (usually auto)
- `/prime <repo>` — Load patterns before starting work
