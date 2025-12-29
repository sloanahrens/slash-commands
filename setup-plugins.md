---
description: Install all recommended plugins for the workspace
---

# Setup Plugins

Idempotently install all recommended Claude Code plugins for this workspace.

**Arguments**: `$ARGUMENTS` - Optional flags: `--check` (only show status), `--update` (update existing)

---

## Plugin Registry

### Superpowers Marketplace

| Plugin | Description |
|--------|-------------|
| `superpowers` | Core skills: TDD, debugging, brainstorming, collaboration patterns |
| `elements-of-style` | Writing guidance (Strunk & White) |
| `episodic-memory` | Persistent memory across sessions via semantic search |
| `double-shot-latte` | Auto-continue without "Would you like me to continue?" |
| `superpowers-developing-for-claude-code` | Plugin/skill development resources + 42 docs |
| `superpowers-lab` | Experimental: vim, menuconfig, REPLs via tmux |
| `superpowers-chrome` | Chrome DevTools Protocol access (BETA) |

### Official Anthropic

| Plugin | Description |
|--------|-------------|
| `frontend-design` | Avoid generic React/Tailwind aesthetics |
| `feature-dev` | Code architect, explorer, reviewer agents |
| `code-review` | Code review workflow |
| `commit-commands` | Git commit helpers (/commit, /commit-push-pr) |
| `pr-review-toolkit` | PR review workflow |
| `hookify` | Create custom hooks for Claude Code |
| `plugin-dev` | Plugin development tools |
| `agent-sdk-dev` | Agent SDK development helpers |
| `security-guidance` | Security best practices |
| `typescript-lsp` | TypeScript language server |
| `gopls-lsp` | Go language server |

---

## Process

### Step 1: Check Marketplaces

Verify both marketplaces are registered:

```bash
cat ~/.claude/plugins/known_marketplaces.json
```

**If `superpowers-marketplace` missing:**
```bash
claude plugin marketplace add obra/superpowers-marketplace
```

**If `claude-plugins-official` missing:**
```bash
claude plugin marketplace add anthropics/claude-plugins-official
```

### Step 2: Check Installed Plugins

Read current installations:

```bash
cat ~/.claude/plugins/installed_plugins.json
```

Build list of what's already installed vs what needs installing.

### Step 3: Install Missing Plugins

**If `--check` flag passed:** Only report status, don't install.

For each missing plugin, run:

```bash
# Superpowers Marketplace
claude plugin install superpowers@superpowers-marketplace
claude plugin install elements-of-style@superpowers-marketplace
claude plugin install episodic-memory@superpowers-marketplace
claude plugin install double-shot-latte@superpowers-marketplace
claude plugin install superpowers-developing-for-claude-code@superpowers-marketplace
claude plugin install superpowers-lab@superpowers-marketplace
claude plugin install superpowers-chrome@superpowers-marketplace

# Official Anthropic
claude plugin install frontend-design@claude-plugins-official
claude plugin install feature-dev@claude-plugins-official
claude plugin install code-review@claude-plugins-official
claude plugin install commit-commands@claude-plugins-official
claude plugin install pr-review-toolkit@claude-plugins-official
claude plugin install hookify@claude-plugins-official
claude plugin install plugin-dev@claude-plugins-official
claude plugin install agent-sdk-dev@claude-plugins-official
claude plugin install security-guidance@claude-plugins-official
claude plugin install typescript-lsp@claude-plugins-official
claude plugin install gopls-lsp@claude-plugins-official
```

### Step 4: Update Existing (Optional)

**If `--update` flag passed:**

```bash
claude plugin update <plugin-name>
```

For each installed plugin.

### Step 5: Report Results

```
Plugin Setup Complete
=====================

Installed (new):
  - superpowers@superpowers-marketplace (v4.0.2)
  - episodic-memory@superpowers-marketplace (v1.0.15)
  ...

Already installed:
  - elements-of-style@superpowers-marketplace (v1.0.0)
  ...

Updated:
  - <plugin> (v1.0.0 -> v1.0.1)
  ...

Total: X plugins installed

NOTE: Restart Claude Code to activate new plugins.
```

---

## Key Plugins for This Workspace

| Repo Type | Recommended Plugins |
|-----------|---------------------|
| Next.js apps | `frontend-design`, `typescript-lsp` |
| Go projects | `gopls-lsp` |
| All repos | `superpowers`, `episodic-memory`, `double-shot-latte` |
| Plugin dev | `superpowers-developing-for-claude-code`, `plugin-dev` |

---

## Examples

```bash
/setup-plugins              # Install all missing plugins
/setup-plugins --check      # Only show what would be installed
/setup-plugins --update     # Install missing + update existing
```
