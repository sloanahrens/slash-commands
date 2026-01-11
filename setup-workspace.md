---
description: Unified workspace setup (config, devbot, symlinks, plugins)
---

# Setup Workspace

Complete idempotent workspace initialization in one command.

**Arguments**: `$ARGUMENTS` - Optional flags: `--skip-config`, `--skip-plugins`, `--force`

**Note:** This command hardcodes `~/code/slash-commands` paths because it runs before devbot is installed. Other commands should use `devbot path slash-commands` per `_shared-repo-logic.md`.

---

## Quick Status Check (ALWAYS DO THIS FIRST)

Before doing anything else, check if setup is already complete:

```bash
# Single command to check all components
CONFIG_OK=$( [ -f ~/code/slash-commands/config.yaml ] && echo "yes" || echo "no" )
DEVBOT_OK=$( command -v devbot >/dev/null 2>&1 && echo "yes" || echo "no" )
SYMLINKS_OK=$( [ $(ls ~/.claude/commands/*.md 2>/dev/null | wc -l) -gt 20 ] && echo "yes" || echo "no" )

echo "Config: $CONFIG_OK | devbot: $DEVBOT_OK | Symlinks: $SYMLINKS_OK"
```

**If all three are "yes"**: Print summary and STOP. No further action needed.

```
✓ Workspace already configured

  Config:    ~/code/slash-commands/config.yaml
  devbot:    $(which devbot)
  Symlinks:  $(ls ~/.claude/commands/*.md | wc -l) commands

Nothing to do. Use --force to re-run setup.
```

**Only continue if something is missing OR --force flag is passed.**

---

## Step 1: Configure (if needed)

**Skip if**: config.yaml exists AND --skip-config not passed AND --force not passed

```bash
if [ ! -f ~/code/slash-commands/config.yaml ]; then
  echo "Creating config.yaml..."
  # Copy from example if it exists
  cp ~/code/slash-commands/config.yaml.example \
     ~/code/slash-commands/config.yaml 2>/dev/null || \
  echo "⚠ No config.yaml.example found - create config.yaml manually"
fi
```

If new config created, prompt user to edit it with their repos.

---

## Step 2: Install devbot (if needed)

**Skip if**: `command -v devbot` succeeds AND --force not passed

```bash
if ! command -v devbot >/dev/null 2>&1; then
  echo "Installing devbot..."
  make -C ~/code/slash-commands/devbot install
fi
```

---

## Step 3: Create Symlinks (if needed)

**Skip if**: 20+ symlinks exist AND --force not passed

```bash
COMMANDS_DIR=~/.claude/commands
mkdir -p "$COMMANDS_DIR"

for file in ~/code/slash-commands/*.md; do
  name=$(basename "$file")
  [ "${name:0:1}" = "_" ] && continue  # Skip _shared files
  ln -sf "$file" "$COMMANDS_DIR/$name"
done

echo "✓ $(ls "$COMMANDS_DIR"/*.md | wc -l) commands linked"
```

---

## Step 4: Install Plugins (if needed)

**Skip if**: --skip-plugins passed

Delegate to `/setup-plugins` command.

---

## Completion Summary

```
✓ Workspace setup complete!

  Config:    ~/code/slash-commands/config.yaml
  devbot:    $(which devbot)
  Symlinks:  $(ls ~/.claude/commands/*.md | wc -l) commands
  Plugins:   Run /setup-plugins if needed

Quick start:
  /status              Show all repo status
  /super <repo>        Start brainstorming
  /run-tests <repo>    Run quality checks
```

---

## Flags

| Flag | Effect |
|------|--------|
| `--skip-config` | Don't touch config.yaml |
| `--skip-plugins` | Skip plugin installation |
| `--force` | Re-run all steps even if already configured |

---

## Key Design: Check First, Act Only If Needed

This command is idempotent because:

1. **Status check happens FIRST** - before any action
2. **Each step checks its own precondition** - doesn't run if already done
3. **"Nothing to do" is the happy path** - not an error condition
4. **No complex shell scripting** - simple existence checks only
