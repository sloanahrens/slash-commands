---
description: Update documentation for a repository
---

# Update Documentation (Trabian Branch)

Update project documentation for a repository, following trabian patterns and maintaining consistency.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Documentation Structure (Trabian)

### Workspace Level (`~/code/trabian-ai/`)

| File | Purpose |
|------|---------|
| `CLAUDE.md` | Workspace overview, structure, key commands |
| `docs/plans/` | Design docs and implementation plans |
| `docs/<system>/` | Knowledge base by system tag (q2, tecton) |

### Repository Level

| File | Purpose | Guidelines |
|------|---------|------------|
| `CLAUDE.md` | Primary Claude Code reference | Commands, patterns, warnings. 100-200 lines max. |
| `README.md` | Human entry point | Brief, link to details. Under 100 lines. |
| `docs/overview.md` | Detailed documentation | Full details, architecture. |
| `docs/tech-review.md` | Technical review findings | Only if needed. |

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Updating docs for: <repo-name>"

### Step 2: Inventory Documentation

```bash
# Check repo-level docs
ls -la <repo-path>/README.md <repo-path>/CLAUDE.md <repo-path>/docs/ 2>/dev/null

# Check for workspace-level mentions
grep -l "<repo-name>" ~/code/trabian-ai/CLAUDE.md ~/code/trabian-ai/docs/**/*.md 2>/dev/null
```

### Step 3: Gather Current State

**For TypeScript packages:**
```bash
cd <repo-path> && npm test 2>&1 | tail -10      # Test counts
cd <repo-path> && npm run build 2>&1 | tail -5  # Build status
```

**For Python MCP server:**
```bash
cd <repo-path> && uv run pytest 2>&1 | tail -10
```

**Common:**
```bash
wc -l <repo-path>/README.md <repo-path>/CLAUDE.md 2>/dev/null  # Line counts
git -C <repo-path> log --oneline -5                             # Recent changes
```

For significant updates, consider using local model to draft sections (see `_shared-repo-logic.md` → "Local Model Acceleration"). Claude reviews all drafts before writing.

### Step 4: Update Files

Follow `elements-of-style` principles: omit needless words, use active voice, be specific.

**CLAUDE.md** (Priority):
- Verify commands are current
- Ensure patterns match actual code
- Update warnings/gotchas
- Keep under 200 lines

**README.md**:
- Keep brief, include quick start
- Link to detailed docs
- Under 100 lines

**docs/overview.md** (if exists):
- Update "Last Updated" date
- Refresh architecture descriptions
- Update test/coverage metrics

### Step 5: Check Consistency

Verify documentation consistency across:
- Repo CLAUDE.md ↔ README.md
- Repo docs ↔ workspace CLAUDE.md references
- Commands documented ↔ commands that exist

### Step 6: Update Workspace Docs (if needed)

If repo changes affect workspace documentation:

```bash
# Check if workspace CLAUDE.md needs updates
cat ~/code/trabian-ai/CLAUDE.md | grep -A5 "<repo-name>"
```

Only update workspace CLAUDE.md if:
- Repo structure changed significantly
- New key commands added
- Critical warnings need workspace visibility

---

## Anti-Patterns

| DON'T | WHY |
|-------|-----|
| Duplicate metrics across files | Creates maintenance burden |
| Create README files in test dirs | Unnecessary clutter |
| Add detailed change history to CLAUDE.md | Use git log instead |
| Include volatile data in README.md | Gets stale quickly |
| Create standalone docs at workspace root | Keep docs with their repos |

---

## Trabian-Specific Patterns

### Knowledge Base Docs

For system documentation (Q2, Tecton, etc.):
- Location: `~/code/trabian-ai/docs/<system>/`
- Use system tags consistently
- Cross-reference with clone repos

### Plan Documents

When updating leads to new plans:
```
~/code/trabian-ai/docs/plans/YYYY-MM-DD-<topic>-design.md
~/code/trabian-ai/docs/plans/YYYY-MM-DD-<topic>-plan.md
```

### MCP Server Docs

For trabian-server, ensure docs reflect:
- Sub-server structure (github.py, harvest.py, etc.)
- Authentication middleware
- Available MCP tools with prefixes

---

## Output

Report:
1. Files updated with line counts
2. Summary of changes made
3. Any inconsistencies found and fixed
4. Suggestions for further documentation improvements

---

## Examples

```bash
/sloan/update-docs              # Interactive selection
/sloan/update-docs cli          # Update trabian-cli docs
/sloan/update-docs server       # Update trabian-server docs
/sloan/update-docs my-app       # Update app repo docs
```
