---
description: Remind Claude of workspace development rules
---

# Dev Rules (Trabian Branch)

Continue whatever you were doing, but remember these trabian workspace rules:

---

## Path Safety

- **Run `pwd` before bash commands** - verify current location before file/path operations
- **Use absolute paths** - always use full paths from `~/code/trabian-ai/`
- **Stay within workspace** - never navigate above `~/code/trabian-ai/`
- **Respect worktree isolation** - `.trees/` worktrees are separate git environments

---

## Trabian Structure

| Directory | Purpose |
|-----------|---------|
| `packages/` | TypeScript packages (trabian-cli) |
| `mcp/` | Python MCP servers (trabian-server) |
| `clones/` | Read-only reference repos |
| `.trees/` | Git worktrees for feature branches |
| `docs/` | Knowledge base and plans |
| `.claude/commands/` | Claude commands |

---

## File Creation

- **NO `/tmp` files** - create temporary/working files in `docs/` directories
- **Prefer editing over creating** - modify existing files when possible
- **Plans go in** `docs/plans/YYYY-MM-DD-<topic>-<type>.md`
- **Knowledge base** in `docs/<system>/` with system tags (q2, tecton, etc.)

---

## Commit Messages (Trabian)

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

## Trabian-Specific Rules

### TypeScript (packages/)
- ES2022 target, CommonJS output, strict mode
- Node.js >=18.0.0 required
- Run `npm run build` before testing

### Python (mcp/)
- Python >=3.12 required
- Use `uv` for dependency management
- Run `uv sync` before testing
- FastMCP with sub-server composition

### Financial Services Context
- Consider security and compliance requirements
- Note data sensitivity and privacy concerns
- Think about regulatory implications
- Consider production environment impact
- Evaluate client-facing impact

---

## MCP Tool Usage

When using trabian MCP tools:

| Prefix | Service |
|--------|---------|
| `mcp__plugin_linear_linear__*` | Linear issues |
| `mcp__trabian__github_*` | GitHub Projects v2 |
| `mcp__trabian__harvest_*` | Harvest time tracking |
| `mcp__trabian__hubspot_*` | HubSpot CRM |
| `mcp__trabian__projects_*` | RAID log management |
| `mcp__Notion__*` | Notion workspace |

---

## General

- Read the repo's `CLAUDE.md` before making changes
- Read `~/code/trabian-ai/CLAUDE.md` for workspace context
- Run tests after making changes
- Keep changes focused and minimal
- Consider clone repos as read-only references

---

Now continue with your previous task.
