# Claude Code Configuration

This repository is designed to be cloned directly as `~/.claude/` - the Claude Code configuration directory.

## Fresh Install

```bash
# Clone as ~/.claude
git clone https://github.com/sloanahrens/slash-commands.git ~/.claude

# Configure
cd ~/.claude
cp config.yaml.example config.yaml
# Edit config.yaml with your workspace path and repos

# Install devbot CLI
make -C devbot install

# Create runtime directories
mkdir -p notes/hindsight notes/sessions
```

Or after cloning, run `/setup-workspace` in Claude Code.

## Migrating from Existing ~/.claude

If you have an existing ~/.claude with symlinks to this repo:

```bash
# 1. Save your config.yaml (it's gitignored)
cp ~/.claude/config.yaml ~/config.yaml.backup

# 2. Remove symlinks
rm ~/.claude/commands/*.md
rm ~/.claude/hookify.*.local.md
rm ~/.claude/CLAUDE.md ~/.claude/settings.json ~/.claude/config.yaml

# 3. Initialize git and pull
cd ~/.claude
git init
git remote add origin https://github.com/sloanahrens/slash-commands.git
git fetch origin
git reset --hard origin/master
git branch --set-upstream-to=origin/master master

# 4. Restore config
cp ~/config.yaml.backup ~/.claude/config.yaml

# 5. Rebuild devbot
make -C devbot install
```

## Structure

```
~/.claude/                   # This repo
├── CLAUDE.md               # Global instructions
├── settings.json           # Permissions + plugins
├── config.yaml             # Your workspace config (gitignored)
├── hookify.*.local.md      # Hookify rules
├── commands/               # Slash commands
├── devbot/                 # CLI tool
├── hooks/                  # Session hooks
├── patterns/               # Versioned patterns
├── templates/              # Prompt templates
│
│ # Runtime (gitignored):
├── notes/                  # Local hindsight/session notes
├── history.jsonl
├── plugins/
├── cache/
└── ...
```

## What's Tracked vs Gitignored

**Tracked (travels with repo):**
- `CLAUDE.md` - Global Claude instructions
- `settings.json` - Permissions and plugins
- `hookify.*.local.md` - Hookify rules
- `commands/` - Slash commands
- `devbot/` - CLI tool source
- `patterns/` - Versioned patterns
- `templates/` - Prompt templates
- `hooks/` - Session hooks
- `config.yaml.example` - Config template

**Gitignored (local/runtime):**
- `config.yaml` - Your workspace paths
- `notes/` - Local hindsight and session notes
- `history.jsonl` - Conversation history
- `plugins/` - Downloaded plugins
- `cache/`, `debug/`, `todos/`, etc. - Runtime files

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

## Syncing Changes

Since ~/.claude is a git repo:

```bash
cd ~/.claude
git pull                    # Get latest changes
git add -A && git commit    # Commit your changes
git push                    # Push to origin
```

## Updating devbot

After pulling changes that modify devbot:

```bash
make -C ~/.claude/devbot install
```

## Requirements

- [Claude Code](https://claude.ai/code) CLI
- Git
- Go 1.23+ (for devbot)
