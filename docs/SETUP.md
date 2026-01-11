# Setting Up on a New System

This repository is designed to be cloned directly as `~/.claude/`.

## Fresh Install (New System)

```bash
# Clone as ~/.claude
git clone https://github.com/sloanahrens/slash-commands.git ~/.claude

# Create your config
cd ~/.claude
cp config.yaml.example config.yaml
# Edit config.yaml with your workspace path and repos

# Install devbot
make -C devbot install

# Create runtime directories
mkdir -p notes/hindsight notes/sessions
```

Then run `/setup-workspace` in Claude Code to verify and install plugins.

## Migrating from Existing ~/.claude

If you have an existing ~/.claude with symlinks to this repo:

```bash
# 1. Copy the .gitignore first
cp /path/to/slash-commands/.gitignore ~/.claude/.gitignore

# 2. Save your config.yaml (it's gitignored)
cp ~/.claude/config.yaml ~/config.yaml.backup  # if it's a symlink, copy the actual content

# 3. Remove symlinks
rm ~/.claude/commands/*.md
rm ~/.claude/hookify.*.local.md
rm ~/.claude/CLAUDE.md ~/.claude/settings.json ~/.claude/config.yaml

# 4. Restore config
cp ~/config.yaml.backup ~/.claude/config.yaml

# 5. Initialize git and pull
cd ~/.claude
git init
git remote add origin https://github.com/sloanahrens/slash-commands.git
git fetch origin
git reset --hard origin/master
git branch --set-upstream-to=origin/master master

# 6. Rebuild devbot
make -C devbot install
```

## What's Tracked vs Gitignored

**Tracked (your configuration):**
- `CLAUDE.md` - Global Claude instructions
- `settings.json` - Permissions and plugins
- `hookify.*.local.md` - Hookify rules
- `commands/` - Slash commands
- `devbot/` - CLI tool source
- `docs/` - Patterns, templates
- `config.yaml.example` - Config template

**Gitignored (local/runtime):**
- `config.yaml` - Your workspace paths (sensitive)
- `history.jsonl` - Conversation history
- `plugins/` - Downloaded plugins
- `cache/`, `debug/`, `todos/`, etc. - Runtime files
- `notes/` - Local notes

## Syncing Changes

Since ~/.claude is now a git repo:

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
