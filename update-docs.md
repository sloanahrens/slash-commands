---
description: Update documentation for a repository
---

# Update Documentation

Update project documentation for a repository, maintaining consistency across files.

**Arguments**: `$ARGUMENTS` - Repo name (exact match). See `_shared-repo-logic.md`.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Documentation Structure

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

### Step 2: Inventory Documentation and Structure

Use devbot for fast parallel analysis (~0.03s total):

```bash
devbot tree <repo-path>     # Directory structure (respects .gitignore)
devbot config <repo-name>   # Config files (package.json, go.mod, etc.)
devbot stats <repo-path>    # Code metrics and complexity
```

Also check existing docs:
```bash
ls -la <repo-path>/README.md <repo-path>/CLAUDE.md <repo-path>/docs/ 2>/dev/null
```

### Step 3: Gather Build/Test State

**For TypeScript packages:**
```bash
cd <repo-path> && npm test 2>&1 | tail -10      # Test counts
cd <repo-path> && npm run build 2>&1 | tail -5  # Build status
```

**For Python projects:**
```bash
cd <repo-path> && uv run pytest 2>&1 | tail -10
```

**Common:**
```bash
git -C <repo-path> log --oneline -5         # Recent changes
```

Use stats output to update CLAUDE.md metrics section if present:
```markdown
## Codebase Metrics
- **Files:** 45 source files
- **Lines:** 8,234 total (6,102 code, 892 comments, 1,240 blank)
- **Functions:** 87 (average 12 lines)
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
- Commands documented ↔ commands that exist

### Step 6: Verify Documentation Accuracy (Optional)

After updating documentation, use `pr-review-toolkit:comment-analyzer` to verify accuracy:

```
"Launch comment-analyzer agent to verify documentation in CLAUDE.md and README.md"
```

The agent checks:
- Comment accuracy vs actual code behavior
- Documentation completeness
- Potential comment rot or technical debt
- Misleading or outdated descriptions

Address any high-confidence issues before finalizing.

---

## Anti-Patterns

| DON'T | WHY |
|-------|-----|
| Duplicate metrics across files | Creates maintenance burden |
| Create README files in test dirs | Unnecessary clutter |
| Add detailed change history to CLAUDE.md | Use git log instead |
| Include volatile data in README.md | Gets stale quickly |

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
/update-docs              # Interactive selection
/update-docs cli          # Update CLI docs
/update-docs server       # Update server docs
/update-docs my-app       # Update app repo docs
```
