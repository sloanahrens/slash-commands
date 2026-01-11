---
name: search-notes
description: Search hindsight and session notes for relevant knowledge
subagent_type: general-purpose
---

# Search Notes Template

Search notes for knowledge related to: {keywords}

## Your Task

Search the note system for relevant prior knowledge that could help with the current task.

## Locations to Search

1. **Patterns:** `<slash-commands-path>/docs/patterns/*.md`
2. **Hindsight:** `~/.claude/notes/hindsight/*.md`
3. **Sessions:** `~/.claude/notes/sessions/*.md`

## Search Strategy

1. Grep for keywords in frontmatter (tags, repos)
2. Grep for keywords in content
3. Check recent files (last 14 days) even without keyword match

## Report Format

### Matching Patterns
Versioned patterns that apply:

**{filename}**
- Tags: tag1, tag2
- Relevance: Why this pattern applies
- Key insight: Main takeaway

### Matching Hindsight Notes
Prior failure captures:

**{filename}**
- Date: YYYY-MM-DD
- Tags: tag1, tag2
- Status: active/stale
- Summary: What was learned
- Applicability: How it relates to current task

### Matching Session Notes
Prior session summaries:

**{filename}**
- Date: YYYY-MM-DD
- Repos: repo1, repo2
- Summary: What was done
- Relevant findings: Specific applicable knowledge

### Recommendations
Based on notes found:
1. Consider approach X from pattern Y
2. Avoid mistake Z mentioned in hindsight note
3. Build on session work from note W
