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

**First, get the repo path (REQUIRED):**

```bash
REPO_PATH=$(devbot path <repo-name>)
```

Then use devbot for fast parallel analysis (~0.03s total):

```bash
devbot tree "$REPO_PATH"     # Directory structure (respects .gitignore)
devbot config <repo-name>    # Config files (package.json, go.mod, etc.)
devbot stats "$REPO_PATH"    # Code metrics and complexity
```

Also check existing docs:
```bash
ls -la "$REPO_PATH"/README.md "$REPO_PATH"/CLAUDE.md "$REPO_PATH"/docs/ 2>/dev/null
```

**NEVER construct paths manually like `~/code/<repo-name>` - always use `devbot path` first.**

### Step 3: Gather Build/Test State

Use `$REPO_PATH` from Step 2 throughout:

**For TypeScript packages:**
```bash
cd "$REPO_PATH" && npm test 2>&1 | tail -10      # Test counts
cd "$REPO_PATH" && npm run build 2>&1 | tail -5  # Build status
```

**For Python projects:**
```bash
cd "$REPO_PATH" && uv run pytest 2>&1 | tail -10
```

**Common:**
```bash
git -C "$REPO_PATH" log --oneline -5         # Recent changes
```

Use stats output to update CLAUDE.md metrics section if present:
```markdown
## Codebase Metrics
- **Files:** 45 source files
- **Lines:** 8,234 total (6,102 code, 892 comments, 1,240 blank)
- **Functions:** 87 (average 12 lines)
```

For significant updates, consider using local model to draft sections (see `_shared-repo-logic.md` → "Local Model Acceleration"). Claude reviews all drafts before writing.

### Step 4: Update Files (Dual-Model Evaluation)

Use dual-model pattern from `_shared-repo-logic.md` to build confidence in local model.

#### 4a. Generate Documentation Draft (Local Model)

For each section needing updates, use local model first:

```python
mcp__plugin_mlx-hub_mlx-hub__mlx_infer(
  model_id="mlx-community/Qwen2.5-Coder-14B-Instruct-4bit",
  prompt="""Update this documentation section based on current state.

Current section:
{existing_section}

New information:
{gathered_context}

Write the updated section. Be concise, use active voice, omit needless words.
Updated section:""",
  max_tokens=500
)
```

#### 4b. Claude Reviews Draft

Claude independently generates the same section, then evaluates local draft:

**Evaluation criteria:**
- ✓ Factually matches gathered context (commands exist, paths correct)
- ✓ Concise (no unnecessary words)
- ✓ Active voice throughout
- ✓ No hallucinated features or commands

#### 4c. Select and Mark

**If local passes all criteria:**
- Use local draft
- Note in output: `[local] Updated: {section name}`

**If local fails any criteria:**
- Use Claude version
- Note: `[claude] Updated: {section name} (local had {issue})`

#### 4d. File-Specific Guidelines

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

#### 4e. Summary with Provenance

After all sections updated, show:
```
Documentation updated:
  [local] CLAUDE.md - Commands section
  [local] CLAUDE.md - Architecture section
  [claude] README.md - Quick start (local missed new flag)

Docs updated (2/3 sections via local model). Ready to commit.
```

If >50% sections used local model, commit message gets ` [local]` suffix.

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
