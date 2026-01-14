---
description: Start brainstorming session with workspace context
---

# Super Command

Start a structured brainstorming session with full workspace and repo context.

**Arguments**: `$ARGUMENTS` - Optional repo name or task description. If repo recognized, selects it. Otherwise treated as brainstorm topic.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, context loading, local model usage, and integration patterns.

---

## Process

### Step 1: Check Superpowers Plugin

The superpowers plugin provides `/superpowers:brainstorming` (required).

```bash
claude plugin list 2>/dev/null | grep superpowers
```

If not installed, offer: `claude plugin add superpowers@superpowers-marketplace`

If user declines, **stop** - cannot complete without the brainstorming skill.

### Step 2: Resolve Repository

Follow "Standard Process Start" from `_shared-repo-logic.md`.

### Step 3: Load Context

Per `_shared-repo-logic.md` → "Context Loading":
1. Read global `~/.claude/CLAUDE.md` (if exists)
2. Read repo's `CLAUDE.md` (if exists)
3. Get git status and recent commits

**Use devbot for speed:**
```bash
# First get the path (REQUIRED for tree/stats commands)
devbot path <repo>
# Output: /path/to/repo (use this literal path below)

# Then use repo name OR path as appropriate
devbot status <repo>          # Git status (~0.02s) - takes NAME
devbot stats /path/to/repo    # Complexity metrics - takes literal PATH
devbot tree /path/to/repo     # Directory structure - takes literal PATH
```

**NEVER use compound commands or construct paths manually.**

### Step 3.5: Load Relevant Notes (Memory Priming)

Search for patterns and insights that might inform this brainstorm:

```bash
# Patterns tagged for this repo or "all"
grep -l "repos:.*<repo-name>\|repos:.*all" ~/.claude/patterns/*.md 2>/dev/null

# Insights for this repo
cat ~/.claude/notes/insights/<repo-name>.md 2>/dev/null
```

If matches found, briefly summarize relevant insights before brainstorming:
```
📝 Loaded context from notes:
   - Pattern: bash-execution.md (command execution in repos)
   - Insights: 3 entries for <repo-name>
```

This prevents repeating past mistakes during ideation.

### Step 4: Check Related Issues (Optional)

If brainstorming about a specific feature/bug and Linear integration is configured:
```
mcp__plugin_linear_linear__list_issues with query="<topic>"
```

### Step 5: Run Brainstorming

Invoke `/superpowers:brainstorming` with:
- Workspace and repo context gathered above
- Task/topic from `$ARGUMENTS`
- Any related issues
- Awareness of local model availability (per shared logic)

---

## Documentation Location

Place documentation in the relevant repository:

| Type | Location |
|------|----------|
| Design docs | `<repo>/docs/YYYY-MM-DD-<topic>-design.md` |
| Implementation plans | `<repo>/docs/YYYY-MM-DD-<topic>-plan.md` |
| Technical reviews | `<repo>/docs/tech-review.md` |

**If unsure where docs belong, ASK the user.**

---

## Post-Brainstorming Suggestions

After brainstorming completes, suggest:

```
Brainstorming complete. Next steps:
- /capture-session <repo>  — Save decisions and progress for future sessions
- /run-tests <repo>        — Validate implementation
- /yes-commit <repo>       — Commit changes
```

| Task Type | Suggested Commands |
|-----------|-------------------|
| Feature implementation | `/capture-session`, `/run-tests`, `/yes-commit` |
| Bug fix | `/capture-insight` (if tricky), `/find-tasks` |
| Documentation | `/update-docs`, `/capture-session` |

---

## Examples

```bash
/super cli add config validation     # Brainstorm for CLI package
/super server add new endpoint       # Brainstorm for server
/super optimize integration          # Prompts for repo selection
/super                               # Shows selection menu
```
