---
description: Switch context to a repository
---

# Switch Command

Quickly switch context to a repository with status summary and suggestions.

**Arguments**: `$ARGUMENTS` - Repo name (required, exact match).

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Switching to: <repo-name>"

### Step 2: Load Context

**First, get the repo path (REQUIRED for stats/tree commands):**
```bash
devbot path <repo-name>
# Output: /path/to/repo (use this literal path below)
```

Then use devbot for fast context loading (~0.05s total):

```bash
devbot status <repo-name>     # Branch, dirty count, ahead/behind - takes NAME
devbot branch <repo-name>     # Tracking info, commits to push - takes NAME
devbot stats /path/to/repo    # Codebase metrics - takes literal PATH
```

**NEVER construct paths manually - always use `devbot path` first.**

These run in parallel and provide:
- Branch name and tracking status
- Dirty file count
- Commits ahead/behind remote
- Code metrics (files, lines, functions, complexity)

### Step 3: Display Summary

**For builtin packages:**
```
Switched to: my-cli
========================

Type:   Package (builtin)
Path:   <base_path>/packages/my-cli
Branch: main
Status: clean
Stats:  45 files, 8.2k lines, 87 functions (avg 12 lines)

Recent:
  abc1234 Add config command
  def5678 Fix clone setup
  ghi9012 Update dependencies

Quick actions:
  /run-tests my-cli     Run quality checks
  /find-tasks my-cli    Find next tasks
  /super my-cli         Start brainstorming
```

**For worktrees:**
```
Switched to: feature-new-auth (worktree)
=========================================

Type:   Worktree
Path:   <base_path>/.trees/feature-new-auth
Branch: feature/new-auth
Parent: main (5 commits ahead)
Status: 3 modified
Stats:  23 files, 4.5k lines, 42 functions (avg 15 lines)

Recent:
  abc1234 WIP: Add auth handler
  def5678 Setup auth types

Quick actions:
  /run-tests feature-new-auth   Run quality checks
  /yes-commit feature-new-auth  Commit changes
```

**For clones:**
```
Switched to: some-sdk (clone)
===========================

Type:   Clone (reference)
Path:   <base_path>/clones/some-sdk
Branch: main
Status: clean (read-only reference)

Description: SDK for external service

Quick actions:
  Search: grep -r "pattern" <base_path>/clones/some-sdk/

Note: This is a reference clone. Changes should not be committed here.
```

**For working repos:**
```
Switched to: my-project
===========================

Type:   Working Repo
Path:   <code_path>/my-project
Branch: feature/new-feature
Status: 2 modified
Stats:  156 files, 12.4k lines, 203 functions (avg 18 lines)

Recent:
  abc1234 Add feature component
  def5678 Setup project structure

Quick actions:
  /run-tests my-project     Run quality checks
  /find-tasks my-project    Find next tasks
  /super my-project         Start brainstorming
```

### Step 4: Load Context

Per `_shared-repo-logic.md` → "Context Loading":
1. Read `~/.claude/CLAUDE.md` (global settings)
2. Read `<repo-path>/CLAUDE.md` (repo-specific guidance)

Show key info from both (if they exist):
- Stack/language
- Key commands
- Any warnings or gotchas

### Step 5: Suggest Commands

Based on repo type, suggest relevant commands:

| Repo Type | Suggested Commands |
|-----------|-------------------|
| Package | `/yes-commit`, `/run-tests` |
| Worktree | `/yes-commit`, `/push` |
| Clone | Search patterns |
| Working Repo | `/find-tasks`, `/run-tests` |

---

## Examples

```bash
/switch cli          # Switch to CLI package
/switch server       # Switch to server
/switch sdk          # Switch to sdk clone
/switch auth         # Switch to worktree matching "auth"
/switch              # Show selection menu
```
