---
description: Start brainstorming session with workspace context
---

# Super Plan

Structured brainstorming → design doc → implementation plan.

**Arguments**: `$ARGUMENTS` - Repo name and topic (or Linear issue ID like `XYZ-15`).

---

## Workflow

```
/super-plan <repo> <topic-or-issue-id>
    ├─▶ Load context (repo status, beads, decisions)
    ├─▶ Brainstorm (superpowers:brainstorming)
    │     • Q&A to understand requirements
    │     • 2-3 approaches with trade-offs
    │     • Design in sections, validate each
    │     • Save: docs/plans/YYYY-MM-DD-<topic>-design.md
    └─▶ Implementation plan (superpowers:writing-plans)
          • Bite-sized TDD tasks
          • Save: docs/plans/YYYY-MM-DD-<topic>-plan.md
```

---

## Process

### Step 1: Resolve Repository

```bash
devbot path <repo>
devbot status <repo>
```

Load beads state and recent decisions if available.

### Step 2: Check for Linear Issue (optional)

If topic looks like a Linear issue ID (e.g., `XYZ-15`):

```
mcp__plugin_linear_linear__get_issue(id: "<issue-id>")
```

Extract title as topic, description as initial requirements. Include URL in plan header.

### Step 3: Brainstorming

**Invoke**: `superpowers:brainstorming`

The skill will:
1. Ask questions one at a time
2. Propose 2-3 approaches with trade-offs
3. Present design in 200-300 word sections
4. Save to `docs/plans/YYYY-MM-DD-<topic>-design.md`

**Do not proceed until design is saved.**

### Step 4: Implementation Plan

**Invoke**: `superpowers:writing-plans`

If Linear issue found, include in header:
```markdown
# [Feature] Implementation Plan

**Linear Issue:** [XYZ-15](https://linear.app/mycompany/issue/XYZ-15) - Title
**Goal:** One sentence

---
```

Save to `docs/plans/YYYY-MM-DD-<topic>-plan.md`

### Step 5: Next Steps

```
✓ Design and plan complete

  Design: docs/plans/YYYY-MM-DD-<topic>-design.md
  Plan:   docs/plans/YYYY-MM-DD-<topic>-plan.md

Next: /plan-to-beads <repo> or /execute-plan <repo>
```

---

## Examples

```bash
/super-plan my-frontend add user authentication
/super-plan my-api XYZ-15              # From Linear issue
/super-plan my-cli config validation
```

---

## Related

- `/plan-to-beads` — Convert plan to beads
- `/execute-plan` — Start execution (auto-creates beads)
- `/prime-context` — Load context
