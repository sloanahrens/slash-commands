# Claude Code Slash Commands

Portable slash commands for managing multi-repo workspaces with Claude Code.

## Installation

```bash
git clone https://github.com/sloanahrens/slash-commands.git ~/code/slash-commands
/setup-workspace  # In Claude Code - handles config, devbot, symlinks, plugins
```

## Commands

Run `/list-commands` for the full list. Key commands:

| Command | Description |
|---------|-------------|
| `/super <repo>` | Brainstorming session with context |
| `/run-tests <repo>` | Lint, type-check, build, test |
| `/yes-commit <repo>` | Draft and commit changes |
| `/push <repo>` | Push to origin |
| `/status [repo]` | Repository status |
| `/prime <repo>` | Load patterns and notes before work |
| `/capture-hindsight` | Save failure as a note |
| `/setup-workspace` | Full setup (config, devbot, plugins) |

All repo commands require exact names from `config.yaml`.

## Configuration

```yaml
workspace: ~/code/my-workspace
repos:
  - name: my-project        # Must match directory name
    group: apps
    language: typescript
    work_dir: cmd/api       # Optional subdirectory
```

## devbot CLI

Fast parallel operations. See [devbot/README.md](devbot/README.md) for full documentation.

```bash
devbot status              # Git status across all repos
devbot check <repo>        # lint, typecheck, build, test
devbot exec <repo> <cmd>   # Run command in repo directory
```

## Requirements

- [Claude Code](https://claude.ai/code) CLI
- Git
- Go 1.23+ (for devbot)
