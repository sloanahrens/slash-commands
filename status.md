---
description: Show status overview of all repositories
---

# Status Command

Fast parallel status of all repositories.

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
/status              # Quick overview (~0.05s)
/status --all        # Include clean repos
/status my-project   # Single repo details
```
