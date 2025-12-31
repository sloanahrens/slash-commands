---
description: Show status overview of all repositories
---

# Status Command (Trabian Branch)

Display a quick overview of all configured repositories, worktrees, and clones.

**Arguments**: `$ARGUMENTS` - Optional repo name for detailed view.

---

## Process

### Step 1: Run Status Script

Execute the status check as a single command block. **Do NOT use multi-line for loops** - they fail in Claude Code's bash tool.

Read `config.yaml` for:
- `base_path` (~/trabian)
- `builtin[]` components
- `worktrees_dir` (.trees)
- `clones_config` (clones/clone-config.json)
- `repos[]` additional repos

### Step 2: Discover All Repos

**Builtin packages:**
```bash
# Check each builtin path exists
ls ~/trabian/packages/trabian-cli
ls ~/trabian/mcp/trabian-server
```

**Worktrees:**
```bash
ls -d ~/trabian/.trees/*/ 2>/dev/null
```

**Clones:**
```bash
# Read clone-config.json and check which exist
cat ~/trabian/clones/clone-config.json
ls ~/trabian/clones/
```

**Additional repos:**
Check each repo in `repos[]` exists.

### Step 3: Gather Status

Run **one repo at a time** using simple commands:

```bash
# For each repo directory that exists, run:
cd /path/to/repo && echo "repo-name: $(git branch --show-current) | $(git status --porcelain | wc -l | xargs) dirty | $(git log -1 --format='%cr')"
```

### Step 3b: Gather Linear Status (if configured)

For repos with `linear_project` in config, use Linear MCP tools:

```
mcp__plugin_linear_linear__list_issues with assignee="me" and project filter
```

Summarize as: `N In Progress` or `N Todo` or `-` if no Linear config

### Step 4: Display Overview

Present results as a simple table:

```
Trabian Workspace Status
========================

Packages:
| Repo            | Branch | Status  | Last Commit    |
|-----------------|--------|---------|----------------|
| trabian-cli     | main   | clean   | 2 hours ago    |
| trabian-server  | main   | 1 dirty | 1 day ago      |

Worktrees:
| Worktree            | Branch              | Status  | Last Commit |
|---------------------|---------------------|---------|-------------|
| feature-new-auth    | feature/new-auth    | 3 dirty | 3 hours ago |

Clones:
| Clone    | Branch | Status | Last Commit  |
|----------|--------|--------|--------------|
| q2-sdk   | main   | clean  | 5 days ago   |
| tecton   | main   | clean  | 2 weeks ago  |

Apps:
| Repo           | Branch  | Status | Last Commit | Linear        |
|----------------|---------|--------|-------------|---------------|
| client-project | feature | clean  | 1 hour ago  | 2 In Progress |

Legend: clean = no changes, N dirty = N modified/untracked files
```

### Step 5: Show Sync Status (if remote configured)

If `$ARGUMENTS` contains a repo name, show:
- Full `git status`
- Recent commits (`git log --oneline -5`)
- Any uncommitted changes

---

## Key Rules

1. **One bash call per repo** - Don't try complex loops
2. **Simple output** - Branch, dirty count, last commit time
3. **Skip missing repos** - Just note them as "not found"
4. **Group by type** - Packages, Worktrees, Clones, Apps

---

## Detailed Mode

If `$ARGUMENTS` specifies a repo, show expanded info:

```
Status: trabian-cli
===================

Path:   ~/trabian/packages/trabian-cli
Branch: main
Remote: origin/main (up to date)
Status: clean

Recent commits:
  abc1234 Add config command (2 hours ago)
  def5678 Fix clone setup (1 day ago)
  ghi9012 Update dependencies (3 days ago)

Build status:
  npm run build → last run: unknown
  npm test → last run: unknown
```

For worktrees, show parent branch info:

```
Status: feature-new-auth (worktree)
===================================

Path:   ~/trabian/.trees/feature-new-auth
Branch: feature/new-auth
Parent: main (5 commits ahead)
Status: 3 files modified

Modified:
  M src/auth/handler.ts
  M src/auth/types.ts
  ? src/auth/utils.ts
```

---

## Examples

```bash
/sloan/status              # Overview of all repos
/sloan/status cli          # Detailed status for trabian-cli
/sloan/status server       # Detailed status for trabian-server
/sloan/status q2           # Detailed status for q2-sdk clone
```
