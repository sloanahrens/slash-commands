---
description: Start brainstorming session with workspace context
---

# Super Command

Start a structured brainstorming session with full context about the workspace and selected repository.

**Arguments**: `$ARGUMENTS` - Optional repo name or task description. If repo recognized, selects it. Otherwise treated as brainstorm topic.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, and commit rules.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`:
1. Read `config.yaml` for base path and repo definitions
2. Match `$ARGUMENTS` to repo name or alias
3. If no repo recognized, ask which repo the task relates to
4. Confirm: "Brainstorming for: <repo-name>"

### Step 2: Load Repo Context

```bash
pwd  # Verify again before repo commands
cd <base_path>/<repo> && git status
cd <base_path>/<repo> && git log --oneline -5
```

Read: `<repo>/CLAUDE.md`, `README.md`, `docs/overview.md`

### Step 3: Run Brainstorming

Invoke `/superpowers:brainstorming` with:
- Selected repo name and path
- Task/topic from `$ARGUMENTS`
- Key context from repo's CLAUDE.md
- Current git status

---

## Documentation Location

When creating documentation:

| Type | Location |
|------|----------|
| Technical reviews | `<repo>/docs/tech-review.md` |
| Design docs | `<repo>/docs/plans/<date>-<topic>-design.md` |
| Implementation plans | `<repo>/docs/plans/<date>-<topic>-plan.md` |

**If unsure where docs belong, ASK the user.**

---

## Examples

```bash
/super my-app add user authentication   # Brainstorm for my-app repo
/super optimize database queries        # Prompts for repo selection
/super pulumi                           # Start brainstorming for infra repo
```
