---
description: List all available slash commands
---

# List Commands

List all available slash commands with their descriptions.

## Process

1. List all `.md` files in `~/.claude/commands/` (excluding files starting with `_`)
2. For each file, extract the `description` from the YAML frontmatter
3. Display as a formatted table

## Output Format

```
Available Commands
==================

| Command | Description |
|---------|-------------|
| /add-repo | Add a new repository to the workspace |
| /age-notes | Review and age old notes |
| /capture-insight | Capture a failure or lesson learned |
...
```

## Notes

- Commands starting with `_` are shared logic files, not user-invocable
- Run `/help` for Claude Code built-in commands
- Each command file contains detailed usage instructions
