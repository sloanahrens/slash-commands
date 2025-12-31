---
description: Commit changes for a repository
---

# Commit (Trabian Branch)

Help commit git changes for a repository following trabian conventions.

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

### Step 3: Detect Pre-commit Hooks

Check for hooks that will run on commit:

```bash
ls <repo-path>/.husky/pre-commit 2>/dev/null
ls <repo-path>/.git/hooks/pre-commit 2>/dev/null
cat <repo-path>/package.json | grep -q "husky\|lint-staged"
```

If hooks exist, warn:
```
Pre-commit hooks detected (husky/lint-staged)
Hooks will run: lint, format, tests
This may modify files or reject the commit.
```

### Step 4: Review Changes

```bash
git -C <repo-path> diff --stat
git -C <repo-path> diff
```

### Step 5: Generate Commit Message (Local Model First)

**If local model available** (see `_local-model.md`), try it first:

```bash
# Get the diff for context
DIFF=$(cd <repo-path> && git diff --staged 2>/dev/null || git diff)

# Generate with local model
mlx_lm.generate \
  --model mlx-community/DeepSeek-Coder-V2-Lite-Instruct-4bit-mlx \
  --max-tokens 100 \
  --prompt "Write a git commit message for this diff. Use imperative mood, under 72 chars:

$DIFF

Commit message:"
```

**Display with label:**
```
[local] Proposed commit message:
---
<message from local model>
---

(y) Accept  (c) Regenerate with Claude  (e) Edit
```

**If user chooses Claude (c)**, regenerate using Claude and label `[claude]`.

**If local model unavailable**, use Claude directly (no label needed).

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

### Step 6: Present for Approval

```
Proposed commit message for <repo-name>:
---
<commit message>
---

Would you like to proceed? (yes/no/edit)
```

### Step 7: Execute Commit

If approved:
```bash
git -C <repo-path> add -A && git -C <repo-path> commit -m "<message>"
```

If user wants edits, ask what to modify and regenerate.

### Step 8: Verify Success

Confirm commit succeeded:
```bash
git -C <repo-path> log -1 --oneline
```

If pre-commit hooks modified files, include them in an amended commit.

---

## Commit Guidelines (Trabian)

| Do | Don't |
|----|-------|
| Summarize nature of changes | Include Claude/Anthropic attribution |
| Keep summary under 72 chars | Include co-author lines |
| Use imperative mood | Include "Generated with" tags |
| Focus on why | Commit secrets (.env, credentials) |

**From trabian CLAUDE.md:**
- Follow existing patterns in the codebase
- Consider financial services context for security-related changes
- Note any compliance implications in commit body if relevant

---

## Worktree Handling

When committing in a worktree (`.trees/<name>`):

1. Confirm the target branch:
   ```bash
   git -C ~/code/trabian-ai/.trees/<name> branch --show-current
   ```

2. Show commits ahead of main:
   ```bash
   git -C ~/code/trabian-ai/.trees/<name> rev-list --count main..HEAD
   ```

3. After commit, suggest next steps:
   ```
   Committed to feature/new-auth (6 commits ahead of main)

   Next steps:
     /sloan/push <worktree>        Push to remote
     /sloan/run-tests <worktree>   Verify tests pass
     Create PR when ready
   ```

---

## Example Output

```
Proposed commit message for trabian-cli:
---
Add MCP header configuration command

Implement config mcp-headers subcommand for managing
authentication headers on MCP servers. Supports add,
remove, list, and show-config operations.
---

Would you like to proceed? (yes/no/edit)
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
/sloan/yes-commit                        # Interactive selection
/sloan/yes-commit cli                    # Commit for trabian-cli
/sloan/yes-commit server --conventional  # Use conventional commits
/sloan/yes-commit auth                   # Commit for auth worktree
```
