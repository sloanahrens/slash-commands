---
description: Show status overview of all repositories
---

# Status Command

Display a quick overview of all configured repositories.

**Arguments**: `$ARGUMENTS` - Optional repo name for detailed view.

---

## Process

### Step 1: Run Status Script

Execute the status check as a single command block. **Do NOT use multi-line for loops** - they fail in Claude Code's bash tool.

Run this **one repo at a time** using simple commands:

```bash
# For each repo directory that exists, run:
cd /path/to/repo && echo "repo-name: $(git branch --show-current) | $(git status --porcelain | wc -l | xargs) dirty | $(git log -1 --format='%cr')"
```

### Step 2: Format Output

Present results as a simple table:

```
Workspace Status
================

| Repo             | Branch  | Dirty | Last Commit    |
|------------------|---------|-------|----------------|
| devops-pulumi-ts | master  | 0     | 26 minutes ago |
| git-monitor      | main    | 0     | 44 minutes ago |
| mango            | master  | 0     | 21 hours ago   |
```

### Step 3: Detailed Mode (if repo specified)

If `$ARGUMENTS` contains a repo name, show:
- Full `git status`
- Recent commits (`git log --oneline -5`)
- Any uncommitted changes

---

## Key Rules

1. **One bash call per repo** - Don't try complex loops
2. **Simple output** - Branch, dirty count, last commit time
3. **Skip missing repos** - Just note them as "not found"
4. **No fancy features** - No Linear, no MLX, no sync calculations

---

## Examples

```bash
/status              # Overview of all repos
/status mango        # Detailed status for mango
```
