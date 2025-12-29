---
description: Review Project (for specified repo, or prompts for selection)
---

# Review Project

Conduct a technical review of a repository and provide actionable recommendations.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Understand the Repository

1. Read documentation: `<repo>/CLAUDE.md`, `README.md`, `docs/overview.md`
2. Examine structure: `ls -la <repo-path>`
3. Check `package.json` for dependencies and scripts

### Step 2: Run Available Checks

```bash
npm run lint        # If available
npm run type-check  # If available
npm test            # If available
```

### Step 3: Review Key Areas

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

### Step 4: Generate Report

Write to `<repo-path>/docs/tech-review.md`:

```markdown
# Technical Review: <repo-name>
*Generated: <date>*

## Executive Summary
[3-5 key findings]

## Strengths
[What the project does well - specific examples]

## Critical Issues
[High-priority problems]
- **Issue**: Description
- **Impact**: Why this matters
- **Recommendation**: Action items
- **Effort**: Low/Medium/High

## Improvement Opportunities
[Medium-priority improvements]

## Future Considerations
[Long-term architectural improvements]
```

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
/review-project pulumi       # Fuzzy match → devops-gcp-pulumi
/review-project atap         # Fuzzy match → atap-automation2
```
