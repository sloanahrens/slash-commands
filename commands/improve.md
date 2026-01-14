---
description: Run improvement loop with parallel exploration
---

# Improve

Spawn parallel subagents to explore a task, synthesize findings, and capture learnings.

**Arguments**: `$ARGUMENTS` - `<repo> <task>` - Repository and task description

---

## Process

1. **Prime context** - Run `/prime <repo>` to load relevant patterns
2. **Parse task** - Identify type: bug fix, feature, exploration, optimization
3. **Spawn subagents** - Launch 2-3 Task agents in parallel for exploration
4. **Synthesize findings** - Analyze combined results, recommend approach
5. **Attempt fix** - If actionable, implement with user approval
6. **Capture learnings** - Run `/capture-insight` to save what was learned

---

## Subagent Patterns

Launch agents in parallel based on task type:

**Bug fixes:**
- Code Explorer: Search relevant code, read context
- History Analyzer: Check recent commits touching related files
- Note Searcher: Search insights for similar issues

**Features:**
- Pattern Finder: Find similar patterns in codebase
- Dependency Mapper: Identify affected files
- Test Scanner: Find existing tests for the area

**Exploration:**
- Structure Analyzer: Map directory structure
- Entry Point Finder: Identify main flows
- Documentation Reader: Read README, CLAUDE.md

---

## Synthesis

After subagents complete, present:
- What was learned
- What patterns emerged
- Recommended approach
- Suggested next steps

If fix is clear, ask: "Proceed with fix? [Y/n]"

---

## Iteration

If task not complete:
1. Update notes with partial progress
2. Ask if user wants another iteration
3. Return to subagent exploration with refined focus

---

## Examples

```bash
/improve fractals-nextjs fix flaky test in payment service
/improve slash-commands add new devbot command
/improve fractals optimize render performance
```

---

## Related Commands

- `/prime` — Load context before improvement
- `/capture-insight` — Save learnings from session
- `/find-tasks` — Discover tasks to improve
