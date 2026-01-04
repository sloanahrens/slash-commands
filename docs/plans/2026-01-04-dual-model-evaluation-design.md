# Dual-Model Evaluation for Local Model Confidence

**Date:** 2026-01-04
**Status:** Approved
**Goal:** Expand the dual-model evaluation pattern from `yes-commit` to `update-docs` and `find-tasks`, building confidence in local model capabilities over time.

---

## Background

The `yes-commit` command uses a dual-model evaluation pattern:
1. Local Qwen model generates commit message (~0.5s)
2. Claude generates independently
3. Compare against criteria (length, mood, accuracy, grammar)
4. If local passes all → use with `[local]` suffix
5. If local fails → use Claude's version

This creates a visible audit trail in git history, enabling confidence accumulation without formal metrics infrastructure. Episodic memory captures context naturally.

---

## Design

### Shared Pattern (add to `_shared-repo-logic.md`)

```markdown
## Dual-Model Evaluation

For text generation tasks, use this pattern to build confidence in local model:

### Generate Both Versions
1. Local model generates first (fast path)
2. Claude generates independently (for comparison)

### Evaluate Local Output
Apply task-specific criteria. If ALL pass:
- Use local output
- Add `[local]` marker inline
- If committing, append ` [local]` to commit message

### Criteria Templates

**For prose/docs:**
- Factually accurate
- Concise (no bloat)
- Active voice
- No hallucinations

**For task summaries:**
- Actionable verb
- Correct priority
- Under 100 chars
- Accurate references

### Inline Markers
Show provenance in output:
```
[local] Add retry logic for API failures
[claude] Refactor authentication to use OAuth2
```
```

---

### `update-docs.md` Changes

Replace Step 4 with:

#### Step 4: Update Files (Dual-Model Evaluation)

##### 4a. Generate Documentation Draft (Local Model)

For each file needing updates:

```python
mcp__plugin_mlx-hub_mlx-hub__mlx_infer(
  model_id="mlx-community/Qwen2.5-Coder-14B-Instruct-4bit",
  prompt="""Update this CLAUDE.md section based on current state.

Current section:
{existing_section}

New information:
{gathered_context}

Write the updated section. Be concise, use active voice.
Updated section:""",
  max_tokens=500
)
```

##### 4b. Claude Reviews Draft

Claude independently generates the same section, then compares.

**Evaluation criteria:**
- ✓ Factually matches gathered context (commands exist, paths correct)
- ✓ Concise (no unnecessary words)
- ✓ Active voice throughout
- ✓ No hallucinated features or commands

##### 4c. Select and Mark

**If local passes all criteria:**
- Use local draft
- Note in output: `[local] Updated: Commands section`

**If local fails any criteria:**
- Use Claude version
- Note: `[claude] Updated: Commands section (local draft had {issue})`

After all sections updated, show summary:
```
Documentation updated:
  [local] CLAUDE.md - Commands section
  [local] CLAUDE.md - Architecture section
  [claude] README.md - Quick start (local missed new flag)

Docs updated (2/3 sections via local model). Ready to commit.
```

---

### `find-tasks.md` Changes

Add to Step 5:

#### Step 5: Generate Task Options (Dual-Model Evaluation)

For each task identified (TODOs, complexity issues, coverage gaps):

##### 5a. Summarize with Local Model

```python
mcp__plugin_mlx-hub_mlx-hub__mlx_infer(
  model_id="mlx-community/Qwen2.5-Coder-14B-Instruct-4bit",
  prompt="""Write a one-line task summary. Use imperative mood, under 80 chars.

Context:
- File: {file_path}:{line}
- TODO comment: "{todo_text}"
- Surrounding code: {context}

Task summary:""",
  max_tokens=50
)
```

##### 5b. Claude Reviews Each Summary

**Evaluation criteria:**
- ✓ Starts with actionable verb (Add, Fix, Refactor, Implement)
- ✓ Under 100 characters
- ✓ Accurate file/line reference preserved
- ✓ Priority inference reasonable given context

##### 5c. Build Output with Markers

```
Tasks for: my-cli
======================

From TODO Comments:
1. [local] **Add retry logic for clone operations** (Medium)
   - Location: src/commands/clones.ts:245

2. [claude] **Refactor auth flow to support OAuth2** (High)
   - Location: src/auth/provider.ts:89
   - (local summary missed OAuth2 requirement)

From Complexity Analysis:
3. [local] **Split runStats into smaller functions** (Medium)
   - Location: cmd/main.go:793 (127 lines)
```

No commit suffix for `/find-tasks` since it doesn't commit—markers are inline only.

---

### Commit Suffix Integration

When commands end with a commit:

**Counting rule:**
- If >50% of generated content used local model → append ` [local]`
- Otherwise → no suffix (Claude-majority)

**Examples:**

`/update-docs` → updates 3 sections → 2 local, 1 claude → commits with:
```
Update CLAUDE.md commands and architecture sections [local]
```

**Chaining with `/yes-commit`:**

The commit message generation itself goes through its own dual-model evaluation.
The doc/task generation provenance is tracked separately via inline markers.

---

## Implementation Checklist

- [ ] Add "Dual-Model Evaluation" section to `_shared-repo-logic.md`
- [ ] Update `update-docs.md` Step 4 with dual-model pattern
- [ ] Update `find-tasks.md` Step 5 with dual-model pattern
- [ ] Test with real repos to validate criteria

---

## Future Opportunities

Commands that could adopt this pattern later:
- `switch.md` — context summary generation
- `resolve-pr.md` — PR feedback summarization
- `review-project.md` — findings summarization
