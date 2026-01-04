---
description: Setup global symlinks for slash commands
---

# Setup Symlinks

Create or verify symlinks in ~/.claude/commands pointing to slash-commands.

---

## Process

1. Ensure ~/.claude/commands exists:
   ```bash
   mkdir -p ~/.claude/commands
   ```

2. Create symlinks for each command file:
   ```bash
   # Use the slash-commands directory from config
   SLASH_COMMANDS_DIR="$HOME/code/mono-claude/slash-commands"

   for file in "$SLASH_COMMANDS_DIR"/*.md "$SLASH_COMMANDS_DIR"/*.yaml; do
     name=$(basename "$file")
     # Skip files we don't want to symlink
     [[ "$name" == "config.yaml.example" ]] && continue
     target="$HOME/.claude/commands/$name"
     if [ -L "$target" ]; then
       echo "✓ $name (exists)"
     elif [ -e "$target" ]; then
       echo "⚠ $name (file exists, not a symlink)"
     else
       ln -sf "$file" "$target"
       echo "→ $name (created)"
     fi
   done
   ```

3. Report status:
   ```bash
   echo ""
   echo "Symlinks in ~/.claude/commands:"
   ls ~/.claude/commands | wc -l | xargs echo "  Total:"
   ls -la ~/.claude/commands | grep -c "slash-commands" | xargs echo "  From slash-commands:"
   ```

---

## Notes

- Uses `-sf` flag: creates symbolic link, overwrites if exists
- Only symlinks .md and .yaml files (not devbot/ directory)
- Does not touch symlinks to other directories
- Safe to run multiple times (idempotent)

---

## When to Use

- After cloning slash-commands for the first time
- After adding new command files to slash-commands
- To verify symlinks are correctly set up

---

## Output

```
Setting up symlinks in ~/.claude/commands...

✓ _shared-repo-logic.md (exists)
✓ switch.md (exists)
→ install-devbot.md (created)
→ setup-symlinks.md (created)
...

Symlinks in ~/.claude/commands:
  Total: 32
  From slash-commands: 28
```
