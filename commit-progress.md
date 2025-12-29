---
description: Commit Progress (for specified repo, or prompts for selection)
---

# Commit Progress

Help commit git changes for a repository.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, and **commit rules**.

---

## Process

### Step 1: Check Repository Status

```bash
cd <repo-path> && git status
```

If no changes, report "No changes to commit" and exit.

### Step 2: Review Changes

```bash
cd <repo-path> && git diff --stat
cd <repo-path> && git diff
```

### Step 3: Generate Commit Message

Analyze changes and draft a message that:
- Has short summary (50-72 characters)
- Uses imperative mood ("Fix bug" not "Fixed bug")
- Focuses on WHAT and WHY, not HOW
- Follows commit rules in `_shared-repo-logic.md`

### Step 4: Present for Approval

```
Proposed commit message for <repo-name>:
---
<commit message>
---

Would you like to proceed? (yes/no/edit)
```

### Step 5: Execute Commit

If approved:
```bash
cd <repo-path> && git add -A && git commit -m "<message>"
```

If user wants edits, ask what to modify and regenerate.

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
Proposed commit message for devops-gcp-pulumi:
---
Add custom IAM role for Cloud Run deployments

Define granular permissions for deploy service account,
replacing broad predefined roles with minimum required access.
---

Would you like to proceed? (yes/no/edit)
```

---

## Examples

```bash
/commit-progress              # Interactive selection
/commit-progress pulumi       # Fuzzy match → devops-gcp-pulumi
/commit-progress atap         # Fuzzy match → atap-automation2
```
