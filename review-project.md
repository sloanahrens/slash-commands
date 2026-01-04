---
description: Conduct technical review for a repository
---

# Review Project

Conduct a technical review of a repository and update its documentation.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Reviewing: <repo-name>"

**Note:** Reference clones in `clones/` are typically not reviewed - they're external repos. If user selects a clone, confirm they want to review it.

### Step 2: Load Context

**First, get the repo path (REQUIRED):**
```bash
devbot path <repo-name>
# Output: /path/to/repo (use this literal path below)
```

1. Read repo documentation:
   - `/path/to/repo/CLAUDE.md`
   - `/path/to/repo/README.md`
   - `/path/to/repo/docs/overview.md` (if exists)

2. Examine structure:
   ```bash
   devbot tree /path/to/repo -d 2    # Takes literal PATH
   ```
   devbot tree respects .gitignore, excluding node_modules, dist, etc.

### Step 3: Analyze Codebase Metrics

Run stats analysis to identify complexity hotspots:

```bash
devbot stats /path/to/repo           # Takes literal PATH
```

**NEVER use compound commands or construct paths manually.**

**Use the output to guide review focus:**
- Large files (>500 lines) → Check for god objects, consider splitting
- Long functions (>50 lines) → Review for single responsibility
- Deep nesting (>4 levels) → Look for early returns, extraction opportunities

**Include in tech review output:**
```
## Code Metrics
- Files: 45 | Lines: 8,234 (6,102 code, 892 comments)
- Functions: 87 (avg 12 lines)
- Complexity flags: 2 large files, 3 long functions
```

### Step 4: Deep Codebase Exploration (Optional)

For thorough reviews, use `feature-dev:code-explorer` agents to analyze architecture:

```
Launch 2-3 code-explorer agents in parallel:
- "Trace the main entry points and call chains"
- "Map the architecture layers and patterns"
- "Analyze data flow and dependencies"
```

Each agent returns:
- Entry points with file:line references
- Step-by-step execution flow
- Key components and responsibilities
- Architecture insights

This provides deeper understanding than static analysis alone.

### Step 5: Run Available Checks

Use `devbot check` for auto-detected quality checks:

```bash
devbot check <repo-name>
```

This runs lint, typecheck, build, and test in the appropriate order for the detected stack.

**Manual override**: If specific commands needed, check `package.json` or `pyproject.toml`.

### Step 6: Review Key Areas

Consider using local model for initial draft sections (see `_shared-repo-logic.md` → "Local Model Acceleration").

Use the `pr-review-toolkit:code-reviewer` agent to analyze systematically.

**TypeScript Projects:**
- CLI command organization
- Error handling patterns
- Test coverage
- TypeScript strictness

**Python Projects:**
- Framework patterns
- Authentication middleware
- API client implementations
- Test coverage

**For all repos, consider:**
- Security patterns (secrets, auth)
- Production readiness

### Step 7: Update Documentation

**Primary: Update `<repo>/CLAUDE.md`**

Incorporate key findings directly:
- Update commands if they've changed
- Add warnings or gotchas discovered
- Refine architecture descriptions
- Keep it concise (100-200 lines max)

**If detailed analysis needed: `<repo>/docs/tech-review.md`**

Only create if findings are too detailed for CLAUDE.md:

```markdown
# Technical Review: <repo-name>
*Last updated: <date>*

## Summary
[3-5 key findings]

## Strengths
[What the project does well]

## Issues & Recommendations
| Priority | Issue | Impact | Recommendation |
|----------|-------|--------|----------------|
| High | ... | ... | ... |
| Medium | ... | ... | ... |

## Future Considerations
[Long-term improvements]
```

---

## Documentation Rules

| DO | DON'T |
|----|-------|
| Update `<repo>/CLAUDE.md` | Create docs at workspace root |
| Write details to `<repo>/docs/` | Duplicate info across files |
| Consider security findings | Leave issues undocumented |

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
/review-project cli          # Review CLI package
/review-project server       # Review server
/review-project my-app       # Review app repo
```
