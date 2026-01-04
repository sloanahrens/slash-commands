---
description: Push commits to remote for a repository
---

# Push Command

Push local commits to remote origin for a repository.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Pushing for: <repo-name>"

### Step 2: Check Status and Branch Info

```bash
devbot branch <repo-name>
```

This provides in a single call (~0.02s):
- Current branch name
- Tracking status (has upstream or new branch)
- Commits ahead/behind
- List of commits to push

**If uncommitted changes exist** (shown in status):

Invoke `/yes-commit <repo-name>` to commit changes first, then continue with push.

**If no commits ahead:**

Report "Nothing to push" and exit.

### Step 3: Execute Push

**First, get the repo path:**
```bash
devbot path <repo-name>
# Output: /path/to/repo (use this literal path below)
```

**If branch has upstream tracking:**
```bash
git -C /path/to/repo push origin <branch-name>
```

**If new branch (no upstream):**
```bash
git -C /path/to/repo push -u origin <branch-name>
```

### Step 4: Confirm Result

Report push result:
```
✓ Pushed <branch-name> to origin (<X> commits)
```

Or if new branch:
```
✓ Pushed <branch-name> to origin (new branch, upstream set)
```

### Step 5: Suggest Next Steps

**For feature branches (not main/master):**

Suggest creating a PR with code review:
```
Next steps:
  - Create PR: gh pr create --web
  - Or run: /pr-review-toolkit:review-pr (before creating PR)
```

The `pr-review-toolkit:code-reviewer` agent can analyze changes before creating a PR to catch issues early.

---

## Options

| Flag | Effect |
|------|--------|
| `--review` | Run code-reviewer agent before pushing |

Example:
```bash
/push my-feature --review   # Review changes then push
```

---

## Examples

```bash
/push                    # Interactive selection
/push pulumi             # Fuzzy match → my-infra-pulumi
/push commands           # Push .claude/commands repo
```
