---
description: Start brainstorming session with full monorepo context
---

# Super Command

Start a structured brainstorming session with full context about the workspace and selected repository.

**Arguments**: `$ARGUMENTS` - Optional repo name or task description. If repo recognized, selects it. Otherwise treated as brainstorm topic.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, and commit rules.

---

## Process

### Step 1: Verify Location

```bash
pwd
```

Confirm you are in the configured `base_path` or a subdirectory.

### Step 2: Read Workspace Context

- Read `config.yaml` from this commands directory for base path and repo definitions
- Read the workspace-level `CLAUDE.md` if it exists

### Step 3: Resolve Repository

Follow `_shared-repo-logic.md` for repo selection.

If no repo recognized in `$ARGUMENTS`, ask which repo the task relates to.

### Step 4: Load Repo Context

```bash
pwd  # Verify again before repo commands
cd <base_path>/<repo> && git status
cd <base_path>/<repo> && git log --oneline -5
```

Read: `<repo>/CLAUDE.md`, `README.md`, `docs/overview.md`

### Step 5: Run Brainstorming

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
