---
description: Find high-priority tasks for a repository
---

# Find Next Tasks

Analyze the project and suggest 3-5 high-priority tasks for a repository.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Finding tasks for: <repo-name>"

### Step 2: Review Current State

1. Read repo documentation:
   - `<repo>/CLAUDE.md` - Primary reference
   - `<repo>/docs/overview.md` - If exists
   - `<repo>/README.md` - Project overview

2. Check recent commits:
   ```bash
   cd <repo-path> && git log --oneline -10
   ```

3. Examine test coverage gaps (if test scripts exist)

4. Look for TODO/FIXME comments (use language-appropriate patterns):
   ```bash
   # TypeScript/JavaScript
   grep -r "TODO\|FIXME" <repo-path> --include="*.ts" --include="*.tsx" --include="*.js" | head -20

   # Go
   grep -r "TODO\|FIXME" <repo-path> --include="*.go" | head -20

   # Python
   grep -r "TODO\|FIXME" <repo-path> --include="*.py" | head -20
   ```

5. Check for incomplete implementation plans:
   ```bash
   ls <repo-path>/docs/plans/*.md 2>/dev/null
   ```

### Step 3: Identify High-Impact Work

Focus on tasks that:
- Unblock other work
- Improve production readiness
- Are quick wins with high value
- Balance testing, features, and infrastructure

### Step 4: Generate Task Options

Provide 3-5 concrete, actionable tasks.

---

## Output Format

For each task:

1. **Task Name** - Clear, actionable title
2. **Priority** - High/Medium/Low with justification
3. **Impact** - What this accomplishes
4. **Starting Point** - Key files or commands
5. **Dependencies** - Prerequisites or blockers
6. **Success Criteria** - How to know it's done

---

## Priority Levels

| Priority | Criteria |
|----------|----------|
| High | Addresses critical gaps, unblocks work, improves stability |
| Medium | Improves test coverage, adds features, enhances monitoring |
| Low | Nice-to-have improvements, optimizations, documentation |

---

## Examples

```bash
/find-tasks              # Interactive selection
/find-tasks pulumi       # Fuzzy match → devops-gcp-pulumi
/find-tasks atap         # Fuzzy match → atap-automation2
```
