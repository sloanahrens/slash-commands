---
description: Update documentation for a repository, or audit workspace docs
---

# Update Documentation

Update project documentation for a repository, or audit workspace-level documentation.

**Arguments**: `$ARGUMENTS` - Repo name (exact match), or empty for workspace audit. See `_shared-repo-logic.md`.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Mode Selection

**If `$ARGUMENTS` is empty**: Run **Workspace Audit** (see below)

**If `$ARGUMENTS` provided**: Run **Repository Update** (skip to "Repository Update Process")

---

## Workspace Audit (No Arguments)

When run without arguments, audit the entire workspace documentation.

### WA-1: Audit Global CLAUDE.md

Read `~/.claude/CLAUDE.md` and verify:

```bash
# Check if documented repos match config
cat ~/code/mono-claude/slash-commands/config.yaml

# Check if documented commands match actual files
ls ~/code/mono-claude/slash-commands/*.md | grep -v "^_"

# Check devbot commands still work
devbot --help 2>&1 | head -20
```

**Verify:**
- Repository registry matches `config.yaml`
- Slash commands table matches actual `.md` files
- devbot commands documented are current
- Setup instructions are accurate (test `/setup-workspace` flow)
- No obsolete repos, commands, or instructions

Report any discrepancies found.

### WA-2: Scan Repo Documentation

For each repo in config.yaml, check its CLAUDE.md (or README.md as fallback):

```bash
for repo in $(grep "name:" ~/code/mono-claude/slash-commands/config.yaml | awk '{print $3}'); do
  REPO_PATH=$(devbot path "$repo" 2>/dev/null)
  if [ -f "$REPO_PATH/CLAUDE.md" ]; then
    lines=$(wc -l < "$REPO_PATH/CLAUDE.md")
    echo "$repo: CLAUDE.md ($lines lines)"
  elif [ -f "$REPO_PATH/README.md" ]; then
    lines=$(wc -l < "$REPO_PATH/README.md")
    echo "$repo: README.md only ($lines lines) - needs CLAUDE.md"
  else
    echo "$repo: NO DOCS"
  fi
done
```

**Flag repos that need attention:**
- No CLAUDE.md or README.md
- README.md only (suggest creating CLAUDE.md from it)
- CLAUDE.md over 250 lines (needs trimming)
- CLAUDE.md under 30 lines (may be incomplete)
- Last git commit to CLAUDE.md > 30 days ago

### WA-3: Clean Workspace docs/ Folder

Check `~/code/mono-claude/docs/` for obsolete planning files:

```bash
ls -la ~/code/mono-claude/docs/
find ~/code/mono-claude/docs -name "*.md" -type f
```

**For each file found, evaluate:**

1. **Completed plans** - If the plan has been fully implemented:
   - Delete the file
   - Note: "Deleted <filename> - plan completed"

2. **Partially completed plans** - If work remains:
   - Ask user: "Keep, delete, or consolidate <filename>?"
   - If consolidate: extract remaining TODOs to relevant repo's docs/ or CLAUDE.md

3. **Reference documentation** - If it's useful ongoing reference:
   - Move to appropriate repo's `docs/` folder, OR
   - Consolidate key info into `~/.claude/CLAUDE.md` if workspace-wide

4. **Obsolete/stale** - If outdated and no longer relevant:
   - Delete the file

**Goal**: The workspace `docs/` folder should be empty or contain only active plans.

### WA-4: Report Summary

```
Workspace Documentation Audit
=============================

Global CLAUDE.md:
  ✓ Repos match config.yaml (12 repos)
  ✓ Commands match files (22 commands)
  ⚠ devops-cloud-run marked "legacy" but still in registry

Repo CLAUDE.md Status:
  ✓ 10 repos: healthy (50-200 lines)
  ⚠ mango: 245 lines (consider trimming)
  ✗ physics-stuff: missing CLAUDE.md

Workspace docs/ Cleanup:
  Deleted: 2025-12-28-documentation-improvements-design.md (completed)
  Moved: new-gcp-project-setup-plan.md → devops-pulumi-ts/docs/

Recommended actions:
  1. Trim mango/CLAUDE.md (move details to docs/overview.md)
  2. Create physics-stuff/CLAUDE.md
  3. Consider archiving devops-cloud-run if truly legacy
```

---

## Repository Update Process

When a repo name is provided.

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
devbot path <repo-name>
# Output: /path/to/repo (use this literal path below)
```

Then use devbot for fast parallel analysis (~0.03s total):

```bash
devbot tree /path/to/repo    # Directory structure (respects .gitignore)
devbot config <repo-name>    # Config files (package.json, go.mod, etc.)
devbot stats /path/to/repo   # Code metrics and complexity
```

Also check existing docs:
```bash
ls -la /path/to/repo/README.md /path/to/repo/CLAUDE.md /path/to/repo/docs/
```

**NEVER use compound commands or construct paths manually - always run `devbot path` first, then use the literal output.**

### Step 3: Gather Build/Test State

Use the literal path from Step 2 (or use `devbot check <repo-name>` for automated checks):

**For TypeScript packages:**
```bash
devbot check <repo-name>                     # Runs lint, typecheck, build, test
```

**For Python projects:**
```bash
devbot check <repo-name>                     # Runs lint, typecheck, test
```

**For recent changes:**
```bash
cd /path/to/repo
git log --oneline -5
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

**Note:** If local model is unavailable (see `_shared-repo-logic.md` → "Availability Check"), skip local model steps and use Claude directly. Omit `[local]`/`[claude]` markers in output.

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
/update-docs              # Workspace audit (global CLAUDE.md, all repos, cleanup docs/)
/update-docs fractals     # Update fractals-nextjs docs
/update-docs mango        # Update mango docs
/update-docs slash        # Update slash-commands docs
```
