---
description: Update documentation for a repository
---

# Update Documentation

Update project documentation for a repository, maintaining consistency across doc files.

**Arguments**: `$ARGUMENTS` - Optional repo name (supports fuzzy match). If empty, shows selection menu.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery and selection.

---

## Documentation Structure

| File | Purpose | Guidelines |
|------|---------|------------|
| `CLAUDE.md` | Primary Claude Code reference | Commands, patterns, warnings. 100-200 lines max. |
| `README.md` | Human entry point | Brief, link to details. Under 100 lines. |
| `docs/overview.md` | Detailed documentation | Full details, metrics, architecture. |

Documentation lives in each repo - not centralized in mono-claude root.

---

## Process

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`, then confirm: "Updating docs for: <repo-name>"

### Step 2: Inventory Documentation

```bash
ls -la <repo-path>/README.md <repo-path>/CLAUDE.md <repo-path>/docs/overview.md 2>/dev/null
```

### Step 3: Gather Current State

```bash
cd <repo-path> && npm test 2>&1 | tail -10      # Test counts
wc -l <repo-path>/README.md <repo-path>/CLAUDE.md  # Line counts
cd <repo-path> && git log --oneline -5          # Recent changes
```

### Step 4: Update Files

Follow `elements-of-style` principles when writing: omit needless words, use active voice, be specific.

**CLAUDE.md**: Verify commands are current, patterns accurate, links work.

**README.md**: Keep brief, include quick start, link to detailed docs.

**docs/overview.md**: Update test counts, metrics, "Last Updated" date.

### Step 5: Verify

- No metrics duplicated across files
- Line counts are reasonable
- Cross-links work

---

## Anti-Patterns

- **DON'T** duplicate metrics across multiple files
- **DON'T** create README files in test directories
- **DON'T** add detailed change history to CLAUDE.md
- **DON'T** include volatile data (test counts) in README.md

---

## Output

Report:
1. Files updated with line counts
2. Summary of changes made
3. Any files created or deleted

---

## Examples

```bash
/update-docs              # Interactive selection
/update-docs pulumi       # Fuzzy match → devops-gcp-pulumi
/update-docs atap         # Fuzzy match → atap-automation2
```
