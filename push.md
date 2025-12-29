---
description: Push commits to remote for a repository
---

# Push Command

Push local commits to remote origin for a repository.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Pushing for: <repo-name>"

### Step 2: Check Repository Status

```bash
cd <repo-path> && git status
```

If no commits ahead, report "Nothing to push" and exit.

### Step 3: Get Current Branch

```bash
cd <repo-path> && git branch --show-current
```

### Step 4: Check If Branch Tracks Remote

```bash
cd <repo-path> && git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null
```

- If this succeeds → branch already tracks remote
- If this fails → new branch, needs `-u` flag

### Step 5: Execute Push

**If branch tracks remote:**
```bash
cd <repo-path> && git push origin <branch-name>
```

**If new branch (no upstream):**
```bash
cd <repo-path> && git push -u origin <branch-name>
```

### Step 6: Confirm Result

Report push result:
```
✓ Pushed <branch-name> to origin (<X> commits)
```

Or if new branch:
```
✓ Pushed <branch-name> to origin (new branch, upstream set)
```

---

## Examples

```bash
/push                    # Interactive selection
/push pulumi             # Fuzzy match → devops-gcp-pulumi
/push commands           # Push .claude/commands repo
```
