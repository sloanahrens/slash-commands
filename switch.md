---
description: Switch context to a repository
---

# Switch Command

Quickly switch context to a repository with status summary.

**Arguments**: `$ARGUMENTS` - Repo name (required, supports fuzzy match).

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`.

If `$ARGUMENTS` empty, show selection menu.

### Step 2: Load Context

```bash
cd <repo-path>
pwd
git status
git log --oneline -3
```

### Step 3: Display Summary

```
Switched to: my-nextjs-app
==========================

Path:   ~/code/my-workspace/my-nextjs-app
Branch: feature/new-canvas
Status: 2 modified, 1 untracked

Recent:
  abc1234 Add progressive rendering
  def5678 Fix zoom bounds calculation
  ghi9012 Update dependencies

Quick actions:
  /run-tests fractals    Run quality checks
  /find-tasks fractals   Find next tasks
  /super fractals        Start brainstorming
```

### Step 4: Read CLAUDE.md

Show key info from repo's CLAUDE.md:
- Stack/language
- Key commands
- Any warnings or gotchas

---

## Examples

```bash
/switch my-app       # Switch to my-nextjs-app
/switch pulumi       # Switch to my-infra-pulumi
/switch              # Show selection menu
```
