---
description: Initialize workspace (config, devbot, directories)
---

# Setup Workspace

Complete idempotent workspace initialization. This repo IS ~/.claude - no symlinks needed.

**Arguments**: `$ARGUMENTS` - Optional flags: `--skip-config`, `--skip-plugins`, `--force`

---

## Quick Status Check (ALWAYS DO THIS FIRST)

```bash
CONFIG_OK=$( [ -f ~/.claude/config.yaml ] && echo "yes" || echo "no" )
DEVBOT_OK=$( command -v devbot >/dev/null 2>&1 && echo "yes" || echo "no" )
NOTES_OK=$( [ -d ~/.claude/notes ] && echo "yes" || echo "no" )

echo "Config: $CONFIG_OK | devbot: $DEVBOT_OK | Notes: $NOTES_OK"
```

**If all three are "yes"**: Print summary and STOP.

```
✓ Workspace already configured

  Config:  ~/.claude/config.yaml
  devbot:  $(which devbot)
  Notes:   ~/.claude/notes/

Nothing to do. Use --force to re-run setup.
```

**Only continue if something is missing OR --force flag is passed.**

---

## Step 1: Configure (if needed)

**Skip if**: config.yaml exists AND --skip-config not passed AND --force not passed

```bash
if [ ! -f ~/.claude/config.yaml ]; then
  echo "Creating config.yaml..."
  cp ~/.claude/config.yaml.example ~/.claude/config.yaml
  echo "⚠ Edit ~/.claude/config.yaml with your workspace path and repos"
fi
```

If new config created, prompt user to edit it.

---

## Step 2: Install devbot (if needed)

**Skip if**: `command -v devbot` succeeds AND --force not passed

```bash
if ! command -v devbot >/dev/null 2>&1; then
  echo "Installing devbot..."
  make -C ~/.claude/devbot install
fi
```

---

## Step 3: Create Runtime Directories (if needed)

```bash
mkdir -p ~/.claude/notes
mkdir -p ~/.claude/notes/hindsight
mkdir -p ~/.claude/notes/sessions

echo "✓ Runtime directories ready"
```

---

## Step 4: Install Plugins (if needed)

**Skip if**: --skip-plugins passed

Delegate to `/setup-plugins` command.

---

## Completion Summary

```
✓ Workspace setup complete!

  Config:   ~/.claude/config.yaml
  devbot:   $(which devbot)
  Commands: $(ls ~/.claude/commands/*.md | wc -l) available
  Notes:    ~/.claude/notes/
  Plugins:  Run /setup-plugins if needed

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

## Key Design

This command is idempotent because:

1. **Status check happens FIRST** - before any action
2. **Each step checks its own precondition** - doesn't run if already done
3. **"Nothing to do" is the happy path** - not an error condition
4. **No symlinks** - this repo IS ~/.claude
