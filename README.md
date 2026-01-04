# Claude Code Slash Commands

Portable slash commands for managing multi-repo workspaces with Claude Code.

## Installation

```bash
git clone https://github.com/sloanahrens/slash-commands.git ~/code/slash-commands
/setup-workspace  # In Claude Code - handles config, devbot, symlinks, plugins
```

## Commands

| Command | Description |
|---------|-------------|
| `/super <repo>` | Brainstorming session with context |
| `/find-tasks <repo>` | Find tasks from code, Linear, GitHub |
| `/run-tests <repo>` | Lint, type-check, build, and tests |
| `/make-test <repo>` | Test Makefile targets interactively |
| `/yes-commit <repo>` | Draft and commit changes |
| `/push <repo>` | Push commits to origin |
| `/update-docs <repo>` | Update documentation |
| `/review-project <repo>` | Technical review with analysis |
| `/resolve-pr <url>` | Resolve GitHub PR review feedback |
| `/add-repo <url>` | Clone repo |
| `/status [repo]` | Repository status |
| `/sync [repo]` | Pull latest changes |
| `/switch <repo>` | Context switch |
| `/quick-explain <code>` | Code explanation (local model) |
| `/quick-gen <desc>` | Code generation (local model) |
| `/setup-workspace` | Unified setup |
| `/list-commands` | List all commands |

All repo commands require exact repo names from config.yaml.

## Configuration

```yaml
workspace: ~/code/my-workspace

repos:
  - name: my-project        # Must match directory name
    group: apps
    language: typescript
    work_dir: cmd/api       # Optional: subdirectory
```

## devbot CLI

Fast parallel operations. See [devbot/README.md](devbot/README.md) for full documentation.

**Critical:** Commands take either repo NAME or filesystem PATH:

```bash
# NAME commands
devbot path <repo>              # Get path (USE THIS FIRST)
devbot status [repo]            # Git status (~0.03s for 12 repos)
devbot check <repo>             # Auto-detected quality checks
devbot last-commit <repo> [file] # When was repo/file last committed

# PATH commands - always get path first
devbot path my-project          # → /path/to/my-project
devbot tree /path/to/my-project # Use literal path
devbot stats /path/to/my-project
```

## Requirements

- [Claude Code](https://claude.ai/code) CLI
- Git
- Go 1.23+ (for devbot)

## Files

| File | Purpose |
|------|---------|
| `config.yaml` | Workspace config (gitignored) |
| `config.yaml.example` | Template |
| `_shared-repo-logic.md` | Shared command patterns |
| [`devbot/`](devbot/README.md) | Go CLI source |
| [`docs/`](docs/) | Global config references (copy to `~/.claude/`) |

## Global Config Setup

Copy files from `docs/` to `~/.claude/`:

```bash
cp docs/CLAUDE.md ~/.claude/
cp docs/settings.json ~/.claude/
cp docs/hookify.*.md ~/.claude/
```

These provide global Claude Code configuration, permissions, and hookify rules.
