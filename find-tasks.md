---
description: Find high-priority tasks for a repository
---

# Find Next Tasks (Trabian Branch)

Analyze the project and suggest 3-5 high-priority tasks for a repository, integrating with trabian's Linear and GitHub MCP tools.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Finding tasks for: <repo-name>"

### Step 2: Review Current State

1. Read repo documentation:
   - `~/code/trabian-ai/CLAUDE.md` - Workspace rules (always)
   - `<repo>/CLAUDE.md` - Repo-specific guidance
   - `<repo>/docs/overview.md` - If exists
   - `<repo>/README.md` - Project overview

2. Check recent commits:
   ```bash
   git -C <repo-path> log --oneline -10
   ```

3. Examine test coverage gaps (if test scripts exist)

4. Look for TODO/FIXME comments:
   ```bash
   grep -r "TODO\|FIXME" <repo-path> --include="*.ts" --include="*.tsx" --include="*.py" --include="*.go" | head -20
   ```
   If many TODOs found (>10), consider using local model to batch-summarize (see `_shared-repo-logic.md` → "Local Model Acceleration").

5. Check for incomplete implementation plans:
   ```bash
   ls ~/code/trabian-ai/docs/plans/*.md 2>/dev/null
   ls <repo-path>/docs/plans/*.md 2>/dev/null
   ```

### Step 3: Check Linear Issues (Trabian MCP)

Use trabian's Linear MCP to find relevant issues:

```
# Get issues assigned to me
mcp__plugin_linear_linear__list_issues with assignee="me"

# Search by project if repo has linear_project config
mcp__plugin_linear_linear__list_issues with project="<project-name>"
```

Display Linear issues in output:
```
Linear Issues:
├── In Progress
│   └── TRB-123: API redesign (High) - https://linear.app/trabian/issue/TRB-123
├── Backlog
│   └── TRB-456: Add retry logic (High)
└── Todo
    └── TRB-789: Update documentation (Normal)
```

### Step 4: Check GitHub Issues/PRs (Trabian MCP)

Use trabian's GitHub MCP for project data:

```
# Get assigned issues with project status
mcp__trabian__github_get_assigned_issues_with_project_status

# Get project items if known
mcp__trabian__github_get_project_items with project_id
```

Display GitHub items:
```
GitHub Project Items:
├── Ready
│   └── #42: Fix authentication flow
├── In Progress
│   └── #38: Add MCP endpoint
└── Review
    └── PR #45: Update documentation
```

### Step 5: Check RAID Log (if applicable)

For app repos with project associations, check RAID entries:

```
mcp__trabian__projects_fetch_raid_entries with project_id
```

Flag any:
- Unresolved Issues
- Pending Actions
- Outstanding Risks

### Step 6: Identify High-Impact Work

Focus on tasks that:
- Unblock other work
- Are assigned in Linear/GitHub
- Improve production readiness
- Are quick wins with high value
- Address RAID log items
- Balance testing, features, and infrastructure

### Step 7: Generate Task Options

Provide 3-5 concrete, actionable tasks.

---

## Output Format

```
Tasks for: trabian-cli
======================

From Linear:
1. **TRB-123: API redesign** (High, In Progress)
   - Impact: Unblocks mobile team
   - Start: Review current API in src/api/
   - Success: New endpoints pass integration tests

From Codebase Analysis:
2. **Add missing test coverage for config command** (Medium)
   - Impact: Increases confidence in releases
   - Start: src/commands/config.ts (0% coverage)
   - Success: >80% coverage for config module

From TODO Comments:
3. **Implement retry logic in clone setup** (Medium)
   - Location: src/commands/clones.ts:245
   - Impact: Reduces failed clone attempts
   - Success: Retry with exponential backoff

From RAID Log:
4. **Resolve Action: Update SSH key documentation** (Low)
   - RAID Entry: ACT-12, due 2025-01-05
   - Impact: Reduces support requests
   - Success: Updated docs/tutorial.md

Quick Win:
5. **Fix typo in error message** (Low)
   - Location: src/utils/logger.ts:42
   - Impact: Professional error messages
   - Success: Corrected spelling
```

---

## Priority Levels

| Priority | Criteria |
|----------|----------|
| High | Assigned in Linear/GitHub, addresses critical gaps, unblocks work |
| Medium | Improves test coverage, adds features, addresses RAID items |
| Low | Nice-to-have improvements, documentation, minor fixes |

---

## Options

| Flag | Effect |
|------|--------|
| `--linear` | Include Linear issues (auto-enabled if repo has `linear_project`) |
| `--github` | Include GitHub project items |
| `--raid` | Include RAID log items |
| `--deep` | More thorough analysis (test coverage, dependency audit) |
| `--all` | Enable all integrations |

---

## Examples

```bash
/sloan/find-tasks                    # Interactive selection
/sloan/find-tasks cli                # Tasks for trabian-cli
/sloan/find-tasks server --deep      # Deep analysis of MCP server
/sloan/find-tasks client --all       # All integrations for client project
```
