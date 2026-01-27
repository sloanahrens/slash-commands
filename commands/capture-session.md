---
description: Capture or update a session summary for today's work
---

# Capture Session

End-of-session routine: sync Beads, infer decisions, update Linear, show summary.

**Arguments**: `$ARGUMENTS` - Repo name (optional, will ask if omitted).

---

## Process

### Step 1: Resolve Repository

```bash
devbot path <repo-name>
```

### Step 2: Verify Beads

```bash
ls /path/to/repo/.beads/ 2>/dev/null
```

**If no `.beads/`**: Ask to initialize (see CLAUDE.md). Then continue.

### Step 3: Gather Context

```bash
cd /path/to/repo
bd list --status=closed --since today
bd list --status=in_progress
git log --since="midnight" --oneline
```

Also review conversation for decisions (look for "decided to", "chose", "because", "instead of").

### Step 4: Infer and Log Decisions

**Log these** to `<repo>/.claude/decisions.md`:
- Architecture choices (why A over B)
- Constraints discovered
- Workarounds applied
- Trade-offs made

**Don't log**: Simple completions, routine commits, obvious implementations.

**Format**:
```markdown
[YYYY-MM-DD] **<brief title>**
<1-3 sentence explanation>
```

### Step 5: Sync Beads

```bash
bd sync
git push origin beads-sync 2>/dev/null
```

### Step 6: Update Linear (if configured)

Check `~/.claude/config.yaml` for `linear_projects`. If not configured, skip.

**If configured**:

1. Get beads with today's activity (closed or in_progress)
2. Trace beads → plans (parse description for plan file references)
3. Match plans → Linear issues (same logic as /prime-context Step 5)
4. For each matched issue:
   - Fetch recent comments: `list_comments(issueId)`
   - Check for already-reported work (idempotency)
   - Post new progress only:
     ```
     mcp__plugin_linear_linear__create_comment(issueId, body)
     ```

**Comment format**:
```markdown
Progress update:
• Completed <task>
• Completed <task>
• Started <task>
```

If >5 items, summarize: "Completed 4 tasks: X, Y, Z, W"

### Step 7: Show Summary

```
SESSION SUMMARY
===============

Repo: <repo-name>

Completed: ✅ <bead-id>: <title>
In progress: 🔄 <bead-id>: <title>

Decisions logged: <title> (or "none — routine implementation")
Beads: ✓ synced
Linear: 📝 XYZ-15: Posted (2 completed, 1 started)

---
Next session: /prime-context <repo-name>
```

---

## Examples

```bash
/capture-session my-frontend
/capture-session              # Asks for repo
```

---

## Related

- `/prime-context` — Load context before starting
- `/execute-plan` — Continue implementation
