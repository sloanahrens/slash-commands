---
description: Unified workspace setup (config, devbot, symlinks, plugins)
---

# Setup Workspace

Complete idempotent workspace initialization in one command.

**Arguments**: `$ARGUMENTS` - Optional flags: `--skip-config`, `--skip-plugins`

---

## Process

Run through four steps, prompting only when work is needed. Each step can be skipped.

### Step 1: Configure Workspace

**Detect workspace:**
```bash
# Use current directory as workspace root
WORKSPACE=$(pwd)
echo "Detected workspace: $WORKSPACE"
```

**Check for existing config:**
```bash
# Config lives in slash-commands directory
SLASH_COMMANDS="$HOME/code/mono-claude/slash-commands"
CONFIG_FILE="$SLASH_COMMANDS/config.yaml"

if [ -f "$CONFIG_FILE" ]; then
  echo "Found existing config.yaml"
else
  echo "No config.yaml found"
fi
```

**If no config exists — scan and generate:**

1. Use devbot to discover repos in workspace:
   ```bash
   # List all directories with .git
   find "$WORKSPACE" -maxdepth 2 -type d -name ".git" | while read gitdir; do
     repo_path=$(dirname "$gitdir")
     repo_name=$(basename "$repo_path")
     # Skip special directories
     [[ "$repo_path" == *".trees"* ]] && continue
     [[ "$repo_path" == *"clones"* ]] && continue
     [[ "$repo_path" == *"node_modules"* ]] && continue
     echo "$repo_name"
   done
   ```

2. For each discovered repo, detect language with devbot:
   ```bash
   devbot detect "$repo_path"
   # Output: typescript, nextjs | go | python | etc.
   ```

3. Present findings to user:
   ```
   Found 12 repositories:

   | Repo                    | Stack         | Status |
   |-------------------------|---------------|--------|
   | mango                   | go, nextjs    | NEW    |
   | slash-commands          | go            | NEW    |
   | ...                     |               |        |

   Generate config.yaml? (yes / edit / skip)
   ```

4. If `yes`, generate config:
   ```yaml
   workspace: ~/code/mono-claude

   repos:
     - name: mango
       group: apps
       language: go

     - name: slash-commands
       group: tools
       language: go
   ```

**If config exists — check for new repos:**

1. Read existing repos from config
2. Scan workspace for repos not in config
3. If new repos found, offer to merge:
   ```
   Found 2 new repos not in config:
     - new-project (python)
     - experimental-app (nextjs)

   Add these to config.yaml? (yes / skip)
   ```

**Skip condition:** Pass `--skip-config` flag.

---

### Step 2: Install devbot

**Check if installed:**
```bash
if command -v devbot &> /dev/null; then
  VERSION=$(devbot --version 2>/dev/null || echo "installed")
  echo "✓ devbot already installed ($VERSION)"
  # Auto-skip, no prompt needed
else
  echo "✗ devbot not found in PATH"
  # Prompt: Install devbot? (yes / skip)
fi
```

**If installing:**
```bash
cd "$SLASH_COMMANDS/devbot" && make install
```

**Verify:**
```bash
which devbot && devbot --help | head -3
```

---

### Step 3: Create Symlinks

**Check current state:**
```bash
COMMANDS_DIR="$HOME/.claude/commands"
mkdir -p "$COMMANDS_DIR"

# Count existing symlinks from slash-commands
EXISTING=$(ls -la "$COMMANDS_DIR" 2>/dev/null | grep -c "slash-commands" || echo 0)

# Count command files in slash-commands
AVAILABLE=$(ls "$SLASH_COMMANDS"/*.md 2>/dev/null | wc -l | tr -d ' ')

if [ "$EXISTING" -eq "$AVAILABLE" ]; then
  echo "✓ $EXISTING symlinks already configured"
  # Auto-skip
else
  echo "Found $EXISTING symlinks, $AVAILABLE available"
  # Prompt: Create missing symlinks? (yes / skip)
fi
```

**If creating:**
```bash
for file in "$SLASH_COMMANDS"/*.md; do
  name=$(basename "$file")
  # Skip files we don't want to symlink
  [[ "$name" == "_"* ]] && continue  # Skip _shared-repo-logic.md etc.
  target="$COMMANDS_DIR/$name"
  if [ ! -L "$target" ]; then
    ln -sf "$file" "$target"
    echo "→ $name (created)"
  fi
done
```

---

### Step 4: Install Plugins

**Invoke `/setup-plugins`:**

This step delegates to the existing `/setup-plugins` command to keep logic DRY.

```
Invoke: /setup-plugins
```

If `--skip-plugins` flag was passed, skip this step entirely.

---

### Completion Summary

After all steps:

```
→ Workspace setup complete!

┌─────────────────────────────────────────────────┐
│ Config      │ ✓ 12 repos configured             │
│ devbot      │ ✓ Installed (v1.0.0)              │
│ Symlinks    │ ✓ 28 commands available           │
│ Plugins     │ ✓ 25 installed                    │
└─────────────────────────────────────────────────┘

Quick start:
  /status              Show all repo status
  /super <repo>        Start brainstorming
  /run-tests <repo>    Run quality checks

NOTE: Restart Claude Code to activate new plugins.
```

If any steps were skipped:
```
│ Plugins     │ ⊘ Skipped (run /setup-plugins)    │
```

---

## Options

Parse flags from `$ARGUMENTS`:

| Flag | Effect |
|------|--------|
| `--skip-config` | Skip config generation/merge step |
| `--skip-plugins` | Skip plugin installation step |

---

## Idempotency

This command is safe to run multiple times:

- **Config:** Merges new repos, never removes existing ones
- **devbot:** Checks `which devbot` before building
- **Symlinks:** Uses `-sf` flag, only creates missing ones
- **Plugins:** Delegates to `/setup-plugins` which checks installed list

---

## Examples

```bash
/setup-workspace                    # Full setup, step-by-step
/setup-workspace --skip-plugins     # Skip slow plugin installation
/setup-workspace --skip-config      # Keep existing config unchanged
```
