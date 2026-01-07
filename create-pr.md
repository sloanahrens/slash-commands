---
description: Create a pull request from current branch to target base branch
---

# Create PR

Create a pull request from current branch to target base branch with auto-generated title and description.

**Arguments**: `$ARGUMENTS` - `<repo> [base-branch]`. Base branch required; prompts if not provided.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, and local model usage.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Creating PR for: <repo-name>"

### Step 2: Get Base Branch

**If base branch provided in args:**
Use it directly.

**If not provided:**
Prompt user: "Which branch to PR into? (e.g., main, dev, develop)"

### Step 3: Validate State

```bash
devbot path <repo-name>
cd /path/to/repo
```

**Check current branch:**
```bash
git branch --show-current
```

**Early exits:**
- If current branch equals base branch → "Already on <base>, nothing to PR"
- If uncommitted changes exist → Suggest `/yes-commit <repo>` first, then retry

**Get commits ahead of base:**
```bash
git log <base>..HEAD --oneline
```

If no commits → "Nothing to PR, branch is up to date with <base>"

### Step 4: Gather PR Content

```bash
# Commits for title/summary generation
git log <base>..HEAD --oneline

# Diff stats for context
git diff <base> --stat
```

### Step 5: Generate PR Title (Dual-Model Evaluation)

**Note:** If local model is unavailable (see `_shared-repo-logic.md` → "Availability Check"), skip to 5b and use Claude directly.

#### 5a. Generate local model title

```python
mcp__plugin_mlx-hub_mlx-hub__mlx_infer(
  model_id="mlx-community/Qwen2.5-Coder-14B-Instruct-4bit",
  prompt="Write a PR title for these commits. Keep it under 72 chars, imperative mood, no attribution.\n\nCommits:\n{commit_list}\n\nJust output the title, nothing else.",
  max_tokens=100
)
```

Store result as `local_title`.

#### 5b. Generate Claude title

Using the commit list, generate a PR title following the requirements. Store as `claude_title`.

#### 5c. Compare and select title

Display both:
```
PR title comparison:
─────────────────────────────────────
[local]  {local_title}
[claude] {claude_title}
─────────────────────────────────────
```

**Selection criteria** (evaluate local_title):
- ✓ Under 72 characters
- ✓ Imperative mood
- ✓ Captures essence of changes
- ✓ No attribution

**If local passes:** Use `local_title`, append ` [local]` to title
**If local fails:** Use `claude_title` (no suffix)

### Step 6: Generate PR Body (Dual-Model Evaluation)

#### 6a. Generate local model body

```python
mcp__plugin_mlx-hub_mlx-hub__mlx_infer(
  model_id="mlx-community/Qwen2.5-Coder-14B-Instruct-4bit",
  prompt="Write a PR description for these changes. Use this format:\n\n## Summary\n- bullet points\n\n## Test Plan\n- [ ] verification steps\n\nCommits:\n{commit_list}\n\nDiff stats:\n{diff_stats}\n\nJust output the PR body, nothing else.",
  max_tokens=500
)
```

Store result as `local_body`.

#### 6b. Generate Claude body

Using commits and diff stats, generate PR body following the template. Store as `claude_body`.

#### 6c. Compare and select body

**Selection criteria** (evaluate local_body):
- ✓ Accurate summary (matches actual commits)
- ✓ Reasonable test plan steps
- ✓ No hallucinated changes
- ✓ Follows template structure

**If local passes:** Use `local_body`, append `\n\n[local]` at end of body
**If local fails:** Use `claude_body` (no suffix)

### Step 7: Create PR

```bash
gh pr create --base <base> --title "<title>" --body "$(cat <<'EOF'
<body>
EOF
)"
```

### Step 8: Return Result

Display:
```
✓ Created PR: <url>
  <branch> → <base>
  Title: <title>
```

---

## PR Body Template

```markdown
## Summary
- High-level change 1
- High-level change 2

## Test Plan
- [ ] Verification step 1
- [ ] Verification step 2
```

---

## Options

| Flag | Effect |
|------|--------|
| `--draft` | Create as draft PR |
| `--web` | Open PR in browser after creation |

---

## Examples

```bash
/create-pr my-service dev          # PR into dev branch
/create-pr my-service main         # PR into main branch
/create-pr my-service              # Prompts for base branch
/create-pr my-service dev --draft  # Create draft PR
```
