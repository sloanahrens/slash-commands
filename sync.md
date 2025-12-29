---
description: Pull latest changes for repositories
---

# Sync Command

Pull latest changes from remote for one or all repositories.

**Arguments**: `$ARGUMENTS` - Optional repo name (fuzzy match). If empty, syncs all repos.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Determine Scope

- If `$ARGUMENTS` provided → sync single repo (fuzzy match)
- If empty → sync all repos

### Step 2: Pre-flight Check

For each repo to sync:

```bash
cd <repo-path> && git status --porcelain
```

If dirty, warn:
```
⚠️  <repo-name> has uncommitted changes. Skip? (yes/no/stash)
```

Options:
- **yes** - Skip this repo
- **no** - Abort sync
- **stash** - Stash changes, pull, pop stash

### Step 3: Sync Repos

```bash
cd <repo-path> && git pull --rebase
```

Use `--rebase` to keep history clean. If conflicts occur, report and stop.

### Step 4: Report Results

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
