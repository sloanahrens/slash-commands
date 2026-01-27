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
ls /path/to/repo/.beads/ 2>/dev/null
```

**If no `.beads/`**: Initialize silently:
```bash
cd /path/to/repo && bd init --branch beads-sync
```

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
git fetch origin beads-sync 2>/dev/null
bd sync --import 2>/dev/null
bd list --status=open
```

Match tasks to beads by title. **If missing**: Auto-create using `/plan-to-beads` logic:

```bash
bd create --title="Task N: <title>" --type=task --priority=2
bd dep add <task-2> <task-1>  # Sequential dependencies
bd create --title="[Feature] <name>" --type=feature --priority=1
bd dep add <feature> <task-1>  # All tasks block feature
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
- Before task: `bd update <id> --status=in_progress`
- After task: `bd close <id>`
- Every 3 tasks: `bd sync`
- On completion: `bd close <feature> && bd sync`

### Step 8: Session End

If ending before completion:
```bash
bd sync
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
