---
name: block-compound-bash
enabled: true
event: bash
pattern: (&&|;\s*\w|(?<!\S)\$\()
action: block
---

**Compound bash commands are not allowed**

Split into separate commands:
- Instead of `cmd1 && cmd2` → run `cmd1`, then `cmd2`
- Instead of `cmd1; cmd2` → run `cmd1`, then `cmd2`
- Instead of `$(cmd)` → capture output separately

Use devbot for git operations: `devbot status`, `devbot diff`, etc.

See ~/.claude/CLAUDE.md → "Simple bash only"
