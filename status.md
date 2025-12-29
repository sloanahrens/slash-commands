---
description: Show status overview of all repositories
---

# Status Command

Display a quick overview of all configured repositories.

**Arguments**: `$ARGUMENTS` - Optional repo name to show detailed status for just one repo.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery.

---

## Process

### Step 1: Load Configuration

Read `config.yaml` for base path and repo list.

### Step 2: Gather Status

For each repo, run:

```bash
cd <repo-path> && git status --porcelain && git rev-parse --abbrev-ref HEAD && git log -1 --format="%cr"
```

### Step 2b: Gather Linear Status (if configured)

For repos with `linear_project` in config, use Linear MCP tools:

```
# Get open issue counts
mcp__linear__search_issues(query: "<project-keywords>", status: "In Progress")
mcp__linear__search_issues(query: "<project-keywords>", status: "Todo")
```

Summarize as: `N In Progress` or `N Todo` or `-` if no Linear config

### Step 3: Display Overview

```
Workspace Status
================

| Repo              | Branch  | Status  | Last Commit    | Linear        |
|-------------------|---------|---------|----------------|---------------|
| my-infra-pulumi   | master  | clean   | 2 hours ago    | -             |
| my-nextjs-app     | feature | 3 dirty | 1 day ago      | -             |
| my-client-project | main    | clean   | 3 hours ago    | 2 In Progress |
| my-go-api         | master  | clean   | 3 days ago     | -             |

Legend: clean = no changes, N dirty = N modified/untracked files
Linear: shows issue count if repo has linear_project in config
```

### Step 4: Show Sync Status (if remote configured)

For repos with remotes:

```bash
git rev-list --left-right --count origin/HEAD...HEAD
```

Add to output:
- `↑2` = 2 commits ahead
- `↓1` = 1 commit behind
- `↑2↓1` = diverged

---

## Detailed Mode

If `$ARGUMENTS` specifies a repo, show expanded info:

```
Status: my-nextjs-app
=====================

Branch: feature/new-canvas
Remote: origin/feature/new-canvas (↑1)
Status: 3 files modified, 1 untracked

Modified:
  M src/lib/fractalRenderer.ts
  M src/app/page.tsx
  M package.json
  ? src/lib/newFeature.ts

Recent commits:
  abc1234 Add progressive rendering (2 hours ago)
  def5678 Fix zoom bounds calculation (1 day ago)
  ghi9012 Update dependencies (3 days ago)
```

For repos with `linear_project` configured, add Linear section:

```
Status: my-client-project
=========================

Branch: main
Remote: origin/main (up to date)
Status: clean

Linear Project: my-client-api
  In Progress (1):
    └── PROJ-123: API redesign (High) - https://linear.app/my-org/issue/PROJ-123
  Backlog (2):
    ├── PROJ-456: Add retry logic (High)
    └── PROJ-789: Create payment resource

Recent commits:
  abc1234 Add token refresh endpoint (3 hours ago)
```

---

## Examples

```bash
/status              # Overview of all repos
/status my-app       # Detailed status for my-nextjs-app
```
