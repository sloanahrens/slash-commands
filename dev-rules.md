---
description: Remind Claude of workspace development rules
---

# Dev Rules

Continue whatever you were doing, but remember these workspace rules:

---

## Path Safety

- **Run `pwd` before bash commands** - verify current location before file/path operations
- **Use absolute paths** - always use full paths from your configured workspace
- **Stay within workspace** - never navigate above your base_path or code_path
- **Respect worktree isolation** - `.trees/` worktrees are separate git environments

---

## Workspace Structure

| Directory | Purpose |
|-----------|---------|
| `packages/` | TypeScript packages |
| `apps/` | Applications |
| `clones/` | Read-only reference repos |
| `.trees/` | Git worktrees for feature branches |
| `docs/` | Documentation and plans |

---

## File Creation

- **NO `/tmp` files** - create temporary/working files in `docs/` directories
- **Prefer editing over creating** - modify existing files when possible
- **Plans go in** `docs/plans/YYYY-MM-DD-<topic>-<type>.md`

---

## Commit Messages

| Do | Don't |
|----|-------|
| Use imperative mood ("Add feature") | Include Claude/Anthropic attribution |
| Keep summary under 72 characters | Include co-author lines |
| Focus on WHY not just WHAT | Include "Generated with" tags |
| Consider compliance implications | Commit secrets (.env, credentials) |

---

## YAML Gotchas

- **Avoid colons in list items** - even inside quotes, `- echo "Service URL: foo"` becomes `{'echo "Service URL': 'foo"'}`. Use dashes instead.
- **Use `|` for multiline scripts** - avoids escaping issues
- **Validate with Python** - `python3 -c "import yaml; print(yaml.safe_load(open('file.yml')))"` reveals parsing surprises

---

## Language-Specific Rules

### TypeScript (packages/)
- ES2022 target, CommonJS output, strict mode
- Node.js >=18.0.0 required
- Run `npm run build` before testing

### Python
- Python >=3.12 required
- Use `uv` for dependency management
- Run `uv sync` before testing
- FastMCP with sub-server composition

---

## MCP Tool Usage (Optional)

If you have MCP integrations configured:

| Prefix | Service |
|--------|---------|
| `mcp__plugin_linear_linear__*` | Linear issues |
| `mcp__Notion__*` | Notion workspace |

---

## General

- Read the repo's `CLAUDE.md` before making changes
- Read workspace `CLAUDE.md` for context
- Run tests after making changes
- Keep changes focused and minimal
- Consider clone repos as read-only references

---

Now continue with your previous task.
