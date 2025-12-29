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

### Step 3: Display Overview

```
Workspace Status
================

| Repo              | Branch  | Status  | Last Commit    |
|-------------------|---------|---------|----------------|
| devops-gcp-pulumi | master  | clean   | 2 hours ago    |
| fractals-nextjs   | feature | 3 dirty | 1 day ago      |
| mango             | master  | clean   | 3 days ago     |

Legend: clean = no changes, N dirty = N modified/untracked files
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
Status: fractals-nextjs
=======================

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

---

## Examples

```bash
/status              # Overview of all repos
/status fractals     # Detailed status for fractals-nextjs
```
