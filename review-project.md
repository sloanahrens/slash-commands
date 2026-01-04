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

**Note:** Reference clones in `clones/` are typically not reviewed - they're external repos. If user selects a clone, confirm they want to review it.

### Step 2: Load Context

1. Read repo documentation:
   - `<repo>/CLAUDE.md`
   - `<repo>/README.md`
   - `<repo>/docs/overview.md` (if exists)

2. Examine structure:
   ```bash
   ls -la <repo-path>
   devbot tree <repo-path> -d 2
   ```
   devbot tree automatically respects .gitignore, excluding node_modules, dist, etc.

### Step 3: Analyze Codebase Metrics

Run stats analysis to identify complexity hotspots:

```bash
devbot stats <repo-path>
```

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

### Step 4: Run Available Checks

**For TypeScript packages:**
```bash
cd <repo-path> && npm run lint
cd <repo-path> && npm run build
cd <repo-path> && npm test
```

**For Python projects:**
```bash
cd <repo-path> && uv run ruff check .
cd <repo-path> && uv run mypy .
cd <repo-path> && uv run pytest
```

**For other repos:**
Check `package.json` or `pyproject.toml` for available commands.

### Step 5: Review Key Areas

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

### Step 6: Update Documentation

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
