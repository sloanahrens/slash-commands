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

### Step 3.5: Load Context (Memory Priming)

**Check for Beads first:**

```bash
ls /path/to/repo/.beads/ 2>/dev/null
```

**If Beads exists:**
```bash
cd /path/to/repo
bd ready
bd list --status in_progress 2>/dev/null
```

Display:
```
📝 Loaded context:
   - Beads: 3 ready, 1 in progress
   - Decisions: [if decisions.md exists, note it]
```

**If no Beads (legacy):**
```bash
# Check for project context (external links, stakeholders)
ls /path/to/repo/.claude/project-context.md 2>/dev/null

# Most recent session note for this repo
ls -t /path/to/repo/.claude/sessions/*.md 2>/dev/null | head -1
```

Display:
```
📝 Loaded context:
   - Project: External links, stakeholders (if project-context.md exists)
   - Session: 2026-01-14.md (recent progress)
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

Place plans and design docs in the repo's `docs/plans/` folder:

```
<repo>/docs/plans/YYYY-MM-DD-<topic>-plan.md
<repo>/docs/plans/YYYY-MM-DD-<topic>-design.md
```

**NOT in `.claude/`** — That folder is for local-only context (gitignored).

---

## Post-Brainstorming: Seed Beads Issues

**If `.beads/` exists**, offer to create issues from the brainstorming output:

```
Brainstorming complete. Seed Beads issues from this design?

Options:
- Yes, create issues from action items
- No, I'll create them manually
```

**If yes:**
- Extract action items / next steps from the design
- Create issues with appropriate types and priorities:

```bash
bd create "Implement X" --type feature --priority 2
bd create "Add tests for X" --type task --priority 3
bd dep add <tests-id> <feature-id>  # Tests depend on feature
```

- Show created issues:
```
✓ Created 3 Beads issues:
  - proj-abc: Implement X [P2, feature]
  - proj-def: Add tests for X [P3, task] (blocked by proj-abc)
  - proj-ghi: Update docs [P3, task]
```

---

## Post-Brainstorming Suggestions

After brainstorming completes:

**With Beads:**
```
Brainstorming complete. Next steps:
- bd ready                    — See what's ready to work on
- bd show <id>                — Review issue details
- /run-tests <repo>           — Validate implementation
- /yes-commit <repo>          — Commit changes
```

**Without Beads:**
```
Brainstorming complete. Next steps:
- /capture-session <repo>     — Save decisions and progress
- /run-tests <repo>           — Validate implementation
- /yes-commit <repo>          — Commit changes
```

| Task Type | Suggested Commands |
|-----------|-------------------|
| Feature implementation | `bd ready`, `/run-tests`, `/yes-commit` |
| Bug fix | `bd ready`, `/run-tests` |
| Documentation | `/update-docs`, `/yes-commit` |

---

## Examples

```bash
/super cli add config validation     # Brainstorm for CLI package
/super server add new endpoint       # Brainstorm for server
/super optimize integration          # Prompts for repo selection
/super                               # Shows selection menu
```
