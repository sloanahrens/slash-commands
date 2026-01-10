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

### Step 2: Load Context

Per `_shared-repo-logic.md` → "Context Loading":
1. Read `~/.claude/CLAUDE.md` (global settings)
2. Read `<repo-path>/CLAUDE.md` (repo-specific guidance)

### Step 3: Check Implementation Plans (FIRST PRIORITY)

**First, get the repo path (REQUIRED):**
```bash
devbot path <repo-name>
# Output: /path/to/repo (use this literal path below)
```

Check for plan documents in `docs/plans/`:
```bash
ls /path/to/repo/docs/plans/*.md 2>/dev/null
```

**If plans exist, process each one:**

1. **Read the plan** and determine its status:
   - Look for completion markers: `✅ COMPLETED`, `Status: Done`, all checklist items checked
   - Compare plan tasks against current codebase state

2. **For completed plans:**
   - Confirm all tasks are actually done by checking the codebase
   - If truly complete, delete the plan file:
     ```bash
     rm /path/to/repo/docs/plans/<plan-file>.md
     ```
   - Note: "Deleted completed plan: <plan-name>"

3. **For incomplete plans:**
   - Check which tasks are done vs pending
   - Update the plan file to reflect current state (mark completed items)
   - Extract remaining tasks as high-priority suggestions

4. **For outdated plans:**
   - If plan references files/code that no longer exists, update or delete
   - If plan's goals are superseded by other work, delete

**Output plan-based tasks first** before other analysis:
```
From Implementation Plans:
1. **<plan-name>: <next-task>** (High)
   - Plan: docs/plans/<file>.md
   - Progress: 2/5 tasks complete
   - Next step: <specific action>
```

### Step 4: Review Current State

1. Read repo documentation:
   - `/path/to/repo/CLAUDE.md` - Repo-specific guidance
   - `/path/to/repo/docs/overview.md` - If exists
   - `/path/to/repo/README.md` - Project overview

2. Check recent commits:
   ```bash
   devbot log <repo-name> -10    # Takes repo NAME, shows last 10 commits
   ```

3. Examine test coverage gaps (if test scripts exist)

4. Look for TODO/FIXME comments:
   ```bash
   devbot todos <repo-name>    # Takes repo NAME
   ```
   This scans for TODO, FIXME, HACK, XXX, BUG markers in parallel (~0.1s).

5. Check for complexity hotspots:
   ```bash
   devbot stats /path/to/repo   # Takes literal PATH
   ```

   Flag any complexity issues as potential refactoring tasks:
   - Large files (>500 lines) → "Consider splitting <file>"
   - Long functions (>50 lines) → "Refactor <function> (<lines> lines)"
   - Deep nesting (>4 levels) → "Simplify control flow in <file>"

**NEVER construct paths manually - always use `devbot path` first.**

### Step 5: Check Linear Issues (Optional)

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

### Step 6: Identify High-Impact Work

Focus on tasks that:
- Unblock other work
- Are assigned in Linear/GitHub
- Improve production readiness
- Are quick wins with high value
- Balance testing, features, and infrastructure

### Step 7: Generate Task Options (Dual-Model Evaluation)

Use dual-model pattern from `_shared-repo-logic.md` to build confidence in local model.

**Note:** If local model is unavailable (see `_shared-repo-logic.md` → "Availability Check"), skip local model steps and use Claude directly. Omit `[local]`/`[claude]` markers in output.

#### 7a. Summarize Each Task with Local Model

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

#### 7b. Claude Reviews Each Summary

**Evaluation criteria:**
- ✓ Starts with actionable verb (Add, Fix, Refactor, Implement)
- ✓ Under 100 characters
- ✓ Accurate file/line reference preserved
- ✓ Priority inference reasonable given context

#### 7c. Build Output with Markers

Mark each task with its provenance:
- `[local]` — local model summary passed all criteria
- `[claude]` — Claude's version used (with reason)

Provide 3-5 concrete, actionable tasks with markers.

---

## Output Format

```
Tasks for: my-cli
======================

From Implementation Plans:
1. **Performance Optimizations: Implement Top-K algorithm** (High)
   - Plan: docs/plans/performance-optimizations-design.md
   - Progress: 1/4 tasks complete
   - Next step: Replace full sort with heap-based Top-K in analyze-overlap

(Deleted completed plan: quick-wins-design.md)

From Linear:
2. **PROJ-123: API redesign** (High, In Progress)
   - Impact: Unblocks mobile team
   - Start: Review current API in src/api/
   - Success: New endpoints pass integration tests

From Codebase Analysis:
3. [local] **Add missing test coverage for config command** (Medium)
   - Impact: Increases confidence in releases
   - Start: src/commands/config.ts (0% coverage)
   - Success: >80% coverage for config module

From TODO Comments:
4. [local] **Implement retry logic in clone setup** (Medium)
   - Location: src/commands/clones.ts:245
   - Impact: Reduces failed clone attempts
   - Success: Retry with exponential backoff

From Complexity Analysis:
5. [claude] **Refactor runStats into smaller focused functions** (Medium)
   - Location: cmd/main.go:793 (127 lines)
   - Impact: Improves maintainability
   - Success: Function under 50 lines
   - (local summary lacked specificity)
```

Note: Linear issues and plan tasks don't get markers (external source, not generated).

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
