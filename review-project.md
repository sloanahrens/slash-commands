---
description: Conduct technical review for a repository
---

# Review Project (Trabian Branch)

Conduct a technical review of a repository and update its documentation, following trabian patterns.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Reviewing: <repo-name>"

**Note:** Reference clones in `clones/` are typically not reviewed - they're external repos. If user selects a clone, confirm they want to review it.

### Step 2: Load Context

1. Read trabian workspace context:
   ```bash
   cat ~/code/trabian-ai/CLAUDE.md
   ```

2. Read repo documentation:
   - `<repo>/CLAUDE.md`
   - `<repo>/README.md`
   - `<repo>/docs/overview.md` (if exists)

3. Examine structure:
   ```bash
   ls -la <repo-path>
   tree <repo-path> -L 2 -I 'node_modules|.git|dist|__pycache__'
   ```

### Step 3: Run Available Checks

**For TypeScript packages (packages/):**
```bash
cd <repo-path> && npm run lint
cd <repo-path> && npm run build
cd <repo-path> && npm test
```

**For Python MCP server (mcp/):**
```bash
cd <repo-path> && uv run ruff check .
cd <repo-path> && uv run mypy .
cd <repo-path> && uv run pytest
```

**For app repos:**
Check `package.json` or `pyproject.toml` for available commands.

### Step 4: Review Key Areas

Consider using local model for initial draft sections (see `_shared-repo-logic.md` → "Local Model Acceleration").

Use the `pr-review-toolkit:code-reviewer` agent to analyze systematically.

**TypeScript Packages (trabian-cli):**
- CLI command organization
- Error handling patterns
- Test coverage
- TypeScript strictness

**Python MCP Server (trabian-server):**
- FastMCP patterns
- Sub-server composition
- Authentication middleware
- API client implementations
- Test coverage

**App Repos:**
- Architecture & component organization
- State management patterns
- API design
- Error handling, test coverage

**For all repos, consider:**
- Financial services compliance implications
- Security patterns (secrets, auth)
- Production readiness

### Step 5: Check RAID Log (if applicable)

For app repos with project associations:
```
mcp__trabian__projects_fetch_raid_entries with project_id
```

Note any:
- Unresolved issues related to this repo
- Outstanding risks
- Pending actions

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

**Update trabian workspace CLAUDE.md only if:**
- Major structural changes to the repo
- New commands or key information
- Critical warnings for the workspace

### Step 7: Create RAID Entries (if applicable)

If review reveals issues that should be tracked:
```
mcp__trabian__projects_create_raid_entry with:
  - type: "Issue" or "Risk"
  - project_id: <from project association>
  - title: <issue summary>
  - content: <details from review>
```

---

## Documentation Rules (Trabian)

| DO | DON'T |
|----|-------|
| Update `<repo>/CLAUDE.md` | Create docs at workspace root |
| Write details to `<repo>/docs/` | Duplicate info across files |
| Consider compliance implications | Ignore security findings |
| Create RAID entries for blockers | Leave issues undocumented |

---

## Review Principles

- Be **specific and actionable** - cite file paths
- Balance **ideal** with **pragmatic**
- Acknowledge **good patterns** already in place
- Prioritize by **impact vs effort**
- Consider **financial services context**

---

## Examples

```bash
/sloan/review-project              # Interactive selection
/sloan/review-project cli          # Review trabian-cli
/sloan/review-project server       # Review trabian-server
/sloan/review-project my-app       # Review app repo
```
