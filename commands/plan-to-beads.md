---
description: Convert implementation plan to exhaustive Beads issues
---

# Plan to Beads

Convert an implementation plan into Beads issues for multi-session tracking.

**Arguments**: `$ARGUMENTS` - Repo name and optional plan file path.

**Note**: `/execute-plan` calls this automatically if beads don't exist. Use this command for manual control or to inspect before executing.

---

## Process

### Step 1: Resolve Repository

```bash
devbot path <repo-name>
ls /path/to/repo/.beads/ 2>/dev/null
```

**If no `.beads/`**: Initialize first (see CLAUDE.md).

### Step 2: Find Plan

```bash
ls -t /path/to/repo/docs/plans/*-plan.md 2>/dev/null
```

| Condition | Action |
|-----------|--------|
| Plan specified | Use it |
| Single plan | Use it |
| Multiple plans | Show selection |
| No plans | Error: "Run /super-plan first" |

### Step 3: Parse Plan

Extract:
- Feature name from header
- Linear issue reference (if present):
  ```markdown
  **Linear Issue:** [XYZ-15](https://linear.app/mycompany/issue/XYZ-15)
  ```
- All tasks (`### Task N:` patterns)

### Step 4: Create Beads

**For each task**:
```bash
bd create --title="Task N: <title>" --type=task --priority=2 \
  --description="<full task content from plan>"
```

**Set dependencies** (sequential):
```bash
bd dep add <task-2> <task-1>
bd dep add <task-3> <task-2>
```

**Create feature bead**:
```bash
bd create --title="[Feature] <name>" --type=feature --priority=1 \
  --description="Parent issue for <feature>.

Linear: XYZ-15 - https://linear.app/mycompany/issue/XYZ-15
Plan: docs/plans/<plan>.md"
```

**Link tasks to feature**:
```bash
bd dep add <feature> <task-1>
bd dep add <feature> <task-2>
# ... all tasks block feature
```

### Step 5: Sync and Display

```bash
bd sync
```

```
✓ Created Beads from: docs/plans/<plan>.md

Feature: <id> - [Feature] <name>

Tasks (execution order):
  1. <id>: Task 1: <title>              [ready]
  2. <id>: Task 2: <title>              [blocked by 1]
  3. <id>: Task 3: <title>              [blocked by 2]

Total: N tasks
First ready: <task-1-id>

Next: /execute-plan <repo>
```

---

## Handling Subtasks

If plan has nested tasks:
```markdown
### Task 3: Auth Module
#### 3.1: Login endpoint
#### 3.2: Logout endpoint
```

Create subtasks as separate beads blocking the parent:
```bash
bd create --title="Task 3.1: Login endpoint" --type=task
bd dep add <task-3> <task-3.1>
```

---

## Examples

```bash
/plan-to-beads my-frontend
/plan-to-beads my-frontend auth-plan.md
```

---

## Related

- `/super-plan` — Create the plan
- `/execute-plan` — Execute (auto-creates beads)
