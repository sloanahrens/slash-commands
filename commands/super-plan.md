---
description: Start brainstorming session with workspace context
---

# Super Plan Command

Start a structured brainstorming session that creates a design and implementation plan.

**Arguments**: `$ARGUMENTS` - Optional repo name or task description. If repo recognized, selects it. Otherwise treated as brainstorm topic.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, context loading.

---

## Workflow Overview

```
/super-plan <repo> <topic-or-issue-id>
    │
    ├─▶ 1. Load Context
    │      • Repo status, structure, existing Beads
    │      • Recent decisions and session context
    │      • Linear issue details (if issue ID provided)
    │
    ├─▶ 2. Brainstorming (superpowers:brainstorming)
    │      • Interactive Q&A to understand requirements
    │      • Explore 2-3 approaches with trade-offs
    │      • Present design in sections, validate each
    │      • Save design to docs/plans/YYYY-MM-DD-<topic>-design.md
    │
    └─▶ 3. Implementation Plan (superpowers:writing-plans)
           • Convert design into bite-sized TDD tasks
           • Save to docs/plans/YYYY-MM-DD-<topic>-plan.md
           • Prompt: run /plan-to-beads to create tracking issues
```

---

## Process

### Step 1: Check Prerequisites

#### 1.1: Verify Superpowers Plugin

```bash
claude plugin list 2>/dev/null | grep superpowers
```

Required skills:
- `superpowers:brainstorming` - Interactive design creation
- `superpowers:writing-plans` - Detailed implementation planning

If not installed, offer: `claude plugin add superpowers@superpowers-marketplace`

#### 1.2: Resolve Repository

Follow "Standard Process Start" from `_shared-repo-logic.md`.

```bash
devbot path <repo>
```

### Step 2: Load Context

```bash
# Load status and structure
devbot status <repo>
devbot stats /path/to/repo
devbot tree /path/to/repo
```

**Load Beads state (if exists):**

```bash
ls /path/to/repo/.beads/ 2>/dev/null
cd /path/to/repo
bd ready 2>/dev/null
bd list --status in_progress 2>/dev/null
```

**Load decisions log (if exists):**

```bash
tail -30 /path/to/repo/.claude/decisions.md 2>/dev/null
```

Display:
```
📋 Context loaded for <repo-name>
   - Beads: X ready, Y in progress (or "not initialized")
   - Decisions: [recent context if exists]
```

### Step 2.5: Check for Linear Issue (optional)

If the topic argument looks like a Linear issue ID (e.g., `XYZ-15`, `ABC-42`):

1. Fetch issue details:
   ```
   mcp__plugin_linear_linear__get_issue(id: "<issue-id>")
   ```

2. Extract context:
   - Issue title → use as topic/feature name
   - Issue description → include in brainstorming context
   - Issue URL → store for inclusion in plan file

3. Display:
   ```
   📎 Linear issue: XYZ-15 - Auth middleware refactor
      Status: In Progress | Assignee: user@

   Using issue description as initial requirements.
   ```

**If not a Linear issue ID:** Treat argument as a topic description (existing behavior).

**Store Linear context** for use in plan file header (Step 4).

### Step 3: Brainstorming Phase

**REQUIRED SKILL:** `superpowers:brainstorming`

Invoke the brainstorming skill with:
- Workspace and repo context gathered above
- Task/topic from `$ARGUMENTS`
- Existing Beads issues for awareness

The brainstorming skill will:
1. Ask questions one at a time to understand requirements
2. Propose 2-3 approaches with trade-offs
3. Present design in 200-300 word sections, validating each
4. Save validated design to `docs/plans/YYYY-MM-DD-<topic>-design.md`

**Do NOT proceed to Step 4 until brainstorming explicitly says the design is complete and saved.**

### Step 4: Implementation Planning Phase

**REQUIRED SKILL:** `superpowers:writing-plans`

After design is saved, invoke the writing-plans skill to:
1. Convert the design into bite-sized tasks (2-5 minutes each)
2. Each task follows TDD: write test → verify fails → implement → verify passes → commit
3. Include exact file paths, complete code, exact commands
4. Save to `docs/plans/YYYY-MM-DD-<topic>-plan.md`

**If Linear issue context exists**, include in plan file header:

```markdown
# [Feature Name] Implementation Plan

**Linear Issue:** [XYZ-15](https://linear.app/mycompany/issue/XYZ-15) - Auth middleware refactor
**Goal:** One sentence description

---
```

This enables `/plan-to-beads` and `/capture-session` to trace back to the Linear issue.

### Step 5: Prompt for Next Steps

After the implementation plan is saved:

```
✓ Design and implementation plan complete

  Design: docs/plans/YYYY-MM-DD-<topic>-design.md
  Plan:   docs/plans/YYYY-MM-DD-<topic>-plan.md

Next step:
  /plan-to-beads <repo>   — Convert plan to tracked Beads issues

This creates one Bead per task with dependencies for multi-session execution.
```

---

## Documentation Location

Place designs and plans in the repo's `docs/plans/` folder:

```
<repo>/docs/plans/YYYY-MM-DD-<topic>-design.md   # From brainstorming
<repo>/docs/plans/YYYY-MM-DD-<topic>-plan.md     # From writing-plans
```

**NOT in `.claude/`** — That folder is for local-only context (gitignored).

---

## Error Handling

### Brainstorming Incomplete

If user wants to skip brainstorming:

```
⚠ Brainstorming ensures we build the right thing.
  Skipping risks wasted effort on wrong approach.

Continue anyway?
- Yes, I have a clear design already
- No, let's brainstorm first (Recommended)
```

### Plan Already Exists

If `docs/plans/*-<topic>-plan.md` already exists:

```
Found existing plan: docs/plans/2026-01-20-feature-plan.md

Options:
- Use existing plan (skip to /plan-to-beads)
- Create new plan (overwrites existing)
- Review existing plan first
```

---

## Examples

```bash
/super-plan my-frontend add user authentication
# → Brainstorm auth approaches
# → Create implementation plan
# → Prompt to run /plan-to-beads

/super-plan my-api XYZ-15
# → Fetches Linear issue XYZ-15 details
# → Uses issue description as requirements
# → Creates plan linked to XYZ-15

/super-plan my-cli add config validation
# → Brainstorm validation strategy
# → Create implementation plan
# → Prompt to run /plan-to-beads

/super-plan                    # Shows repo selection, then topic prompt
```

---

## Related Commands

- `/plan-to-beads` — Convert implementation plan to Beads issues
- `/execute-plan` — Start Beads-aware execution
- `/prime-context` — Load context from previous sessions
- `/capture-session` — Save decisions and sync Beads
