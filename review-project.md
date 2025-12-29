---
description: Conduct technical review for a repository
---

# Review Project

Conduct a technical review of a repository and update its documentation.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Reviewing: <repo-name>"

### Step 2: Understand the Repository

1. Read documentation: `<repo>/CLAUDE.md`, `README.md`, `docs/overview.md`
2. Examine structure: `ls -la <repo-path>`
3. Check `package.json` for dependencies and scripts

### Step 3: Run Available Checks

```bash
npm run lint        # If available
npm run type-check  # If available
npm test            # If available
```

### Step 4: Review Key Areas

Use the `pr-review-toolkit:code-reviewer` agent to analyze the codebase systematically. The agent will check:

**App Repos (Next.js, React):**
- Architecture & component organization
- State management patterns
- API route design
- Error handling, test coverage

**Infrastructure Repos (Pulumi, Terraform):**
- Resource organization
- Security patterns (IAM, secrets)
- State management
- CI/CD integration

**General:**
- Code organization and discoverability
- Dependency management
- Environment configuration
- Documentation quality

Invoke the agent with the repo path and focus areas relevant to the stack.

### Step 5: Update Documentation

**Primary: Update `<repo>/CLAUDE.md`**

Incorporate key findings directly into the repo's CLAUDE.md:
- Update commands if they've changed
- Add warnings or gotchas discovered
- Refine architecture descriptions
- Keep it concise (100-200 lines max)

**If detailed analysis needed: `<repo>/docs/tech-review.md`**

Only create this file if findings are too detailed for CLAUDE.md:

```markdown
# Technical Review: <repo-name>
*Last updated: <date>*

## Summary
[3-5 key findings]

## Strengths
[What the project does well]

## Issues & Recommendations
- **Issue**: Description
- **Impact**: Why this matters
- **Fix**: Action items

## Future Considerations
[Long-term improvements]
```

**If registry info changed: Update root CLAUDE.md**

Only update `<base_path>/CLAUDE.md` if:
- Repo name or alias changed
- New gotchas that affect the registry table
- Stack or description is outdated

---

## Documentation Rules

| DO | DON'T |
|----|-------|
| Update `<repo>/CLAUDE.md` | Create docs at workspace root |
| Write details to `<repo>/docs/` | Duplicate info across files |
| Update root CLAUDE.md registry | Create standalone review files at root |

---

## Review Principles

- Be **specific and actionable** - cite file paths
- Balance **ideal** with **pragmatic**
- Acknowledge **good patterns** already in place
- Prioritize by **impact vs effort**

---

## Examples

```bash
/review-project              # Interactive selection
/review-project pulumi       # Fuzzy match → devops-pulumi-ts
/review-project atap         # Fuzzy match → atap-automation2
```
