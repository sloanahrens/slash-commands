---
description: Interact with Linear issues and projects
---

# Linear Command

Query, create, and manage Linear issues directly from the workspace.

**Arguments**: `$ARGUMENTS` - Subcommand and options (see below)

---

## Subcommands

### `my` - Show my assigned issues

```bash
/linear my              # All assigned issues
/linear my --limit 10   # Limit results
```

Uses `mcp__linear__get_user_issues` to fetch issues assigned to current user.

### `search <query>` - Search issues

```bash
/linear search Plaid              # Text search
/linear search --status "In Progress"
/linear search --priority 1       # Urgent only
/linear search Plaid --status Backlog
```

Uses `mcp__linear__search_issues` with filters.

### `project <repo>` - Show project issues

```bash
/linear project hanscom           # Fuzzy match repo with linear_project
/linear project                   # Interactive selection
```

Looks up `linear_project` from config.yaml and searches related issues.

### `create <title>` - Create new issue

```bash
/linear create "Fix token refresh bug"
/linear create "Add retry logic" --priority 2
```

Prompts for:
- Description (optional)
- Team (from config or prompt)
- Priority (1=urgent, 2=high, 3=normal, 4=low)
- Status (default: Backlog)

Uses `mcp__linear__create_issue`.

### `comment <issue-id> <message>` - Add comment

```bash
/linear comment MESH-905 "Fixed in commit abc123"
```

Uses `mcp__linear__add_comment`.

### `update <issue-id>` - Update issue

```bash
/linear update MESH-905 --status "In Progress"
/linear update MESH-905 --priority 2
```

Uses `mcp__linear__update_issue`.

---

## Available MCP Tools

| Tool | Purpose |
|------|---------|
| `mcp__linear__get_user_issues` | Get issues assigned to a user |
| `mcp__linear__search_issues` | Search with filters (query, status, priority, team, labels) |
| `mcp__linear__create_issue` | Create new issue (title, description, teamId, priority, status) |
| `mcp__linear__update_issue` | Update issue (id, title, description, priority, status) |
| `mcp__linear__add_comment` | Add comment to issue (issueId, body) |

---

## Priority Levels

| Value | Level |
|-------|-------|
| 1 | Urgent |
| 2 | High |
| 3 | Normal |
| 4 | Low |

---

## Configuration

Repos can specify their Linear project in `config.yaml`:

```yaml
repos:
  - name: hanscom-fcu-poc-plaid-token-manager
    linear_project: hanscom-fcu-plaid-token-manager-api
```

The workspace Linear config is in `external_sources.linear`:

```yaml
external_sources:
  linear:
    workspace: trabian
    projects:
      - hanscom-fcu-plaid-token-manager-api
```

---

## Examples

```bash
/linear my                                    # My assigned issues
/linear search "authentication" --status Todo # Search with filter
/linear project hanscom                       # Hanscom project issues
/linear create "Add error handling"           # Create issue (prompts for details)
/linear comment MESH-905 "WIP - 50% done"     # Add comment
/linear update MESH-905 --status "In Progress" # Update status
```
