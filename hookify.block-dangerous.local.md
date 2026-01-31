---
name: block-dangerous
enabled: true
event: bash
pattern: (&&|;\s*\w|(?<!\S)\$\()|git\s+-C\s|pulumi\s+(stack\s+init|login\s+(file://|--local|-l\b)|destroy|stack\s+rm)
action: block
---

**Dangerous command blocked**

## Compound Bash (blocked)

Split into separate commands:
- `cmd1 && cmd2` → run `cmd1`, then `cmd2`
- `cmd1; cmd2` → run `cmd1`, then `cmd2`
- `$(cmd)` → capture output separately

## git -C (blocked)

Use the two-step pattern:
1. `devbot path <repo>` → get path
2. `cd /path/to/repo`
3. `git <command>`

Or use devbot: `devbot status/diff/log/branch/fetch <repo>`

## Pulumi Dangerous Operations (blocked)

| Command | Risk |
|---------|------|
| `pulumi stack init` | Orphans existing infrastructure |
| `pulumi destroy` | Deletes all resources |
| `pulumi stack rm` | Loses state permanently |
| `pulumi login file://` | Abandons remote state |

**Safe workflow:**
1. `devbot pulumi <repo>` - check state first
2. `pulumi stack select <env>` - select existing stack
3. `pulumi preview` → `pulumi up`

See ~/.claude/CLAUDE.md for full guidelines.
