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

### Step 2: Check Repository Status and Review Changes

```bash
devbot diff <repo-name>
```

This provides in a single call (~0.02s):
- Branch name
- Staged files with addition/deletion counts
- Unstaged files with addition/deletion counts

If no changes (clean), report "No changes to commit" and exit.

### Step 3: Generate Commit Message (Dual-Model Evaluation)

**This step uses dual-model evaluation to build confidence in local model commit messages.**

#### 3a. Get the diff

```bash
git -C <repo-path> diff --staged   # If files are staged
git -C <repo-path> diff            # If nothing staged yet
```

#### 3b. Generate local model message

Use `mcp__plugin_mlx-hub_mlx-hub__mlx_infer` with local model:

```python
mcp__plugin_mlx-hub_mlx-hub__mlx_infer(
  model_id="mlx-community/Qwen2.5-Coder-14B-Instruct-4bit",
  prompt="Write a git commit message for these changes. Keep it under 72 chars, imperative mood, no attribution.\n\nChanges:\n{diff_summary}\n\nJust output the commit message, nothing else.",
  max_tokens=100
)
```

Store result as `local_message`.

#### 3c. Generate Claude message

Using the same diff, generate a commit message following the requirements below. Store as `claude_message`.

#### 3d. Compare and select

Display both for evaluation:

```
Commit message comparison:
─────────────────────────────────────
[local]  {local_message}
[claude] {claude_message}
─────────────────────────────────────
```

**Selection criteria** (evaluate local_message):
- ✓ Correct length (≤72 chars)
- ✓ Imperative mood
- ✓ Captures the essence of changes
- ✓ No attribution or co-author lines
- ✓ Grammatically correct

**If local message passes all criteria:**
- Use `local_message`
- Append ` [local]` suffix to the message

**If local message fails any criteria:**
- Use `claude_message`
- No suffix added

Report which was selected: `Selected: [local]` or `Selected: [claude]`

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

### Step 4: Display and Execute

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

### Step 5: Verify Success

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
