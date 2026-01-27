---
description: Execute an implementation plan with Beads tracking
---

# Execute Plan Command

Execute an implementation plan with Beads tracking for multi-session persistence.

**Arguments**: `$ARGUMENTS` - Repo name and/or plan file path. Smart resolution: figures out repo from path, or finds plan from repo.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Purpose

The primary entry point for plan execution. This command:
- Finds the plan file (from args, context, or prompts user)
- Auto-creates Beads if they don't exist (invokes /plan-to-beads logic)
- Maps tasks to Beads and shows progress
- Executes with Beads tracking via superpowers:executing-plans

**Smart behavior:** If you run `/execute-plan <repo>` without having created Beads first, it automatically creates them — no need to run `/plan-to-beads` separately.

---

## Process

### Step 1: Resolve Arguments

Parse `$ARGUMENTS` to determine repo and/or plan file:

**If full path provided** (e.g., `/path/to/repo/docs/plans/feature-plan.md`):
- Extract repo from path
- Use that plan file

**If repo name provided** (e.g., `fractals-nextjs`):
- Resolve repo path via `devbot path <repo>`
- Find/select plan file (Step 3)

**If repo + plan filename** (e.g., `fractals-nextjs auth-plan.md`):
- Resolve repo, look for plan matching filename

```bash
devbot path <repo-name>
```

### Step 2: Verify Beads Initialized

```bash
ls /path/to/repo/.beads/ 2>/dev/null
```

**If no `.beads/`**: Initialize it silently:

```bash
cd /path/to/repo
bd init --branch beads-sync
cat >> .git/info/exclude << 'EOF'
.beads/issues.jsonl
.beads/interactions.jsonl
.beads/metadata.json
EOF
```

Continue — Beads will be created from plan in Step 5.

### Step 3: Find Plan File

```bash
ls -t /path/to/repo/docs/plans/*-plan.md 2>/dev/null
```

**If plan specified in arguments**: Use that plan file.

**If single plan exists**: Use it automatically.

**If multiple plans exist**: Show selection:

```
Available plans for <repo-name>:

1. 2026-01-21-auth-feature-plan.md (today)
2. 2026-01-18-api-refactor-plan.md (3 days ago)

Which plan to execute?
```

**If no plans exist**:

```
No implementation plans found in docs/plans/

Run /super-plan <repo> <topic> to create a design and implementation plan first.
```

### Step 4: Read and Parse Plan

Read the selected plan file. Extract:
- Feature name/goal from the header
- All tasks (look for `### Task N:` patterns)
- Task descriptions, files, steps

### Step 5: Ensure Beads Exist for Plan

Pull latest Beads state:

```bash
cd /path/to/repo
git fetch origin beads-sync 2>/dev/null
bd sync --import 2>/dev/null
```

Search for Beads matching this plan's tasks:

```bash
bd list --status=open
bd list --status=in_progress
bd list --status=closed
```

**Match tasks to Beads** by title (look for "Task N:" pattern match).

**If Beads exist for all tasks**: Continue to Step 6.

**If NO Beads exist or partial match**: Auto-create them silently.

#### 5.1: Auto-Create Beads (plan-to-beads logic)

For each task in the plan without a matching Bead:

```bash
bd create --title="Task N: <task title>" --type=task --priority=2 \
  --description="<full task content from plan>"
```

Set dependencies (each task blocks the next):

```bash
bd dep add <task-2-id> <task-1-id>
bd dep add <task-3-id> <task-2-id>
# etc.
```

Create parent feature Bead:

```bash
bd create --title="[Feature] <feature name>" --type=feature --priority=1 \
  --description="Parent issue. See docs/plans/<plan-file>.md"
```

Link all tasks to block feature:

```bash
bd dep add <feature-id> <task-1-id>
bd dep add <feature-id> <task-2-id>
# ... for all tasks
```

Display:

```
Created Beads from plan:
  [Feature] <name> (<feature-id>)
  ├── Task 1: <title> (<task-1-id>) [ready]
  ├── Task 2: <title> (<task-2-id>) [blocked]
  └── ... N tasks total
```

### Step 6: Show Progress and Task Mapping

```
═══════════════════════════════════════════════════════════════════════
                    BEADS-AWARE PLAN EXECUTION
═══════════════════════════════════════════════════════════════════════

Plan: docs/plans/<plan-file>.md
Feature: <feature-id> - [Feature] <name>

Progress: X/Y tasks complete

  Task                                 Bead           Status
  ────────────────────────────────────────────────────────────
  Task 1: Set up structure             proj-abc123    ✅ completed
  Task 2: Add authentication           proj-def456    🔄 in_progress
  Task 3: Write tests                  proj-ghi789    ⏳ ready
  Task 4: Add dashboard                proj-jkl012    🔒 blocked
  Task 5: Integration tests            proj-mno345    🔒 blocked

Resume point: Task 2 (in progress) or Task 3 (next ready)
```

### Step 7: Determine Resume Point

```bash
bd list --status=in_progress
```

**If task in_progress**: Resume from there.

**If no in_progress**:

```bash
bd ready
```

Start with first ready task.

**If all tasks complete**:

```
All tasks complete for this plan.

Close the feature and sync:
  bd close <feature-id>
  bd sync
```

### Step 8: Execute with Beads Tracking

**REQUIRED SKILL:** `superpowers:executing-plans`

Invoke with explicit Beads context:

```
I'm using superpowers:executing-plans to implement this plan.

BEADS TRACKING (CRITICAL):
- Plan: <path-to-plan>
- Feature: <feature-id>
- Task mapping:
  - Task 1 → <bead-1-id> (completed)
  - Task 2 → <bead-2-id> (in_progress) ← RESUME HERE
  - Task 3 → <bead-3-id> (ready)
  - ...

EXECUTION RULES:
- BEFORE each task: bd update <task-id> --status=in_progress
- AFTER each task:  bd close <task-id>
- Every 3 tasks:    bd sync
- On completion:    bd close <feature-id> && bd sync
```

---

## Progress Reporting

After each task batch:

```
Session progress:
  ✅ Task 2: Add authentication (proj-def456) - closed
  ✅ Task 3: Write tests (proj-ghi789) - closed
  🔄 Task 4: Add dashboard (proj-jkl012) - in progress

Overall: 3/5 tasks complete (60%)
```

---

## Session End

If session ends before completion:

```bash
bd sync
```

```
Session ending. Progress saved to Beads.

To resume: /execute-plan <repo>

State: 3/5 tasks complete, Task 4 in progress
```

---

## Examples

```bash
/execute-plan fractals-nextjs
# → Finds plan, auto-creates Beads if needed, executes

/execute-plan fractals-nextjs docs/plans/2026-01-21-auth-plan.md
# → Uses specific plan

/execute-plan /Users/me/code/fractals/docs/plans/auth-plan.md
# → Extracts repo from path, uses that plan
```

---

## Related Commands

- `/super-plan` — Create design and implementation plan
- `/plan-to-beads` — Manually create Beads (optional, /execute-plan does this automatically)
- `/prime-context` — Load context without executing
- `/capture-session` — Save session summary and sync
