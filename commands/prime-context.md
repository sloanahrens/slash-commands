---
description: Load the most recent session note for a repo before starting work
---

# Prime Context

Load context from previous work before starting a session.

**Arguments**: `$ARGUMENTS` - Repo name. See `_shared-repo-logic.md`.

---

## Process

### Step 1: Resolve Repository

```bash
devbot path <repo-name>
```

Confirm: "Priming context for: <repo-name>"

### Step 2: Check for Beads

```bash
devbot exec <repo> ls .beads/
```

**If no `.beads/`**: Ask to initialize (see CLAUDE.md "Initialize Beads in a Repo"). Then continue.

### Step 3: Pull and Show Ready Work

```bash
devbot exec <repo> git fetch origin beads-sync
devbot exec <repo> bd sync --import
devbot exec <repo> bd ready
devbot exec <repo> bd blocked
```

### Step 4: Load Decisions Log

```bash
tail -30 /path/to/repo/.claude/decisions.md 2>/dev/null
```

### Step 5: Check Linear Integration (if configured)

Check `~/.claude/config.yaml` for `linear_projects` field on this repo.

**If not configured**: Skip to output.

**If configured**:

1. Fetch open issues from each project:
   ```
   mcp__plugin_linear_linear__list_issues(project: "<name>")
   ```

2. For each issue, search for matching plans in `plan_paths` (default: `docs/plans`):

   | Match Type | Criteria |
   |------------|----------|
   | Filename | 2+ keywords from issue title in filename |
   | Content | File contains issue URL or ID (e.g., `XYZ-15`) |

3. For matched plans, check if beads reference them (search bead descriptions for plan filename).

| Condition | Action |
|-----------|--------|
| Project not found | Warn, skip |
| No open issues | Note "No outstanding issues" |
| Plan matched, no beads | Flag as "Needs beads" |
| No plan matched | Flag as "Needs planning" |

### Step 6: Output

```
Priming context for: <repo-name>
=====================================

## Linear Issues (if configured)

| Issue  | Title              | Status      | Plan           | Beads          |
|--------|--------------------|-------------|----------------|----------------|
| XYZ-15 | Auth middleware    | In Progress | ✓ auth-plan.md | ✓ 5 tasks (2/5)|
| XYZ-18 | Payment processor  | Backlog     | ✓ pay-plan.md  | — No beads     |
| XYZ-22 | User dashboard     | Todo        | — No plan      | —              |

## Actionable Gaps (if any)

• XYZ-22 needs planning → /super-plan <repo> user-dashboard
• XYZ-18 needs beads → /plan-to-beads <repo> pay-plan.md

## Ready Work
[bd ready output]

## Blocked (if any)
[bd blocked output]

## Recent Decisions
[tail of decisions.md]

---
Next: /super-plan, /execute-plan, or bd show <id>
```

---

## Examples

```bash
/prime-context my-frontend
/prime-context my-api
```

---

## Related

- `/capture-session` — End session, sync beads
- `/execute-plan` — Continue implementation
