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

### Step 2: Check for Uncommitted Changes

```bash
devbot status <repo-name>
```

This provides branch, dirty file count, and ahead/behind status in ~0.01s.

**If uncommitted changes exist:**

Invoke `/yes-commit <repo-name>` to commit changes first, then continue with push.

### Step 3: Check Commits Ahead

```bash
cd <repo-path> && git status
```

If no commits ahead of remote, report "Nothing to push" and exit.

### Step 4: Get Current Branch

```bash
cd <repo-path> && git branch --show-current
```

### Step 5: Check If Branch Tracks Remote

```bash
cd <repo-path> && git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null
```

- If this succeeds → branch already tracks remote
- If this fails → new branch, needs `-u` flag

### Step 6: Execute Push

**If branch tracks remote:**
```bash
cd <repo-path> && git push origin <branch-name>
```

**If new branch (no upstream):**
```bash
cd <repo-path> && git push -u origin <branch-name>
```

### Step 7: Confirm Result

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
/push pulumi             # Fuzzy match → my-infra-pulumi
/push commands           # Push .claude/commands repo
```
