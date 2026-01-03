---
description: Switch context to a repository
---

# Switch Command (Trabian Branch)

Quickly switch context to a repository with status summary and trabian-specific suggestions.

**Arguments**: `$ARGUMENTS` - Repo name (required, supports fuzzy match).

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`:
1. Parse `config.yaml` for builtin, worktrees, clones, repos
2. Fuzzy match `$ARGUMENTS` against all sources
3. If no match, show selection menu

### Step 2: Load Context

Use devbot for fast status with stack detection:

```bash
devbot status <repo-name>
git -C <repo-path> log --oneline -3
```

devbot provides branch, dirty count, ahead/behind, and detected stack in ~0.03s.

### Step 3: Display Summary

**For builtin packages:**
```
Switched to: trabian-cli
========================

Type:   Package (builtin)
Path:   ~/code/trabian-ai/packages/trabian-cli
Branch: main
Status: clean

Recent:
  abc1234 Add config command
  def5678 Fix clone setup
  ghi9012 Update dependencies

Quick actions:
  /sloan/run-tests cli     Run quality checks
  /sloan/find-tasks cli    Find next tasks
  /sloan/super cli         Start brainstorming
```

**For worktrees:**
```
Switched to: feature-new-auth (worktree)
=========================================

Type:   Worktree
Path:   ~/code/trabian-ai/.trees/feature-new-auth
Branch: feature/new-auth
Parent: main (5 commits ahead)
Status: 3 modified

Recent:
  abc1234 WIP: Add auth handler
  def5678 Setup auth types

Quick actions:
  /sloan/run-tests feature-new-auth   Run quality checks
  /sloan/yes-commit feature-new-auth  Commit changes
  /dev/implement-plan                  Continue implementation
```

**For clones:**
```
Switched to: q2-sdk (clone)
===========================

Type:   Clone (reference)
Path:   ~/code/trabian-ai/clones/q2-sdk
Branch: main
Status: clean (read-only reference)

Description: Q2 SDK core banking APIs

Quick actions:
  /kb/q2                   Load Q2 knowledge base
  Search: grep -r "pattern" ~/code/trabian-ai/clones/q2-sdk/

Note: This is a reference clone. Changes should not be committed here.
```

**For apps (additional repos):**
```
Switched to: client-project
===========================

Type:   App
Path:   ~/code/trabian-ai/client-project
Branch: feature/new-feature
Status: 2 modified

Linear: 3 issues assigned (2 In Progress, 1 Todo)

Recent:
  abc1234 Add feature component
  def5678 Setup project structure

Quick actions:
  /sloan/run-tests client       Run quality checks
  /sloan/find-tasks client      Find next tasks
  /sloan/linear my              Show my Linear issues
  /pm/raid "Client Project"     Update RAID log
```

### Step 4: Read CLAUDE.md

Show key info from repo's CLAUDE.md (if exists):
- Stack/language
- Key commands
- Any warnings or gotchas

For trabian packages, also reference main CLAUDE.md:
```
Workspace context from ~/code/trabian-ai/CLAUDE.md:
- Node.js >=18.0.0 required
- TypeScript ES2022 target, strict mode
- Run `npm run build` before testing
```

### Step 5: Suggest Trabian Commands

Based on repo type, suggest relevant trabian commands:

| Repo Type | Suggested Commands |
|-----------|-------------------|
| Package (trabian-cli) | `/dev/commit`, `/sloan/run-tests` |
| MCP Server | `/sloan/run-tests`, `uv run pytest` |
| Worktree | `/dev/implement-plan`, `/sloan/yes-commit` |
| Clone (Q2) | `/kb/q2`, search patterns |
| App | `/pm/raid`, `/sloan/linear`, `/pm/meeting-prep` |

---

## Examples

```bash
/sloan/switch cli          # Switch to trabian-cli
/sloan/switch server       # Switch to trabian-server
/sloan/switch q2           # Switch to q2-sdk clone
/sloan/switch auth         # Switch to worktree matching "auth"
/sloan/switch              # Show selection menu
```
