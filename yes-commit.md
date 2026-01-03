---
description: Commit changes for a repository
---

# Commit

Commit git changes for a repository. Shows the proposed message then proceeds to commit immediately.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, and **commit rules**.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Committing for: <repo-name>"

### Step 2: Check Repository Status

```bash
git -C <repo-path> status
```

If no changes, report "No changes to commit" and exit.

### Step 3: Review Changes

```bash
git -C <repo-path> diff --stat
```

### Step 4: Generate Commit Message

Try local model first if available (see `_shared-repo-logic.md` → "Local Model Acceleration"):

1. Get diff: `git -C <repo-path> diff --staged` (or `diff` if nothing staged)
2. Use `mcp__plugin_mlx-hub_mlx-hub__mlx_infer` with local model
3. Display with `[local]` prefix

If local model unavailable, use Claude directly.

**Message requirements:**
- Short summary (50-72 characters)
- Imperative mood ("Fix bug" not "Fixed bug")
- Focuses on WHAT and WHY, not HOW
- Follows commit rules in `_shared-repo-logic.md`

**If `--conventional` flag passed**, use Conventional Commits format:
```
<type>(<scope>): <description>

[optional body]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `build`, `ci`

### Step 5: Display and Execute

Display the proposed message, then immediately execute the commit:

```
Committing to <repo-name>:
---
<commit message>
---
```

```bash
git -C <repo-path> add -A && git -C <repo-path> commit -m "<message>"
```

The user approves via Claude Code's tool permission dialog.

### Step 6: Verify Success

```bash
git -C <repo-path> log -1 --oneline
```

If pre-commit hooks modified files, include them in an amended commit.

---

## Commit Guidelines

| Do | Don't |
|----|-------|
| Summarize nature of changes | Include Claude/Anthropic attribution |
| Keep summary under 72 chars | Include co-author lines |
| Use imperative mood | Include "Generated with" tags |
| Focus on why | Commit secrets (.env, credentials) |

---

## Worktree Handling

When committing in a worktree (`.trees/<name>`), after commit show:

```
Committed to feature/new-auth (N commits ahead of main)

Next: /push <worktree>
```

---

## Options

| Flag | Effect |
|------|--------|
| `--conventional` | Use Conventional Commits format |
| `--amend` | Amend previous commit (use with caution) |

---

## Examples

```bash
/yes-commit              # Interactive selection
/yes-commit cli          # Commit for CLI package
/yes-commit server       # Commit for server
```
