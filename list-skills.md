---
description: List available skills from installed plugins
---

# List Skills

Display skills available from plugins installed via `/setup-plugins`.

---

## Core Skills (superpowers)

These are the most commonly used skills for development workflows.

| Skill | When to Use |
|-------|-------------|
| `superpowers:brainstorming` | Before creative work, designing features, exploring requirements |
| `superpowers:writing-plans` | Creating detailed implementation plans from specs |
| `superpowers:systematic-debugging` | Bug investigation - gather evidence before hypothesizing |
| `superpowers:test-driven-development` | Writing new code with test-first approach |
| `superpowers:verification-before-completion` | Before claiming work is done - run tests, verify output |
| `superpowers:dispatching-parallel-agents` | When facing 2+ independent tasks |
| `superpowers:receiving-code-review` | When processing code review feedback |
| `superpowers:using-git-worktrees` | Starting feature work that needs isolation |
| `superpowers:subagent-driven-development` | Executing plans with independent tasks |
| `superpowers:execute-plan` | Execute plan in batches with review checkpoints |

---

## Writing & Documentation (elements-of-style)

| Skill | When to Use |
|-------|-------------|
| `elements-of-style:writing-clearly-and-concisely` | Any prose humans will read - docs, commits, errors, UI text |

---

## Memory (episodic-memory)

| Skill | When to Use |
|-------|-------------|
| `episodic-memory:remembering-conversations` | When user asks "how should I..." or you're stuck on something discussed before |
| `episodic-memory:search-conversations` | Search previous Claude Code conversations |

---

## Feature Development (feature-dev)

| Skill | When to Use |
|-------|-------------|
| `feature-dev:feature-dev` | Guided feature development with codebase understanding |

Agents available:
- `feature-dev:code-reviewer` - Reviews code for bugs, security, quality
- `feature-dev:code-explorer` - Analyzes codebase features and architecture
- `feature-dev:code-architect` - Designs feature architectures

---

## Code Review (pr-review-toolkit)

| Skill | When to Use |
|-------|-------------|
| `pr-review-toolkit:review-pr` | Comprehensive PR review using specialized agents |

Agents available:
- `pr-review-toolkit:code-reviewer` - Adherence to guidelines and best practices
- `pr-review-toolkit:silent-failure-hunter` - Find inadequate error handling
- `pr-review-toolkit:code-simplifier` - Simplify code for clarity
- `pr-review-toolkit:comment-analyzer` - Analyze comments for accuracy
- `pr-review-toolkit:pr-test-analyzer` - Review test coverage
- `pr-review-toolkit:type-design-analyzer` - Analyze type design quality

---

## Git & Commits (commit-commands)

| Skill | When to Use |
|-------|-------------|
| `commit-commands:commit` | Create a git commit |
| `commit-commands:commit-push-pr` | Commit, push, and open a PR |
| `commit-commands:clean_gone` | Clean up merged/deleted branches |

---

## Frontend Design (frontend-design)

| Skill | When to Use |
|-------|-------------|
| `frontend-design:frontend-design` | Building web components, pages, or applications with high design quality |

---

## Hooks (hookify)

| Skill | When to Use |
|-------|-------------|
| `hookify:hookify` | Create hooks to prevent unwanted behaviors |
| `hookify:writing-rules` | Create hookify rules |
| `hookify:list` | List configured hookify rules |
| `hookify:configure` | Enable/disable hookify rules |

---

## Plugin Development (plugin-dev)

| Skill | When to Use |
|-------|-------------|
| `plugin-dev:create-plugin` | End-to-end plugin creation workflow |
| `plugin-dev:command-development` | Create slash commands |

Agents available:
- `plugin-dev:agent-creator` - Create autonomous agents for plugins
- `plugin-dev:skill-reviewer` - Review skill quality
- `plugin-dev:plugin-validator` - Validate plugin structure

---

## Agent SDK (agent-sdk-dev)

| Skill | When to Use |
|-------|-------------|
| `agent-sdk-dev:new-sdk-app` | Create a new Claude Agent SDK application |

Agents available:
- `agent-sdk-dev:agent-sdk-verifier-ts` - Verify TypeScript Agent SDK apps
- `agent-sdk-dev:agent-sdk-verifier-py` - Verify Python Agent SDK apps

---

## Experimental (superpowers-lab)

| Skill | When to Use |
|-------|-------------|
| `superpowers-lab:mcp-cli` | Use MCP servers on-demand via CLI |
| `superpowers-lab:using-tmux-for-interactive-commands` | Run interactive CLI tools (vim, git rebase -i, REPL) |

---

## Browser Control (superpowers-chrome)

| Skill | When to Use |
|-------|-------------|
| `superpowers-chrome:browsing` | Direct browser control via Chrome DevTools Protocol |

---

## Claude Code Development (superpowers-developing-for-claude-code)

| Skill | When to Use |
|-------|-------------|
| `superpowers-developing-for-claude-code:working-with-claude-code` | Working with CLI, plugins, hooks, MCP, config |
| `superpowers-developing-for-claude-code:developing-claude-code-plugins` | Creating, modifying, testing plugins |

---

## How to Use Skills

Skills are invoked automatically when relevant, or you can invoke them directly:

```
Use the superpowers:brainstorming skill to design this feature
```

Or via the Skill tool in code.

---

## Installing Missing Plugins

If any skills are missing, run:
```
/setup-plugins
```
