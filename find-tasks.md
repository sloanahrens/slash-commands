---
description: Find high-priority tasks for a repository
---

# Find Next Tasks

Analyze the project and suggest 3-5 high-priority tasks for a repository.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Finding tasks for: <repo-name>"

### Step 2: Review Current State

**First, get the repo path (REQUIRED):**
```bash
REPO_PATH=$(devbot path <repo-name>)
```

1. Read repo documentation:
   - `$REPO_PATH/CLAUDE.md` - Repo-specific guidance
   - `$REPO_PATH/docs/overview.md` - If exists
   - `$REPO_PATH/README.md` - Project overview

2. Check recent commits:
   ```bash
   git -C "$REPO_PATH" log --oneline -10
   ```

3. Examine test coverage gaps (if test scripts exist)

4. Look for TODO/FIXME comments:
   ```bash
   devbot todos <repo-name>    # Takes repo NAME
   ```
   This scans for TODO, FIXME, HACK, XXX, BUG markers in parallel (~0.1s).

5. Check for complexity hotspots:
   ```bash
   devbot stats "$REPO_PATH"   # Takes PATH (use variable from above!)
   ```

   Flag any complexity issues as potential refactoring tasks:
   - Large files (>500 lines) → "Consider splitting <file>"
   - Long functions (>50 lines) → "Refactor <function> (<lines> lines)"
   - Deep nesting (>4 levels) → "Simplify control flow in <file>"

6. Check for incomplete implementation plans:
   ```bash
   ls "$REPO_PATH"/docs/plans/*.md 2>/dev/null
   ```

**NEVER construct paths manually - always use `devbot path` first.**

### Step 3: Check Linear Issues (Optional)

If Linear integration is configured, use Linear MCP to find relevant issues:

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
│   └── PROJ-123: API redesign (High) - https://linear.app/team/issue/PROJ-123
├── Backlog
│   └── PROJ-456: Add retry logic (High)
└── Todo
    └── PROJ-789: Update documentation (Normal)
```

### Step 4: Identify High-Impact Work

Focus on tasks that:
- Unblock other work
- Are assigned in Linear/GitHub
- Improve production readiness
- Are quick wins with high value
- Balance testing, features, and infrastructure

### Step 5: Generate Task Options (Dual-Model Evaluation)

Use dual-model pattern from `_shared-repo-logic.md` to build confidence in local model.

#### 5a. Summarize Each Task with Local Model

For each task identified (TODOs, complexity issues, coverage gaps):

```python
mcp__plugin_mlx-hub_mlx-hub__mlx_infer(
  model_id="mlx-community/Qwen2.5-Coder-14B-Instruct-4bit",
  prompt="""Write a one-line task summary. Use imperative mood, under 80 chars.

Context:
- File: {file_path}:{line}
- TODO/Issue: "{todo_text}"
- Surrounding code context: {context}

Task summary:""",
  max_tokens=50
)
```

#### 5b. Claude Reviews Each Summary

**Evaluation criteria:**
- ✓ Starts with actionable verb (Add, Fix, Refactor, Implement)
- ✓ Under 100 characters
- ✓ Accurate file/line reference preserved
- ✓ Priority inference reasonable given context

#### 5c. Build Output with Markers

Mark each task with its provenance:
- `[local]` — local model summary passed all criteria
- `[claude]` — Claude's version used (with reason)

Provide 3-5 concrete, actionable tasks with markers.

---

## Output Format

```
Tasks for: my-cli
======================

From Linear:
1. **PROJ-123: API redesign** (High, In Progress)
   - Impact: Unblocks mobile team
   - Start: Review current API in src/api/
   - Success: New endpoints pass integration tests

From Codebase Analysis:
2. [local] **Add missing test coverage for config command** (Medium)
   - Impact: Increases confidence in releases
   - Start: src/commands/config.ts (0% coverage)
   - Success: >80% coverage for config module

From TODO Comments:
3. [local] **Implement retry logic in clone setup** (Medium)
   - Location: src/commands/clones.ts:245
   - Impact: Reduces failed clone attempts
   - Success: Retry with exponential backoff

From Complexity Analysis:
4. [claude] **Refactor runStats into smaller focused functions** (Medium)
   - Location: cmd/main.go:793 (127 lines)
   - Impact: Improves maintainability
   - Success: Function under 50 lines
   - (local summary lacked specificity)

Quick Win:
5. [local] **Fix typo in error message** (Low)
   - Location: src/utils/logger.ts:42
   - Impact: Professional error messages
   - Success: Corrected spelling
```

Note: Linear issues don't get markers (external source, not generated).

---

## Priority Levels

| Priority | Criteria |
|----------|----------|
| High | Assigned in Linear/GitHub, addresses critical gaps, unblocks work |
| Medium | Improves test coverage, adds features |
| Low | Nice-to-have improvements, documentation, minor fixes |

---

## Options

| Flag | Effect |
|------|--------|
| `--linear` | Include Linear issues (auto-enabled if repo has `linear_project`) |
| `--deep` | More thorough analysis (test coverage, dependency audit) |

---

## Examples

```bash
/find-tasks                    # Interactive selection
/find-tasks cli                # Tasks for CLI package
/find-tasks server --deep      # Deep analysis of server
```
