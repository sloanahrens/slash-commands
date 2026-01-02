# Claude Code Commands (Trabian Branch)

Portable slash commands for managing multi-repo workspaces with Claude Code.

**This is the trabian branch**, adapted for the trabian workspace with:
- Integration with trabian's Linear MCP plugin
- GitHub Projects v2 via trabian MCP tools
- RAID log management via trabian MCP
- Support for trabian's workspace structure (packages/, mcp/, clones/, .trees/)

## Setup (Trabian)

The commands are pre-configured for the trabian workspace:

1. Commands are located at `~/.claude/commands/`
2. `config.yaml` is configured for trabian structure
3. Reference clones are managed via `clones/clone-config.json`

```yaml
# config.yaml - trabian structure
base_path: ~/code/trabian-ai

builtin:
  - name: trabian-cli
    group: packages
    path: packages/trabian-cli
    language: typescript

  - name: trabian-server
    group: mcp
    path: mcp/trabian-server
    language: python

worktrees_dir: .trees
clones_config: clones/clone-config.json
repos: []
```

## Commands

| Command | Description |
|---------|-------------|
| `/sloan/super <repo>` | Start brainstorming session with trabian context |
| `/sloan/find-tasks <repo>` | Find tasks from Linear, GitHub, RAID logs |
| `/sloan/run-tests <repo>` | Run lint, type-check, build, and tests |
| `/sloan/make-test <repo>` | Test Makefile targets interactively |
| `/sloan/yes-commit <repo>` | Draft and commit following trabian conventions |
| `/sloan/push <repo>` | Push commits to origin |
| `/sloan/update-docs <repo>` | Update docs following trabian patterns |
| `/sloan/review-project <repo>` | Technical review with RAID integration |
| `/sloan/resolve-pr <url>` | Resolve GitHub PR review feedback |
| `/sloan/add-repo <url>` | Clone repo (reference or app) |
| `/sloan/status [repo]` | Show status with Linear issue counts |
| `/sloan/sync [repo]` | Pull latest changes |
| `/sloan/switch <repo>` | Context switch with trabian suggestions |
| `/sloan/linear <subcommand>` | Linear issues via trabian MCP |
| `/sloan/quick-explain <code>` | Quick code explanation using local model |
| `/sloan/quick-gen <desc>` | Quick code generation using local model |
| `/sloan/yes-proceed` | Accept recommendation and proceed |
| `/sloan/dev-rules` | Trabian workspace development rules |
| `/sloan/setup-plugins` | Install recommended plugins |
| `/sloan/list-commands` | List all available commands |
| `/sloan/list-skills` | List available skills from plugins |

All repo commands support fuzzy matching (e.g., `/sloan/run-tests cli`).

## Trabian Workspace Structure

### Repository Types

| Type | Location | Description |
|------|----------|-------------|
| **Packages** | `packages/` | TypeScript packages (trabian-cli) |
| **MCP** | `mcp/` | Python MCP servers (trabian-server) |
| **Worktrees** | `.trees/` | Git worktrees for feature branches |
| **Clones** | `clones/` | Read-only reference repos (Q2, Tecton) |
| **Apps** | `~/code/trabian-ai/<name>` | Additional project repos |

### Configuration Files

| File | Purpose |
|------|---------|
| `config.yaml` | Sloan command config (builtin, repos) |
| `clones/clone-config.json` | Trabian's reference clone definitions |

## Trabian MCP Integration

These commands integrate with trabian's MCP tools:

### Linear (via Linear Plugin)
```
mcp__plugin_linear_linear__list_issues
mcp__plugin_linear_linear__get_issue
mcp__plugin_linear_linear__create_issue
mcp__plugin_linear_linear__update_issue
mcp__plugin_linear_linear__create_comment
```

### GitHub Projects v2 (via Trabian Server)
```
mcp__trabian__github_get_assigned_issues_with_project_status
mcp__trabian__github_get_project_items
mcp__trabian__github_update_issue_status_by_number
```

### RAID Logs (via Trabian Server)
```
mcp__trabian__projects_fetch_raid_entries
mcp__trabian__projects_create_raid_entry
```

## MLX Local Model Acceleration

Commands can use a local MLX model (Qwen2.5-Coder) for faster processing of simple tasks. This requires the mlx-hub plugin.

### Supported Commands

| Command | Local Model Use |
|---------|-----------------|
| `/sloan/yes-commit` | Draft commit messages |
| `/sloan/quick-explain` | Code explanations |
| `/sloan/quick-gen` | Simple code generation |
| `/sloan/super` | Summarization during brainstorming |

### Model Configuration

**Model**: `mlx-community/Qwen2.5-Coder-14B-Instruct-4bit`
**Size**: 7.7 GB
**Speed**: ~15-50 tok/s (cold/warm)

### Output Labeling

Local model output is labeled `[qwen]`, Claude output is labeled `[claude]`:

```
[qwen] Commit message:
---
Add user validation to signup form
---

(y) Accept  (c) Regenerate with Claude  (e) Edit
```

### Fallback Behavior

If the local model fails or isn't available, commands automatically fall back to Claude.

## Documentation Patterns (Trabian)

| Type | Location |
|------|----------|
| Design docs | `~/code/trabian-ai/docs/plans/YYYY-MM-DD-<topic>-design.md` |
| Implementation plans | `~/code/trabian-ai/docs/plans/YYYY-MM-DD-<topic>-plan.md` |
| Knowledge base | `~/code/trabian-ai/docs/<system>/` |
| Technical reviews | `<repo>/docs/tech-review.md` |

## Commit Rules (Trabian)

| Do | Don't |
|----|-------|
| Use imperative mood | Include Claude/Anthropic attribution |
| Keep summary under 72 chars | Include co-author lines |
| Focus on WHY not WHAT | Include "Generated with" tags |
| Consider compliance | Commit secrets |

## Worktree Workflow

For feature development, use git worktrees in `.trees/`:

```bash
# Create worktree for feature branch
git worktree add .trees/feature-name -b feature/feature-name

# Switch to worktree
/sloan/switch feature-name

# Commands work in worktrees
/sloan/run-tests feature-name
/sloan/yes-commit feature-name
/sloan/push feature-name
```

## Related Trabian Commands

These sloan commands complement existing trabian commands:

| Trabian Command | Sloan Equivalent | Notes |
|-----------------|------------------|-------|
| `/dev/commit` | `/sloan/yes-commit` | Sloan adds worktree handling |
| `/dev/start-session` | `/sloan/super` | Different workflows |
| `/pm/raid` | - | Use directly for RAID updates |
| `/kb/q2` | - | Use directly for Q2 knowledge |

## Files

| File | Purpose |
|------|---------|
| `config.yaml` | Trabian workspace config |
| `config.yaml.example` | Template showing trabian structure |
| `_shared-repo-logic.md` | Multi-source repo discovery logic |
| `_local-model.md` | Local MLX model invocation helper |

## Requirements

- [Claude Code](https://claude.ai/code) CLI
- Git
- Node.js >=18.0.0 (for TypeScript packages)
- Python >=3.12 with uv (for MCP server)
- SSH access for Q2/Tecton repos

## Recommended Plugins

Run `/sloan/setup-plugins` to install all recommended plugins, or install manually:

### Add Marketplaces

```bash
claude plugin marketplace add obra/superpowers-marketplace
claude plugin marketplace add anthropics/claude-plugins-official
```

### MLX Local Model (Optional)

For local model acceleration (commit messages, quick explanations):

```bash
claude plugin add https://github.com/sloanahrens/mlx-hub-claude-plugin
```

### Core Plugins (Superpowers Marketplace)

```bash
claude plugin install superpowers@superpowers-marketplace          # TDD, debugging, brainstorming
claude plugin install elements-of-style@superpowers-marketplace    # Writing guidance
claude plugin install episodic-memory@superpowers-marketplace      # Memory across sessions
claude plugin install double-shot-latte@superpowers-marketplace    # Auto-continue
claude plugin install superpowers-developing-for-claude-code@superpowers-marketplace  # Plugin dev
claude plugin install superpowers-lab@superpowers-marketplace      # Experimental (vim, tmux)
claude plugin install superpowers-chrome@superpowers-marketplace   # Chrome DevTools (BETA)
```

### Official Plugins

```bash
claude plugin install frontend-design@claude-plugins-official      # React/Tailwind guidance
claude plugin install feature-dev@claude-plugins-official          # Code architect agents
claude plugin install code-review@claude-plugins-official          # Code review workflow
claude plugin install commit-commands@claude-plugins-official      # Git helpers
claude plugin install pr-review-toolkit@claude-plugins-official    # PR review
claude plugin install hookify@claude-plugins-official              # Custom hooks
claude plugin install plugin-dev@claude-plugins-official           # Plugin development
claude plugin install agent-sdk-dev@claude-plugins-official        # Agent SDK helpers
claude plugin install security-guidance@claude-plugins-official    # Security best practices
claude plugin install typescript-lsp@claude-plugins-official       # TypeScript LSP
claude plugin install gopls-lsp@claude-plugins-official            # Go LSP
```

### Key Skills

The `/sloan/super` command uses the `superpowers:brainstorming` skill. Other useful skills:

| Skill | When to Use |
|-------|-------------|
| `superpowers:brainstorming` | Before creative work, designing features |
| `superpowers:writing-plans` | Creating implementation plans |
| `superpowers:systematic-debugging` | Bug investigation (find root cause first) |
| `superpowers:test-driven-development` | Writing new code (test first) |
| `superpowers:verification-before-completion` | Before claiming work is done |

## Portability Notes

This is the **trabian branch** of the portable commands. The main branch contains the generic version that works with any workspace.

### What's Different in Trabian Branch

- Multi-source repo discovery (builtin, worktrees, clones, repos)
- Integration with trabian's clone-config.json
- Linear MCP plugin integration
- GitHub Projects v2 via trabian MCP
- RAID log integration
- Trabian documentation patterns
- Financial services context awareness
- MLX local model acceleration (mlx-hub plugin)

### Integrating Updates from Master

To integrate updates from the master branch:
```bash
cd ~/.claude/commands
git fetch origin
git rebase origin/master  # Resolve conflicts as needed
```

A backup branch is recommended before rebasing:
```bash
git branch trabian-backup-$(date +%Y%m%d)
```
