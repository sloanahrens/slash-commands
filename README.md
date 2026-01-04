# Claude Code Slash Commands

Portable slash commands for managing multi-repo workspaces with Claude Code.

## Features

- Multi-source repository discovery (builtin packages, worktrees, clones, standalone repos)
- Git worktree support for feature branch isolation
- Integration with [devbot](devbot/README.md) for fast parallel operations
- Optional Linear and GitHub MCP integration
- Optional local model acceleration via MLX

## Installation

```bash
# Clone the repo
git clone https://github.com/sloanahrens/slash-commands.git ~/code/slash-commands

# Run unified setup (in Claude Code)
/setup-workspace
```

This single command will:
1. Scan your workspace and generate `config.yaml`
2. Build and install the `devbot` CLI
3. Create symlinks in `~/.claude/commands`
4. Install recommended plugins

Each step prompts for confirmation and can be skipped.

## Commands

| Command | Description |
|---------|-------------|
| `/super <repo>` | Start brainstorming session with context |
| `/find-tasks <repo>` | Find tasks from code, Linear, GitHub |
| `/run-tests <repo>` | Run lint, type-check, build, and tests |
| `/make-test <repo>` | Test Makefile targets interactively |
| `/yes-commit <repo>` | Draft and commit changes |
| `/push <repo>` | Push commits to origin |
| `/update-docs <repo>` | Update documentation |
| `/review-project <repo>` | Technical review with analysis |
| `/resolve-pr <url>` | Resolve GitHub PR review feedback |
| `/add-repo <url>` | Clone repo (reference or working) |
| `/status [repo]` | Show repository status |
| `/sync [repo]` | Pull latest changes |
| `/switch <repo>` | Context switch with suggestions |
| `/quick-explain <code>` | Quick code explanation (local model) |
| `/quick-gen <desc>` | Quick code generation (local model) |
| `/yes-proceed` | Accept recommendation and proceed |
| `/dev-rules` | Load workspace development rules |
| `/setup-workspace` | Unified setup (config, devbot, symlinks, plugins) |
| `/setup-plugins` | Install/update plugins (standalone) |
| `/list-commands` | List all available commands |
| `/list-skills` | List available skills from plugins |

All repo commands require exact repo names from config.yaml.

## Configuration

Run `/setup-workspace` to auto-generate `config.yaml`, or create manually:

```yaml
workspace: ~/code/my-workspace

repos:
  - name: my-project        # Must match directory name exactly
    group: apps
    language: typescript

  - name: my-api
    group: apps
    language: go
    work_dir: cmd/api       # Optional: subdirectory for nested projects
```

Repo names must exactly match the directory name under `workspace`.

## devbot CLI

Fast parallel operations across your workspace. See [devbot/README.md](devbot/README.md) for full documentation.

**Commands taking repo NAME:**
```bash
devbot path <repo>         # Get full filesystem path (USE THIS FIRST!)
devbot status              # Parallel git status (~0.03s for 12 repos)
devbot status <repo>       # Single repo status
devbot check <repo>        # Auto-detected quality checks
devbot diff <repo>         # Git diff summary
devbot config <repo>       # Show config files
```

**Commands taking filesystem PATH:**
```bash
# ALWAYS get path first, then use it:
REPO_PATH=$(devbot path my-project)
devbot tree "$REPO_PATH"   # Directory structure
devbot stats "$REPO_PATH"  # Code metrics

# NEVER: devbot stats ~/code/my-project  ❌ (path may be wrong!)
```

## Worktree Workflow

```bash
git worktree add .trees/feature-name -b feature/feature-name
/switch feature-name
/run-tests feature-name
/yes-commit feature-name
/push feature-name
```

## Local Model Acceleration

Optional local MLX model for faster processing. Requires mlx-hub plugin.

| Command | Local Model Use |
|---------|-----------------|
| `/yes-commit` | Draft commit messages |
| `/quick-explain` | Code explanations |
| `/quick-gen` | Simple code generation |

## Requirements

- [Claude Code](https://claude.ai/code) CLI
- Git
- Go 1.23+ (for devbot)

## Recommended Plugins

Run `/setup-workspace` (includes plugins) or `/setup-plugins` standalone:

```bash
claude plugin marketplace add obra/superpowers-marketplace
claude plugin install superpowers@superpowers-marketplace
claude plugin install episodic-memory@superpowers-marketplace

claude plugin marketplace add anthropics/claude-plugins-official
claude plugin install code-review@claude-plugins-official
claude plugin install commit-commands@claude-plugins-official
claude plugin install pr-review-toolkit@claude-plugins-official
```

## Files

| File | Purpose |
|------|---------|
| `config.yaml` | Your workspace config (gitignored) |
| `config.yaml.example` | Template for config |
| `_shared-repo-logic.md` | Multi-source repo discovery logic |
| [`devbot/`](devbot/README.md) | Go CLI for parallel operations |
