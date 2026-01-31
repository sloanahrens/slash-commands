---
description: Execute an implementation plan with Beads tracking
---

# Execute Plan

Execute an implementation plan with Beads tracking for multi-session persistence.

**Arguments**: `$ARGUMENTS` - Repo name and/or plan file path.

**Smart behavior**: Auto-creates Beads if they don't exist — no need to run `/plan-to-beads` separately.

---

## Process

### Step 1: Resolve Arguments

| Input | Action |
|-------|--------|
| Full path to plan | Extract repo from path, use that plan |
| Repo name | Find plan in `docs/plans/` |
| Repo + filename | Look for matching plan |

```bash
devbot path <repo-name>
```

### Step 2: Verify Beads

```bash
devbot exec <repo> ls .beads/
```

**If no `.beads/`**: Initialize first (see CLAUDE.md).

### Step 3: Find Plan

```bash
ls -t /path/to/repo/docs/plans/*-plan.md 2>/dev/null
```

| Condition | Action |
|-----------|--------|
| Plan in args | Use it |
| Single plan exists | Use it |
| Multiple plans | Show selection prompt |
| No plans | Error: "Run /super-plan first" |

### Step 4: Parse Plan

Extract from plan file:
- Feature name from header
- Tasks (look for `### Task N:` patterns)
- Linear issue reference if present

### Step 5: Ensure Beads Exist

```bash
devbot exec <repo> git fetch origin beads-sync
devbot exec <repo> bd sync --import
devbot exec <repo> bd list --status open
```

Match tasks to beads by title. **If missing**: Auto-create using `/plan-to-beads` logic:

```bash
devbot exec <repo> bd create --title="Task N: <title>" --type=task --priority=2
devbot exec <repo> bd dep add <task-2> <task-1>  # Sequential dependencies
devbot exec <repo> bd create --title="[Feature] <name>" --type=feature --priority=1
devbot exec <repo> bd dep add <feature> <task-1>  # All tasks block feature
```

### Step 6: Show Progress

```
BEADS-AWARE PLAN EXECUTION
==========================

Plan: docs/plans/<plan>.md
Feature: <id> - [Feature] <name>
Progress: X/Y tasks

Task                          Bead        Status
─────────────────────────────────────────────────
Task 1: Set up structure      prj-abc     ✅ completed
Task 2: Add auth              prj-def     🔄 in_progress
Task 3: Write tests           prj-ghi     ⏳ ready
Task 4: Dashboard             prj-jkl     🔒 blocked

Resume: Task 2 (in progress)
```

### Step 7: Execute with Tracking

**Invoke**: `superpowers:executing-plans`

**Rules**:
- Before task: `devbot exec <repo> bd update <id> --status=in_progress`
- After task: `devbot exec <repo> bd close <id>`
- Every 3 tasks: `devbot exec <repo> bd sync`
- On completion: `devbot exec <repo> bd close <feature>` then `devbot exec <repo> bd sync`

### Step 8: Session End

If ending before completion:
```bash
devbot exec <repo> bd sync
```

Show: "Progress saved. Resume: /execute-plan <repo>"

---

## Examples

```bash
/execute-plan my-frontend
/execute-plan my-frontend auth-plan.md
```

---

## Related

- `/super-plan` — Create plan
- `/plan-to-beads` — Manual bead creation (optional)
- `/capture-session` — End session
