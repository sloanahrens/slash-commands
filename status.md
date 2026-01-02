---
description: Show status overview of all repositories
---

# Status Command

Fast parallel status of all repositories in ~/code/.

**Arguments**: `$ARGUMENTS` - Optional repo name for single-repo view.

---

## Process

Run devbot and return results:

```bash
# All repos (shows dirty by default, clean count summarized)
devbot status

# All repos including clean
devbot status --all

# Single repo
devbot status <repo-name>
```

That's it. devbot handles discovery, parallel git status, and formatting.

---

## Examples

```bash
/sloan/status              # Quick overview (~0.05s)
/sloan/status --all        # Include clean repos
/sloan/status mango        # Single repo details
```
