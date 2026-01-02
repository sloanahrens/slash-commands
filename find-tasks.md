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

   **Optional: MLX Acceleration** - If mlx-hub available and many TODOs found (>10), use Fast tier to batch-summarize:
   ```python
   mlx_infer(
     model_id="mlx-community/Llama-3.2-1B-Instruct-4bit",
     prompt="Summarize each TODO in one line:\n\n{todo_list}",
     max_tokens=256
   )
   ```
   See `_shared-repo-logic.md` for MLX routing rules.

5. Check for incomplete implementation plans:
   ```bash
   ls <repo-path>/docs/plans/*.md 2>/dev/null
   ```

6. **(If `--linear` flag OR repo has `linear_project` in config)** Check Linear issues:

   Use Linear MCP tools to fetch open issues:
   ```
   # Search for issues related to the project
   mcp__linear__search_issues(query: "<project-keywords>", status: "In Progress")
   mcp__linear__search_issues(query: "<project-keywords>", status: "Backlog")
   mcp__linear__search_issues(query: "<project-keywords>", status: "Todo")
   ```

   Display Linear issues in output:
   ```
   Linear Issues (<project-name>):
   ├── In Progress
   │   └── PROJ-123: API redesign (High)
   ├── Backlog
   │   └── PROJ-456: Add retry logic (High)
   └── Todo
       └── (none)
   ```

7. **(If `--issues` flag)** Check GitHub/Bitbucket issues/PRs:
   ```bash
   # Detect remote type from git remote
   git remote get-url origin

   # GitHub
   gh issue list --limit 10
   gh pr list --limit 5

   # Bitbucket (if bb CLI available)
   # Or use API directly
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

## Options

| Flag | Effect |
|------|--------|
| `--linear` | Include Linear issues (auto-enabled if repo has `linear_project` in config) |
| `--issues` | Include GitHub/Bitbucket issues and PRs in analysis |
| `--deep` | More thorough analysis (test coverage, dependency audit) |

---

## Examples

```bash
/find-tasks                    # Interactive selection
/find-tasks pulumi             # Fuzzy match → my-infra-pulumi
/find-tasks my-app             # Auto-includes Linear (has linear_project in config)
/find-tasks my-app --linear    # Explicit Linear flag
/find-tasks frontend --issues  # Include GitHub issues
/find-tasks api --deep         # Deep analysis
```
