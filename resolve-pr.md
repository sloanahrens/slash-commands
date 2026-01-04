---
description: Resolve GitHub PR review feedback with technical rigor
argument-hint: <github-pr-url>
allowed-tools: [Bash, Read, Glob, Grep, Task, Skill, WebFetch, AskUserQuestion]
---

# Resolve GitHub PR Feedback

Analyze and resolve all feedback from a GitHub pull request.

**PR URL provided:** $ARGUMENTS

## Step 1: Validate Input and Parse PR URL

Parse the provided URL to extract:
- Owner (organization or user)
- Repository name
- PR number

Expected format: `https://github.com/{owner}/{repo}/pull/{number}`

If the URL is invalid or missing, stop and explain the expected format:
```
Usage: /resolve-pr https://github.com/owner/repo/pull/123
```

## Step 2: Verify Environment

Run `pwd` to confirm current working directory.

Check that `gh` CLI is available:
```bash
gh --version
```

If not installed, explain:
```
The GitHub CLI (gh) is required. Install with: brew install gh
Then authenticate with: gh auth login
```

Check authentication:
```bash
gh auth status
```

## Step 3: Find Local Repository

Use devbot to quickly find the local repository matching the GitHub remote:

```bash
devbot find-repo {owner}/{repo}
```

This searches all configured repos in parallel (~0.03s) and returns:
- Repository name
- Local path
- Remote configuration

**If found:**
- Store the returned path for use in subsequent steps

**If no match found:**
```bash
devbot status --all    # List all known repos
```
- Use AskUserQuestion to ask user to specify the correct local path

**Important directory safety rules:**
- Always use `git -C <absolute-path>` instead of `cd`
- Use absolute paths when referencing files
- Never assume directory state - always verify

## Step 4: Fetch PR Data

Using the owner, repo, and PR number, fetch all feedback.

**IMPORTANT:** Always run `pwd` before these commands to verify location.

### 4.1 PR Details
```bash
gh pr view {number} --repo {owner}/{repo} --json title,body,state,author,baseRefName,headRefName,url
```

### 4.2 Review Comments (inline code comments)
```bash
gh api repos/{owner}/{repo}/pulls/{number}/comments --paginate
```
Fields of interest: `path`, `line`, `body`, `user.login`, `created_at`, `in_reply_to_id`, `position`, `diff_hunk`

### 4.3 PR Conversation Comments (general discussion)
```bash
gh api repos/{owner}/{repo}/issues/{number}/comments --paginate
```

### 4.4 Reviews (approve/request changes)
```bash
gh api repos/{owner}/{repo}/pulls/{number}/reviews --paginate
```
Fields of interest: `state` (APPROVED, CHANGES_REQUESTED, COMMENTED, PENDING), `body`, `user.login`

### 4.5 CI/CD Check Status
```bash
gh pr checks {number} --repo {owner}/{repo}
```

## Step 5: Read Relevant Code Files

Using the matched local repository path from Step 3:

1. Identify all files mentioned in review comments
2. Read each file to understand current state:
   ```bash
   # Example - use actual matched path
   cat "<matched-repo-path>/<file-path>"
   ```
3. Pay attention to the specific line numbers mentioned in comments

**Always use absolute paths based on the matched repository directory.**

## Step 6: Organize Feedback

Structure all feedback into a clear summary, prioritizing unresolved items:

```markdown
# PR Feedback Summary: {owner}/{repo}#{number}

## PR Context
- **Title:** {title}
- **URL:** {url}
- **Branch:** {headRefName} → {baseRefName}
- **Author:** {author}
- **State:** {state}

## Unresolved Items (Priority)

### Changes Requested Reviews
List any reviews with state=CHANGES_REQUESTED, including:
- Reviewer name
- Review summary/body
- Any specific action items

### Inline Code Comments (Unresolved)
Group by file, then by line number:
1. **`{file}:{line}`** - @{commenter}:
   > "{comment body}"

   Current code context: (show relevant lines)

### General Discussion (Unresolved)
List conversation comments that appear to need action.

## CI/CD Issues
List any failing checks with details.

## Resolved Items (Reference Only)
Briefly list resolved items for context.

## Files Requiring Changes
Summary list of all files that need modifications.
```

## Step 7: Invoke Code Review Reception

Now invoke the receiving-code-review skill to evaluate and resolve feedback:

```
/superpowers:receiving-code-review
```

Provide context:
- The organized feedback summary from Step 6
- The local repository path where changes will be made
- Current state of relevant files

The skill will guide you to:
1. **Clarify first** - Ask about any unclear items before implementing anything
2. **Verify** - Check suggestions against codebase reality before acting
3. **Evaluate** - Push back on technically incorrect suggestions with reasoning
4. **Implement in order** - Blocking issues → simple fixes → complex fixes
5. **Test each** - Verify no regressions after each change

## Step 8: Use PR Review Toolkit Agents (As Needed)

For specific types of feedback, invoke specialized `pr-review-toolkit` agents:

| Feedback Type | Agent to Invoke |
|---------------|-----------------|
| Error handling concerns | `silent-failure-hunter` - analyzes catch blocks, fallbacks |
| Test coverage gaps | `pr-test-analyzer` - identifies critical missing tests |
| Type design issues | `type-design-analyzer` - evaluates type invariants |
| Comment accuracy | `comment-analyzer` - verifies documentation |
| Code quality | `code-reviewer` - general quality analysis |

Example invocation:
```
"Launch silent-failure-hunter agent to analyze error handling in src/api/"
```

After implementing fixes, use `code-simplifier` agent to polish the code while preserving functionality.

## Error Handling

**URL parsing errors:**
- Report the specific issue with the URL format
- Show the expected format with an example

**gh CLI errors:**
- Not installed: Explain installation steps
- Not authenticated: Prompt to run `gh auth login`
- Rate limited: Report and suggest waiting
- Network errors: Report and suggest retry

**Repository not found locally:**
- List all scanned directories and their remotes
- Ask user to provide the correct path

**PR not found:**
- Report the error from `gh`
- Verify the URL is correct and user has access

**Empty feedback:**
- If no comments/reviews exist, report "No feedback to resolve on this PR"
- Offer to run a code review instead using /pr-review-toolkit:review-pr
