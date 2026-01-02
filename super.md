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
1. Read workspace `~/code/CLAUDE.md`
2. Read repo's `CLAUDE.md` (if exists)
3. Get git status and recent commits

### Step 4: Check Related Issues

If brainstorming about a specific feature/bug:
```
mcp__plugin_linear_linear__list_issues with query="<topic>"
mcp__trabian__github_get_assigned_issues_with_project_status
```

### Step 5: Run Brainstorming

Invoke `/superpowers:brainstorming` with:
- Workspace and repo context gathered above
- Task/topic from `$ARGUMENTS`
- Any related Linear/GitHub issues
- Awareness of local model availability (per shared logic)

---

## Documentation Location

When creating documentation, follow trabian's structure:

| Type | Location |
|------|----------|
| Design docs | `~/code/trabian-ai/docs/plans/YYYY-MM-DD-<topic>-design.md` |
| Implementation plans | `~/code/trabian-ai/docs/plans/YYYY-MM-DD-<topic>-plan.md` |
| Knowledge base | `~/code/trabian-ai/docs/<system>/` |
| Technical reviews | `<repo>/docs/tech-review.md` |

**If unsure where docs belong, ASK the user.**

---

## Post-Brainstorming Suggestions

| Task Type | Suggested Commands |
|-----------|-------------------|
| Feature implementation | `/dev/start-session`, `/dev/create-plan` |
| Bug fix | `/sloan/find-tasks`, `/sloan/linear` |
| Infrastructure | `/pm/raid`, review RAID log implications |
| Documentation | `/sloan/update-docs`, `/kb/q2` (if Q2-related) |

---

## Examples

```bash
/sloan/super cli add config validation     # Brainstorm for trabian-cli
/sloan/super server add new endpoint       # Brainstorm for trabian-server
/sloan/super optimize harvest integration  # Prompts for repo selection
/sloan/super                               # Shows selection menu
```
