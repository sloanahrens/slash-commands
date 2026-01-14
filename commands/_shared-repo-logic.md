# Shared Repo Logic

Shared patterns for repo-targeting slash commands.

---

## Workspace Paths

The workspace root is defined in `config.yaml` as `workspace: ~/code`.

**Getting any repo path:**
```bash
devbot path <repo-name>
# Example: devbot path slash-commands → /Users/sloan/code/slash-commands
```

**Note:** `setup-workspace.md` hardcodes `~/.claude/` paths because it runs before devbot is installed.

---

## Repo Resolution

### From `@directory/` (Claude Code context)

Extract name from path → `devbot path <name>` → if not found, show suggestion

### From plain name

`devbot path <name>` → if found, use it → if not found, show suggestion

**ALWAYS use devbot for paths:**

```bash
devbot path fractals-nextjs
# Output: /Users/sloan/code/fractals-nextjs
```

**NEVER construct paths manually.**

---

## Standard Process

1. Extract repo name from `$ARGUMENTS`
2. `devbot path <name>` to get full path
3. If not found, show suggestion and ask user
4. Confirm: "Working on: <repo-name>"
5. Load context (see below)

---

## Context Loading

1. Read `~/.claude/CLAUDE.md` (global)
2. Read `<repo-path>/CLAUDE.md` (repo-specific)
3. Consider `/prime <repo>` for previous session context

**Repo .claude/ folder (gitignored):**
- `<repo-path>/.claude/project-context.md` → External links, stakeholders, decisions
- `<repo-path>/.claude/sessions/` → Session notes (one file per day)

---

## Commit Rules

- **NO** Claude/Anthropic attribution or co-author lines
- Imperative mood ("Add feature" not "Added feature")
- Keep summary under 72 characters

---

## devbot Quick Reference

See [devbot/README.md](devbot/README.md) for full documentation.

**Key pattern - path vs name:**
```bash
devbot path my-repo                    # Get path first
devbot tree /full/path/to/my-repo      # Then use literal path
```

**Run commands in repos:**
```bash
devbot exec <repo> npm test            # Uses work_dir
devbot exec <repo>/subdir go build     # Explicit subdir
devbot exec <repo>/ docker build .     # Repo root (trailing /)
```

---

## Local Model (Optional)

For simple tasks (commit messages, explanations), use mlx-hub plugin if available.

**Check availability:** `claude plugin list 2>/dev/null | grep -q mlx-hub`

**If unavailable:** Skip local model, use Claude directly.

```python
mcp__plugin_mlx-hub_mlx-hub__mlx_infer(
  model_id="mlx-community/Qwen2.5-Coder-14B-Instruct-4bit",
  prompt="...", max_tokens=200
)
```

**Markers:**
- Prefix local output: `[local] Generated: "..."`
- Commit message suffix: ` [local]` if local model was used
