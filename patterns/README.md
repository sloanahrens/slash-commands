# Patterns

Proven, reusable knowledge for working in the workspace.

## What belongs here

- Solutions validated across 2+ sessions
- Timeless guidance (not "today I learned...")
- Patterns that prevent repeated mistakes

## What doesn't belong here

- Session-specific notes → `~/.claude/notes/sessions/`
- Raw failure captures → `~/.claude/notes/insights/`
- Temporal or dated content

## Frontmatter format

```yaml
---
tags: [bash, devbot, hookify]     # Searchable tags
repos: [all]                       # Affected repos, or specific names
created: 2026-01-11
updated: 2026-01-11
---
```

## Creating new patterns

Patterns are **promoted** from insights:

1. Capture insights in `~/.claude/notes/insights/<repo>.md`
2. Reference it in 2+ sessions
3. Run `/promote-pattern` to generalize and move here
4. Commit the new pattern

## Searching patterns

```bash
# Find patterns for a repo
grep -l "repos:.*fractals" ~/.claude/patterns/*.md

# Find by tag
grep -l "tags:.*hookify" ~/.claude/patterns/*.md

# Or use /prime <repo> to search automatically
```

## Pattern template

```markdown
---
tags: [tag1, tag2]
repos: [all]
created: YYYY-MM-DD
updated: YYYY-MM-DD
---

# Pattern Title

## Problem
What situation triggers this pattern.

## Solution
The correct approach.

## Why
Root cause or rationale.

## Examples
Concrete usage examples.
```
