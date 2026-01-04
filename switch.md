---
description: Switch context to a repository
---

# Switch Command

Quickly switch context to a repository with status summary and suggestions.

**Arguments**: `$ARGUMENTS` - Repo name (required, supports fuzzy match).

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Read Config

**CRITICAL**: First read the actual `config.yaml` to get real path values:
```bash
cat ~/.claude/commands/config.yaml
```

Extract:
- `base_path`: The monorepo/workspace root (e.g., `~/code/mono-claude`)
- `code_path`: Location of working repos (e.g., `~/code/mono-claude`)

**DO NOT use example values from documentation. Use the actual config values.**

### Step 2: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`:
1. Parse `config.yaml` for builtin, worktrees, clones, repos
2. Fuzzy match `$ARGUMENTS` against all sources
3. If no match, show selection menu

### Step 3: Load Context

Use devbot for fast context loading (~0.05s total):

```bash
devbot status <repo-name>    # Branch, dirty count, ahead/behind
devbot branch <repo-name>    # Tracking info, commits to push
devbot stats <repo-path>     # Codebase metrics (files, lines, functions)
```

These run in parallel and provide:
- Branch name and tracking status
- Dirty file count
- Commits ahead/behind remote
- Code metrics (files, lines, functions, complexity)

### Step 4: Display Summary

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

### Step 5: Read CLAUDE.md

Show key info from repo's CLAUDE.md (if exists):
- Stack/language
- Key commands
- Any warnings or gotchas

### Step 6: Suggest Commands

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
