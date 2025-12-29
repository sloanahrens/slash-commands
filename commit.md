---
description: Commit changes for a repository
---

# Commit

Help commit git changes for a repository.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, and **commit rules**.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Committing for: <repo-name>"

### Step 2: Check Repository Status

```bash
cd <repo-path> && git status
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
⚠️  Pre-commit hooks detected (husky/lint-staged)
    Hooks will run: lint, format, tests
    This may modify files or reject the commit.
```

### Step 4: Review Changes

```bash
cd <repo-path> && git diff --stat
cd <repo-path> && git diff
```

### Step 5: Generate Commit Message

Analyze changes and draft a message that:
- Has short summary (50-72 characters)
- Uses imperative mood ("Fix bug" not "Fixed bug")
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
cd <repo-path> && git add -A && git commit -m "<message>"
```

If user wants edits, ask what to modify and regenerate.

### Step 8: Verify Success

Confirm commit succeeded:
```bash
cd <repo-path> && git log -1 --oneline
```

If pre-commit hooks modified files, include them in an amended commit.

---

## Commit Guidelines

| Do | Don't |
|----|-------|
| Summarize nature of changes | Include Claude/Anthropic attribution |
| Keep summary under 72 chars | Include co-author lines |
| Use imperative mood | Commit secrets (.env, credentials) |

---

## Example Output

```
Proposed commit message for my-infra-pulumi:
---
Add custom IAM role for Cloud Run deployments

Define granular permissions for deploy service account,
replacing broad predefined roles with minimum required access.
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
/commit                       # Interactive selection
/commit pulumi                # Fuzzy match → my-infra-pulumi
/commit my-app --conventional # Use conventional commits format
```
