---
name: block-git-c
enabled: true
event: bash
pattern: git\s+-C\s
action: block
---

**git -C is blocked** - cannot control commits/pushes with this flag.

Use the two-step pattern:
1. `devbot path <repo>` - get path
2. `cd /path/to/repo`
3. `git <command>`

Or use devbot directly:
- `devbot status <repo>`
- `devbot diff <repo>`
- `devbot log <repo>`
- `devbot branch <repo>`
- `devbot fetch <repo>`
