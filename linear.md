---
description: Interact with Linear issues and projects
---

# Linear Command (Trabian Branch)

Query, create, and manage Linear issues using trabian's Linear MCP plugin.

**Arguments**: `$ARGUMENTS` - Subcommand and options (see below)

---

## Subcommands

### `my` - Show my assigned issues

```bash
/sloan/linear my              # All assigned issues
/sloan/linear my --limit 10   # Limit results
```

Uses `mcp__plugin_linear_linear__list_issues` with `assignee="me"`.

### `search <query>` - Search issues

```bash
/sloan/linear search Plaid              # Text search
/sloan/linear search --state "In Progress"
/sloan/linear search --label "Bug"
/sloan/linear search Plaid --state Backlog
```

Uses `mcp__plugin_linear_linear__list_issues` with query filters.

### `project <name>` - Show project issues

```bash
/sloan/linear project "VACU Card"       # Search by project name
/sloan/linear project                   # Interactive selection
```

Uses `mcp__plugin_linear_linear__list_projects` to find project, then `list_issues` with project filter.

### `issue <id>` - Get issue details

```bash
/sloan/linear issue TRB-123             # Get full issue details
```

Uses `mcp__plugin_linear_linear__get_issue` with includeRelations=true.

### `create <title>` - Create new issue

```bash
/sloan/linear create "Fix token refresh bug"
/sloan/linear create "Add retry logic" --priority 2
```

Prompts for:
- Description (optional)
- Team (from trabian teams)
- Priority (1=urgent, 2=high, 3=normal, 4=low)
- State (default: Backlog)

Uses `mcp__plugin_linear_linear__create_issue`.

### `comment <issue-id> <message>` - Add comment

```bash
/sloan/linear comment TRB-123 "Fixed in commit abc123"
```

Uses `mcp__plugin_linear_linear__create_comment`.

### `update <issue-id>` - Update issue

```bash
/sloan/linear update TRB-123 --state "In Progress"
/sloan/linear update TRB-123 --priority 2
/sloan/linear update TRB-123 --assignee me
```

Uses `mcp__plugin_linear_linear__update_issue`.

---

## Available MCP Tools (Trabian Plugin)

| Tool | Purpose |
|------|---------|
| `mcp__plugin_linear_linear__list_issues` | List issues with filters (assignee, state, project, label, team) |
| `mcp__plugin_linear_linear__get_issue` | Get issue details with relations |
| `mcp__plugin_linear_linear__create_issue` | Create new issue |
| `mcp__plugin_linear_linear__update_issue` | Update issue fields |
| `mcp__plugin_linear_linear__create_comment` | Add comment to issue |
| `mcp__plugin_linear_linear__list_comments` | List comments on issue |
| `mcp__plugin_linear_linear__list_projects` | List Linear projects |
| `mcp__plugin_linear_linear__get_project` | Get project details |
| `mcp__plugin_linear_linear__list_teams` | List Linear teams |
| `mcp__plugin_linear_linear__list_users` | List workspace users |
| `mcp__plugin_linear_linear__list_issue_statuses` | List available statuses for a team |
| `mcp__plugin_linear_linear__list_issue_labels` | List available labels |
| `mcp__plugin_linear_linear__list_cycles` | List team cycles (sprints) |

---

## Priority Levels

| Value | Level |
|-------|-------|
| 1 | Urgent |
| 2 | High |
| 3 | Normal |
| 4 | Low |

---

## Integration with Trabian Workflows

### Link issues to development sessions

```bash
# Start session from Linear issue
/dev/start-session https://linear.app/trabian/issue/TRB-123

# Update issue when done
/sloan/linear update TRB-123 --state "Done"
/sloan/linear comment TRB-123 "Completed in PR #45"
```

### RAID log integration

When working on issues, consider updating RAID logs:

```bash
# If issue reveals a risk or blocker
/pm/raid "Project Name" "Discovered dependency issue from TRB-123"
```

---

## Examples

```bash
/sloan/linear my                                      # My assigned issues
/sloan/linear search "authentication" --state Todo    # Search with filter
/sloan/linear project "VACU"                          # Project issues
/sloan/linear issue TRB-123                           # Issue details
/sloan/linear create "Add error handling"             # Create issue
/sloan/linear comment TRB-123 "WIP - 50% done"        # Add comment
/sloan/linear update TRB-123 --state "In Progress"    # Update status
```
