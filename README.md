# Claude Code Configuration

This repository is designed to be cloned directly as `~/.claude/` - the Claude Code configuration directory.

## Installation

```bash
# Clone as ~/.claude
git clone https://github.com/sloanahrens/slash-commands.git ~/.claude

# Configure
cd ~/.claude
cp config.yaml.example config.yaml
# Edit config.yaml with your workspace path and repos

# Install devbot CLI
make -C devbot install
```

Or after cloning, run `/setup-workspace` in Claude Code.

## Structure

```
~/.claude/                   # This repo
├── CLAUDE.md               # Global instructions (tracked)
├── settings.json           # Permissions + plugins (tracked)
├── config.yaml             # Your workspace config (gitignored)
├── hookify.*.local.md      # Hookify rules (tracked)
├── commands/               # Slash commands (tracked)
├── devbot/                 # CLI tool (tracked)
├── docs/                   # Patterns, templates (tracked)
│
│ # Runtime (gitignored):
├── history.jsonl
├── plugins/
├── cache/
├── notes/
└── ...
```

## Commands

Run `/list-commands` for the full list. Key commands:

| Command | Description |
|---------|-------------|
| `/super <repo>` | Brainstorming session with context |
| `/run-tests <repo>` | Lint, type-check, build, test |
| `/yes-commit <repo>` | Draft and commit changes |
| `/status [repo]` | Repository status |
| `/setup-workspace` | Initial setup |

## Configuration

Edit `config.yaml`:

```yaml
workspace: ~/code/my-workspace
repos:
  - name: my-project
    group: apps
    language: typescript
    work_dir: src           # Optional subdirectory
```

## devbot CLI

Fast parallel operations across repos:

```bash
devbot status              # Git status across all repos
devbot check <repo>        # lint, typecheck, build, test
devbot exec <repo> <cmd>   # Run command in repo directory
```

See [devbot/README.md](devbot/README.md) for full documentation.

## Requirements

- [Claude Code](https://claude.ai/code) CLI
- Git
- Go 1.23+ (for devbot)
