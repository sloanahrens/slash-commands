---
description: Pull latest changes for repositories
---

# Sync Command

Pull latest changes from remote for one or all repositories.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). If empty, syncs all repos.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Pre-flight Check

Check for dirty repos first:

```bash
devbot status --dirty
```

If any dirty repos found, warn and ask how to proceed (skip/abort/stash).

### Step 2: Sync Repos

**For all repos (parallel):**
```bash
devbot run -q -- git pull --rebase
```

**For single repo:**
```bash
devbot run -f <repo-name> -- git pull --rebase
```

This executes pulls in parallel across all repos (~0.5s vs sequential ~5s for 12 repos).
Use `-q` to suppress "Already up to date" messages.

### Step 3: Report Results

```
Sync Results
============

| Repo              | Result      | Details              |
|-------------------|-------------|----------------------|
| my-infra-pulumi   | ✓ updated   | 3 commits pulled     |
| my-nextjs-app     | ✓ current   | Already up to date   |
| my-go-api         | ⚠ skipped   | Uncommitted changes  |
| my-python-service | ✗ failed    | Merge conflict       |
```

---

## Options

| Flag | Effect |
|------|--------|
| `--all` | Sync all repos (same as no arguments) |
| `--force` | Skip dirty repo warnings, stash automatically |

---

## Error Handling

**Merge conflicts:**
```
Conflict in <repo-name>. Resolve manually:
  cd <repo-path>
  # Fix conflicts
  git rebase --continue
```

**No remote:**
```
<repo-name> has no remote configured. Skipping.
```

---

## Examples

```bash
/sync                # Sync all repos
/sync pulumi         # Sync my-infra-pulumi only
/sync --force        # Sync all, auto-stash dirty repos
```
