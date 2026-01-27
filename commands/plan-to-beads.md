---
description: Convert implementation plan to exhaustive Beads issues
---

# Plan to Beads Command

Convert an implementation plan into an exhaustive set of Beads issues for multi-session execution.

**Arguments**: `$ARGUMENTS` - Repo name and optional plan file path.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Purpose

After `/super-plan` creates a design and implementation plan, this command:
1. Reads the plan file
2. Creates one Bead per task with full description
3. Sets up dependencies between tasks
4. Creates a parent feature Bead
5. Prompts to start execution with `/execute-plan`

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`.

```bash
devbot path <repo-name>
```

### Step 2: Verify Beads Initialized

```bash
ls /path/to/repo/.beads/ 2>/dev/null
```

**If no `.beads/`**: Initialize it first (see `/prime-context` for initialization steps).

### Step 3: Find Implementation Plan

```bash
ls -t /path/to/repo/docs/plans/*-plan.md 2>/dev/null
```

**If plan specified in arguments**: Use that plan file.

**If no plan specified**: Show available plans:

```
Available implementation plans:

1. 2026-01-21-auth-feature-plan.md (today)
2. 2026-01-18-api-refactor-plan.md (3 days ago)

Which plan to convert? (number or filename)
```

### Step 4: Read and Parse Plan

Read the selected plan file. Extract:
- Feature name/goal from the header
- **Linear issue reference** (if present in header)
- All tasks (look for `### Task N:` patterns)
- Task descriptions, files, steps, expected outcomes

**Look for Linear issue in plan header:**

```markdown
**Linear Issue:** [XYZ-15](https://linear.app/mycompany/issue/XYZ-15) - Title
```

Extract:
- Issue ID: `XYZ-15`
- Issue URL: `https://linear.app/mycompany/issue/XYZ-15`

**Store for inclusion in feature bead description.**

**Plan structure expected:**

```markdown
# [Feature Name] Implementation Plan

**Goal:** One sentence description

---

### Task 1: [Task Title]

**Files:**
- Create: path/to/file.py
- Modify: path/to/existing.py

**Step 1:** ...
**Step 2:** ...

### Task 2: [Task Title]
...
```

### Step 5: Create Beads Issues

**CRITICAL:** Create one Bead for EVERY task in the plan. Be exhaustive.

#### 5.1: Create Task Beads

For each task in the plan:

```bash
bd create --title="Task N: <task title>" --type=task --priority=2 \
  --description="<full task content from plan including files, steps, expected outcomes>"
```

**Capture each Bead ID** for dependency setup.

#### 5.2: Set Dependencies

Tasks are typically sequential. Set dependencies so each task blocks the next:

```bash
# Task 2 depends on Task 1
bd dep add <task-2-id> <task-1-id>

# Task 3 depends on Task 2
bd dep add <task-3-id> <task-2-id>

# etc.
```

**If plan specifies different dependencies**, follow those instead.

#### 5.3: Create Parent Feature Bead

**If Linear issue reference found in plan:**

```bash
bd create --title="[Feature] <feature name from plan>" --type=feature --priority=1 \
  --description="Parent issue for <feature>.

Linear: XYZ-15 - https://linear.app/mycompany/issue/XYZ-15
Plan: docs/plans/<plan-file>.md"
```

**If no Linear issue reference:**

```bash
bd create --title="[Feature] <feature name from plan>" --type=feature --priority=1 \
  --description="Parent issue for <feature>. See docs/plans/<plan-file>.md for full implementation plan."
```

Including the Linear URL enables `/capture-session` to trace beads back to Linear issues for progress updates.

**Link all tasks to block the feature:**

```bash
bd dep add <feature-id> <task-1-id>
bd dep add <feature-id> <task-2-id>
# ... for all tasks
```

This ensures the feature can only be closed when all tasks are complete.

### Step 6: Display Summary

```
✓ Created Beads from: docs/plans/<plan-file>.md

Feature: <feature-id> - [Feature] <name>
         Blocked by all tasks below

Tasks (in execution order):
  1. <task-1-id>: Task 1: <title>                    [ready]
  2. <task-2-id>: Task 2: <title>                    [blocked by 1]
  3. <task-3-id>: Task 3: <title>                    [blocked by 2]
  4. <task-4-id>: Task 4: <title>                    [blocked by 3]
  5. <task-5-id>: Task 5: <title>                    [blocked by 4]
  ...

Total: N tasks created
First task ready: <task-1-id>
```

### Step 7: Sync Beads

```bash
bd sync
```

### Step 8: Prompt for Execution

```
Beads created and synced. Ready to start implementation?

  /execute-plan <repo>    — Start Beads-aware execution now
  bd ready                — Review available tasks first
  bd show <task-1-id>     — See first task details
```

---

## Handling Large Plans

For plans with many tasks (10+):

1. **Use parallel subagents** to create Beads faster
2. **Batch dependency setup** after all tasks are created
3. **Show progress** as tasks are created

```
Creating Beads... (15 tasks)
  ✓ Task 1: Set up project structure
  ✓ Task 2: Add authentication module
  ✓ Task 3: Write auth tests
  ...
  ✓ Task 15: Update documentation

Setting dependencies...
  ✓ 14 dependencies created

Creating feature issue...
  ✓ [Feature] User Authentication System
```

---

## Handling Nested/Grouped Tasks

If the plan has task groups or subtasks:

```markdown
### Task 3: Authentication Module

#### 3.1: Add login endpoint
#### 3.2: Add logout endpoint
#### 3.3: Add token refresh
```

Create as nested Beads:

```bash
# Parent task
bd create --title="Task 3: Authentication Module" --type=task --priority=2

# Subtasks (as separate issues, blocking parent)
bd create --title="Task 3.1: Add login endpoint" --type=task --priority=2
bd create --title="Task 3.2: Add logout endpoint" --type=task --priority=2
bd create --title="Task 3.3: Add token refresh" --type=task --priority=2

# Subtasks block the parent
bd dep add <task-3-id> <task-3.1-id>
bd dep add <task-3-id> <task-3.2-id>
bd dep add <task-3-id> <task-3.3-id>
```

---

## Error Handling

### Plan Not Found

```
No implementation plans found in docs/plans/

Run /super-plan <repo> <topic> to create a design and implementation plan first.
```

### Plan Parsing Issues

If task extraction fails:

```
⚠ Could not parse tasks from plan.
  Expected format: ### Task N: <title>

  Found sections:
  - [list what was found]

  Please ensure plan follows the expected structure.
```

### Beads Creation Fails

```bash
bd doctor
bd sync
```

If persistent issues, stop and report. Don't create partial Beads.

---

## Examples

```bash
/plan-to-beads fractals-nextjs
# → Shows available plans, converts selected one to Beads

/plan-to-beads fractals-nextjs docs/plans/2026-01-21-auth-plan.md
# → Directly converts specified plan

/plan-to-beads cli
# → For CLI repo
```

---

## Related Commands

- `/super-plan` — Create design and implementation plan
- `/execute-plan` — Start Beads-aware execution
- `/prime-context` — Load context for a repo
