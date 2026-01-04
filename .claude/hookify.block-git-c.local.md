---
name: block-git-c
enabled: true
event: bash
pattern: git\s+-C
action: block
---

**git -C is not allowed**

Use the two-step pattern instead:
1. `cd /path/to/repo`
2. `git <command>`

This ensures commands work correctly and follows workspace rules in ~/.claude/CLAUDE.md.
