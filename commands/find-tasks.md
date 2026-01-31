---
description: Find high-priority tasks for a repository
---

# Find Tasks

Find available work for a repository, prioritizing Beads issues.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

---

## Process

### Step 1: Resolve Repository

Follow `_shared-repo-logic.md` for repo selection.

```bash
devbot path <repo-name>
# Output: /path/to/repo
```

### Step 2: Check for Beads

```bash
ls /path/to/repo/.beads/ 2>/dev/null
```

**If `.beads/` exists → Use Beads workflow (Step 3A)**
**If no `.beads/` → Use legacy discovery (Step 3B)**

---

### Step 3A: Beads Workflow (preferred)

#### 3A.1: Sync and Show Ready Work

```bash
cd /path/to/repo
git fetch origin beads-sync 2>/dev/null
bd sync --import 2>/dev/null
bd ready
```

#### 3A.2: Show In-Progress Work

```bash
bd list --status in_progress
```

#### 3A.3: Show Blocked Work (for context)

```bash
bd blocked
```

#### 3A.4: Display Tasks

```
═══════════════════════════════════════════════════════════════════════
                    AVAILABLE TASKS: <repo-name>
═══════════════════════════════════════════════════════════════════════

🔄 In Progress (resume these):
  <bead-id>: <title>
  <bead-id>: <title>

✅ Ready (start these):
  <bead-id>: <title>                              Priority: <0-4>
  <bead-id>: <title>                              Priority: <0-4>

🔒 Blocked (needs dependencies):
  <bead-id>: <title>  ← blocked by <dep-id>

───────────────────────────────────────────────────────────────────────
Commands:
  bd show <id>                    — View task details
  bd update <id> --status in_progress  — Start working
  /execute-plan <repo>            — Resume plan execution
```

**If no Beads issues exist:**

```
No Beads issues found.

Check for work:
  - Implementation plans: docs/plans/*-plan.md
  - TODOs: devbot todos <repo>
  - Start fresh: /super-plan <repo> <topic>
```

---

### Step 3B: Legacy Discovery (fallback)

Use this path if `.beads/` doesn't exist.

#### 3B.1: Check Session Notes (Highest Priority)

```bash
ls -t /path/to/repo/.claude/sessions/*.md 2>/dev/null | head -3
```

For each note, extract unchecked items from "Next Steps":

```markdown
## Next Steps
- [ ] Unfinished task 1   ← Extract these
- [x] Completed task      ← Ignore
- [ ] Unfinished task 2   ← Extract these
```

Session notes represent **explicit continuity** from prior work.

#### 3B.2: Check Implementation Plans (Second Priority)

```bash
ls -t /path/to/repo/docs/plans/*-plan.md 2>/dev/null
```

- If plan incomplete → extract remaining tasks as high priority
- Suggest: `/execute-plan <repo>` to start Beads-aware execution

#### 3B.3: Check TODO/FIXME Comments

```bash
devbot todos <repo>
```

Finds markers: TODO, FIXME, HACK, XXX, BUG.

#### 3B.4: Check Complexity Hotspots

```bash
devbot stats /path/to/repo
```

- Large files (>500 lines) → suggest splitting
- Long functions (>50 lines) → suggest refactoring

#### 3B.5: Display Tasks

```
═══════════════════════════════════════════════════════════════════════
                    AVAILABLE TASKS: <repo-name>
═══════════════════════════════════════════════════════════════════════

📋 From Session Notes:
  • <unfinished task from notes>
  • <unfinished task from notes>

📄 From Implementation Plans:
  • <plan-file>: <incomplete tasks>

🔍 TODOs/FIXMEs:
  • <file:line>: <todo text>

⚠️ Complexity Hotspots:
  • <file>: <lines> lines — consider splitting

───────────────────────────────────────────────────────────────────────
💡 Initialize Beads for better tracking: cd /path/to/repo && bd init
```

---

## Priority Levels

| Priority | Criteria |
|----------|----------|
| P0-P1 | Critical/Urgent — do immediately |
| P2 | Normal — standard work |
| P3-P4 | Low/Backlog — nice to have |

In Beads, priority is set on issues. Without Beads:

| Priority | Source |
|----------|--------|
| High | From plans, session notes, blocks other work |
| Medium | TODOs, coverage gaps, moderate complexity |
| Low | Documentation, minor refactoring |

---

## Examples

```bash
/find-tasks fractals-nextjs   # Tasks for fractals-nextjs
/find-tasks fractals          # Fuzzy match
/find-tasks slash             # Tasks for slash-commands
```

---

## Related Commands

- `/prime-context <repo>` — Load full context before starting
- `/execute-plan <repo>` — Execute plan with Beads tracking
- `/super-plan <repo> <topic>` — Start new feature design
