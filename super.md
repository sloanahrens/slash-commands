---
description: Start brainstorming session with workspace context
---

# Super Command (Trabian Branch)

Start a structured brainstorming session with full context about the workspace and selected repository. Leverages local MLX models for acceleration where appropriate and integrates with trabian's documentation and planning patterns.

**Arguments**: `$ARGUMENTS` - Optional repo name or task description. If repo recognized, selects it. Otherwise treated as brainstorm topic.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 0: Check Plugins

Before running the brainstorming session, verify required and optional plugins are available:

```bash
claude plugin list 2>/dev/null | grep -E "(superpowers|mlx)"
```

#### 0.1: Superpowers Plugin (Required)

The superpowers plugin provides the `/superpowers:brainstorming` skill used in Step 5.

**If superpowers is NOT installed:**

1. Inform the user: "The superpowers plugin is required for structured brainstorming but isn't currently installed."
2. Offer to install it:
   ```bash
   claude plugin add superpowers@superpowers-marketplace
   ```
3. If the user declines or installation fails, **stop here** - the `/super` command cannot complete without the brainstorming skill. Suggest alternatives:
   - Install the plugin and retry
   - Use basic context gathering without the structured brainstorming workflow

**If superpowers IS installed:** Continue to the next check.

#### 0.2: MLX-Hub Plugin (Optional)

Check if mlx-hub is available for local model acceleration. See `_shared-repo-logic.md` for availability check pattern.

If not installed, continue without acceleration - install via `/setup-plugins` if desired.

---

### Step 1: Load Workspace Context

**Always load trabian workspace context first:**

```bash
cat ~/trabian/CLAUDE.md
```

Key context to extract:
- Repository structure (packages/, mcp/, clones/, .trees/)
- Build commands (npm run build, uv sync)
- MCP server architecture
- Available commands (kb/, pm/, dev/)

### Step 2: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`:
1. Parse `config.yaml` for builtin, worktrees, clones, repos
2. Match `$ARGUMENTS` to repo name or alias
3. If no repo recognized, ask which repo the task relates to
4. Confirm: "Brainstorming for: <repo-name>"

### Step 3: Load Repo Context

**For builtin packages:**
```bash
cat ~/trabian/packages/<name>/CLAUDE.md  # or mcp/<name>
git -C <repo-path> status
git -C <repo-path> log --oneline -5
```

**For worktrees:**
```bash
git -C ~/.trees/<name> branch --show-current
git -C ~/.trees/<name> log --oneline main..HEAD
```

**For clones (reference repos):**
```bash
cat ~/trabian/clones/<name>/README.md
# Note: These are read-only references
```

**For apps:**
```bash
cat <repo-path>/CLAUDE.md
git -C <repo-path> status
```

### Step 4: Local Model Acceleration (Optional)

See `_shared-repo-logic.md` for model tiers and routing rules.

**Use local models for:**
- Summarizing file contents (Fast tier)
- Drafting exploration approaches (Quality tier, if available)
- Parallel task processing while Claude orchestrates

**Keep on Claude:**
- Architectural decisions
- Security-sensitive analysis
- Final synthesis and recommendations

### Step 4b: Check Related Issues (Trabian MCP)

If brainstorming about a specific feature/bug, check for related issues:

```
# Search Linear for related issues
mcp__plugin_linear_linear__list_issues with query="<topic>"

# Check GitHub assigned issues
mcp__trabian__github_get_assigned_issues_with_project_status
```

### Step 5: Run Brainstorming

Invoke `/superpowers:brainstorming` with:
- Trabian workspace context
- Selected repo name, type, and path
- Task/topic from `$ARGUMENTS`
- Key context from repo's CLAUDE.md
- Current git status
- Any related Linear/GitHub issues
- **Awareness of local model for acceleration**

---

## Documentation Location (Trabian Pattern)

When creating documentation, follow trabian's structure:

| Type | Location |
|------|----------|
| Design docs | `~/trabian/docs/plans/YYYY-MM-DD-<topic>-design.md` |
| Implementation plans | `~/trabian/docs/plans/YYYY-MM-DD-<topic>-plan.md` |
| Knowledge base | `~/trabian/docs/<system>/` |
| Technical reviews | `<repo>/docs/tech-review.md` |

**Pattern for design doc filenames:**
```
docs/plans/2025-01-15-mcp-authentication-design.md
docs/plans/2025-01-15-cli-refactor-plan.md
```

**If unsure where docs belong, ASK the user.**

---

## Post-Brainstorming Suggestions

After brainstorming completes, suggest relevant trabian commands:

| Task Type | Suggested Commands |
|-----------|-------------------|
| Feature implementation | `/dev/start-session`, `/dev/create-plan` |
| Bug fix | `/sloan/find-tasks`, `/sloan/linear` |
| Infrastructure | `/pm/raid`, review RAID log implications |
| Documentation | `/sloan/update-docs`, `/kb/q2` (if Q2-related) |

---

## Integration with Trabian Workflows

### Link to Linear issues

If brainstorming leads to new tasks:
```bash
# Create Linear issue from brainstorm outcome
/sloan/linear create "<task title>"
```

### Start development session

If ready to implement:
```bash
/dev/start-session <issue-url-or-description>
```

### Update RAID log

If brainstorming reveals risks or blockers:
```bash
/pm/raid "<project-name>"
```

---

## Examples

```bash
/sloan/super cli add config validation     # Brainstorm for trabian-cli
/sloan/super server add new endpoint       # Brainstorm for trabian-server
/sloan/super optimize harvest integration  # Prompts for repo selection
/sloan/super                               # Shows selection menu
```

## Local Model Tips

- **Prompt tersely** - Llama responds well to direct instructions
- **Review output** - Always have Claude review before committing
- See `_shared-repo-logic.md` for model tiers and usage examples
