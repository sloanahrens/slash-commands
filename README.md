# Claude Code Slash Commands

Portable slash commands for managing multi-repo workspaces with Claude Code.

## Features

- Multi-source repository discovery (builtin packages, worktrees, clones, standalone repos)
- Fuzzy matching for quick repo access
- Git worktree support for feature branch isolation
- Integration with [devbot](devbot/README.md) for fast parallel operations
- Optional Linear and GitHub MCP integration
- Optional local model acceleration via MLX

## Installation

### Option 1: Clone directly to commands folder

```bash
git clone https://github.com/sloanahrens/slash-commands.git ~/.claude/commands
cd ~/.claude/commands/devbot && make install
```

### Option 2: Clone anywhere and create symlinks

```bash
git clone https://github.com/sloanahrens/slash-commands.git ~/code/slash-commands

# Create symlinks (run in Claude Code)
/setup-symlinks

# Install devbot
/install-devbot
```

### Post-Installation

```bash
cp config.yaml.example config.yaml
# Edit config.yaml for your workspace
```

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
| `/setup-plugins` | Install recommended plugins |
| `/install-devbot` | Build and install devbot CLI |
| `/setup-symlinks` | Create global command symlinks |
| `/list-commands` | List all available commands |
| `/list-skills` | List available skills from plugins |

All repo commands support fuzzy matching (e.g., `/run-tests cli`).

## Configuration

Create `config.yaml` based on your workspace:

```yaml
base_path: ~/code/my-workspace

builtin:
  - name: my-cli
    group: packages
    path: packages/my-cli
    language: typescript

worktrees_dir: .trees
clones_config: clones/clone-config.json
code_path: ~/code

repos:
  - name: my-project
    group: projects
    aliases: [proj]
```

## Repository Types

| Type | Location | Description |
|------|----------|-------------|
| **Builtin** | `<base_path>/<path>` | Fixed packages in your monorepo |
| **Worktrees** | `<base_path>/.trees/` | Git worktrees for feature branches |
| **Clones** | `<base_path>/clones/` | Read-only reference repos |
| **Repos** | `<code_path>/<name>` | Standalone working repos |

## devbot CLI

Fast parallel operations across your workspace. See [devbot/README.md](devbot/README.md) for full documentation.

```bash
devbot status              # Parallel git status (~0.03s for 12 repos)
devbot check <repo>        # Auto-detected quality checks
devbot run -- git pull     # Parallel command execution
devbot diff <repo>         # Git diff summary
devbot stats <path>        # Code metrics
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

Run `/setup-plugins` or install manually:

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
