---
name: suggest-devbot
enabled: true
event: bash
pattern: ^git\s+(status|diff|branch|log|show|fetch|remote)\b
action: warn
---

**Use devbot instead**

## Git commands → devbot wrappers

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

## Running commands in repos

**Preferred:** `devbot exec <repo> <command>`

```bash
devbot exec my-app npm run build      # Uses work_dir from config
devbot exec my-app/subdir go test     # Explicit subdir
```

**Fallback patterns:**
- `npm run build --prefix /path`
- `make -C /path target`

**When cd is needed:** For git commands without devbot wrappers (commit, push, etc.):
```bash
devbot path my-repo
cd /path/to/my-repo
git commit -m "message"
```
