---
name: suggest-devbot-git
enabled: true
event: bash
pattern: ^git\s+(status|diff|branch|log|show|fetch|remote)\b
action: warn
---

**Use devbot instead of raw git**

| Instead of | Use |
|------------|-----|
| `git status` | `devbot status <repo>` |
| `git diff` | `devbot diff <repo>` |
| `git branch` | `devbot branch <repo>` |
| `git log` | `devbot log <repo>` |
| `git show <ref>` | `devbot show <repo> [ref]` |
| `git fetch` | `devbot fetch <repo>` |
| `git remote -v` | `devbot remote <repo>` |

Benefits: auto-approved, faster, no path juggling.

For other git commands: `devbot path <repo>`, then `cd` + git.
