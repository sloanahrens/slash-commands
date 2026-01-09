---
name: suggest-devbot
enabled: true
event: bash
pattern: ^cd\s+\S+|^git\s+(status|diff|branch|log|show|fetch|remote)\b
action: warn
---

**Use devbot instead**

## For cd commands → devbot exec

| Instead of | Use |
|------------|-----|
| `cd /path/to/repo; npm run build` | `devbot exec <repo> npm run build` |
| `cd /path/to/repo/subdir; go test` | `devbot exec <repo>/subdir go test` |

`devbot exec` automatically uses `work_dir` from config.yaml.

**Fallback** (when devbot exec isn't suitable):
- `npm run build --prefix /path`
- `make -C /path target`

## For git commands → devbot wrappers

| Instead of | Use |
|------------|-----|
| `git status` | `devbot status <repo>` |
| `git diff` | `devbot diff <repo>` |
| `git branch` | `devbot branch <repo>` |
| `git log` | `devbot log <repo>` |
| `git show` | `devbot show <repo> [ref]` |
| `git fetch` | `devbot fetch <repo>` |
| `git remote -v` | `devbot remote <repo>` |

Benefits: auto-approved, faster, no path juggling.

For other git commands: `devbot path <repo>`, then `cd` + git.
