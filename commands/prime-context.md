---
description: Load the most recent session note for a repo before starting work
---

# Prime Context

Load context from previous work before starting a session.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Purpose

Surface context from previous sessions so Claude doesn't operate from a blank slate. Uses Beads for structured work tracking when available, with decisions log for narrative context.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Priming context for: <repo-name>"

### Step 2: Get Repo Path

```bash
devbot path <repo-name>
# Output: /path/to/repo
```

### Step 3: Check for Linear Integration

Check if repo has `linear_projects` configured in `~/.claude/config.yaml`.

**Read config:**
```bash
# Check config.yaml for this repo's linear_projects field
```

**If `linear_projects` is configured:**

1. For each project name in the array, fetch open issues:
   ```
   mcp__plugin_linear_linear__list_issues(project: "<project-name>")
   ```

2. Filter to open issues only (exclude: Done, Cancelled, Duplicate)

3. Collect for each issue:
   - Issue ID (e.g., "XYZ-15")
   - Title
   - Status (Backlog, Todo, In Progress, In Review)
   - Assignee (or "—" if unassigned)
   - URL

**Edge cases:**
- Project name not found in Linear → warn: "Linear project 'X' not found, skipping"
- No open issues → note: "No outstanding Linear issues"
- No `linear_projects` in config → skip this step entirely

**Store results** for use in Step 4 (Plan Matching) and output formatting.

---

### Step 4: Plan Matching (if Linear issues found)

For each Linear issue from Step 3, search for matching plan documents.

**Get plan paths:**
```bash
# Read plan_paths from config.yaml for this repo
# Default to ["docs/plans"] if not configured
```

**Search algorithm:**

1. Extract keywords from Linear issue title:
   - Split on spaces and common separators
   - Remove stop words: "the", "a", "an", "for", "to", "and", "or", "in", "on", "at"
   - Lowercase all keywords

2. For each path in `plan_paths`:
   - Resolve path (expand `~`, make relative paths absolute from repo root)
   - Skip silently if path doesn't exist
   - Glob for `*.md` files in that directory

3. For each plan file, check for match:
   - **Filename match**: 2+ keywords appear in the filename
   - **Content match**: File contains Linear issue URL (`linear.app/mycompany/issue/XYZ-15`) or issue ID (`XYZ-15`)
   - A plan matches if EITHER condition is true

**Example:**
```
Linear issue: "Payment processor integration" (XYZ-18)
Keywords: ["payment", "processor", "integration"]

Search paths:
  /code/my-project/docs/plans/*.md
  /code/shared-docs/plans/*.md

Results:
  ✓ payment-processor-integration-plan.md (filename: "payment" + "processor" + "integration")
  ✓ 2026-01-20-payments-api.md (content contains "XYZ-18")
```

**Store matches** for each Linear issue (may have 0, 1, or multiple plan matches).

---

### Step 4.5: Beads Correlation (if plans found)

For each plan file that matched a Linear issue, check if beads exist that reference it.

**Prerequisites:**
- Repo must have `.beads/` directory (check in Step 5)
- If no beads initialized, skip this step and mark all as "— No beads"

**Search strategy:**

1. Get all open beads:
   ```bash
   bd list --status=open
   ```

2. For each matched plan file, search bead descriptions for:
   - The plan filename (e.g., "payment-processor-plan.md")
   - The plan path (e.g., "docs/plans/payment-processor-plan.md")

3. If a feature bead references the plan:
   - Count how many task beads are blocked by this feature bead
   - Count how many of those tasks are completed
   - Report: "✓ 5 tasks (2/5 done)"

4. If no bead references the plan:
   - Mark as "— No beads yet"
   - Include in actionable gaps (suggest `/plan-to-beads`)

**Example:**
```
Plan: payment-processor-plan.md
  → Found bead: .prj-abc "[Feature] Payment processor integration"
  → Blocks 5 task beads, 2 completed
  → Status: "✓ 5 tasks (2/5 done)"

Plan: user-dashboard-plan.md
  → No matching beads found
  → Status: "— No beads yet"
```

**Store correlation results** for output formatting.

---

### Step 5: Check for Beads

```bash
ls /path/to/repo/.beads/ 2>/dev/null
```

**If `.beads/` exists → Use Beads workflow (Step 6A)**
**If no `.beads/` → Offer to initialize (Step 5B)**

---

### Step 5B: Initialize Beads (if not present)

Use AskUserQuestion:

```
Beads not initialized for <repo-name>. Set up issue tracking?

Options:
- Yes, initialize Beads (Recommended)
- No, use legacy session notes
```

**If Yes → Run initialization:**

```bash
cd /path/to/repo

# Initialize with sync branch
bd init --branch beads-sync

# Add JSONL files to local exclude
cat >> .git/info/exclude << 'EOF'
.beads/issues.jsonl
.beads/interactions.jsonl
.beads/metadata.json
EOF

# Stage config files
git add .beads/.gitignore .beads/config.yaml .beads/README.md .gitattributes AGENTS.md

# Commit
git commit -m "Add Beads config for protected branch workflow"

# Push (if remote configured)
git push origin $(git branch --show-current) 2>/dev/null
git push origin beads-sync 2>/dev/null
```

Display:
```
✓ Beads initialized for <repo-name>
  Prefix: <repo-name>
  Sync branch: beads-sync
```

#### 5B.2: Migrate from Session Notes (if exist)

```bash
ls -t /path/to/repo/.claude/sessions/*.md 2>/dev/null | head -1
```

**If session notes exist:**

1. Read the most recent session note
2. Look for "Next Steps" or "TODO" sections with unchecked items (`- [ ]`)
3. Display found items and ask:

```
Found existing session notes with unfinished work:

From 2026-01-20.md:
- [ ] Add validation to user input
- [ ] Write tests for auth module
- [ ] Update API documentation

Import these as Beads issues?

Options:
- Yes, create issues from next steps
- No, start fresh
```

**If Yes:**
- Create a Beads issue for each unchecked item:
  ```bash
  bd create "Add validation to user input" --type task
  bd create "Write tests for auth module" --type task
  bd create "Update API documentation" --type task
  ```
- Display created issues

**If No session notes or user declines:** Continue without migration.

**Then continue with Step 6A (Beads workflow).**

**If No to Beads initialization → Continue with Step 6B (legacy session notes).**

---

### Step 6A: Beads Workflow (preferred)

#### 6A.1: Pull remote changes (if remote configured)

```bash
cd /path/to/repo
git fetch origin beads-sync 2>/dev/null
bd sync --import 2>/dev/null
```

This pulls any changes from other machines/sessions.

#### 6A.2: Run bd ready

```bash
bd ready
```

This shows unblocked work ready to pick up.

#### 6A.3: Show blocked issues (if any)

```bash
bd blocked 2>/dev/null | head -10
```

#### 6A.4: Load decisions log (if exists)

```bash
tail -30 /path/to/repo/.claude/decisions.md 2>/dev/null
```

Show recent decisions for context.

#### 6A.5: Output format (Beads)

```
Priming context for: <repo-name>
=====================================

## Linear Issues (<project-names>)

| Issue   | Title                          | Status      | Assignee | Plan                              | Beads           |
|---------|--------------------------------|-------------|----------|-----------------------------------|-----------------|
| XYZ-15  | Auth middleware                | In Progress | user@    | ✓ auth-middleware-plan.md         | ✓ 5 tasks (2/5) |
| XYZ-18  | Payment processor integration  | Backlog     | —        | ✓ payment-processor-plan.md       | — No beads      |
| XYZ-22  | User dashboard setup           | Todo        | alice@   | — No plan                         | —               |

(If no linear_projects configured, skip this section)
(If no open issues, show: "No outstanding Linear issues")

## Actionable Gaps

**Needs planning (no plan file):**
  • XYZ-22: User dashboard setup
    → Run `/super-plan <repo> user-dashboard`

**Needs beads (plan exists, no beads):**
  • XYZ-18: Payment processor integration
    → Run `/plan-to-beads <repo> payment-processor-plan.md`

(If no gaps, skip this section)

## Ready Work (from Beads)
[output from bd ready]

## Blocked (if any)
[output from bd blocked]

## Recent Decisions
[tail of decisions.md if exists]

---
Next steps:
  • `/super-plan <repo> <topic>` — Design unplanned work
  • `/plan-to-beads <repo>` — Convert plan to trackable tasks
  • `/execute-plan <repo>` — Continue implementation
  • `bd show <id>` — View issue details
```

---

### Step 6B: Legacy Session Notes (fallback)

Use this path if `.beads/` doesn't exist.

#### 6B.1: Confirm Global CLAUDE.md

```
📋 Global CLAUDE.md loaded
   Key reminders:
   - Use `devbot exec <repo> <cmd>` not `cd && cmd`
   - No Claude/Anthropic attribution in commits
```

#### 6B.2: Load Project Context (if exists)

```bash
ls /path/to/repo/.claude/project-context.md 2>/dev/null
```

If exists, read and summarize key info.

#### 6B.3: Load Most Recent Session Note

```bash
ls -t /path/to/repo/.claude/sessions/*.md 2>/dev/null | head -1
```

Read and display the full content.

If no session note exists:
```
📝 No session notes found for <repo-name>
   Consider initializing Beads: cd /path/to/repo && bd init
```

#### 6B.4: Output format (legacy)

```
Priming context for: <repo-name>
=====================================

📋 Global CLAUDE.md loaded

## Project Context (if exists)
[Key info from .claude/project-context.md]

## Most Recent Session
[Full content of most recent session note]

---
Ready to continue where you left off.
```

---

## Beads Quick Reference

When working in a Beads-enabled repo:

| Action | Command |
|--------|---------|
| See ready work | `bd ready` |
| Issue details | `bd show <id>` |
| Start working | `bd update <id> --status in_progress` |
| Create issue | `bd create "Title" --type task` |
| Complete work | `bd close <id>` |

---

## Options

| Flag | Effect |
|------|--------|
| `--verbose` | Show all open issues, not just ready |
| `--full` | Also run `bd prime` for full workflow context |

---

## Examples

```bash
/prime-context my-frontend          # Prime for frontend work
/prime-context my-api               # Prime for API work
/prime-context infra-pulumi         # Prime with Beads workflow
```

---

## Related Commands

- `/capture-session` — Save decisions and sync Beads
