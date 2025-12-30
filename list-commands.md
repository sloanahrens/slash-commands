---
description: List all available slash commands
---

# List Commands

Display all available slash commands in this workspace.

---

## Commands

| Command | Description |
|---------|-------------|
| `/super <repo>` | Start brainstorming session with full context |
| `/find-tasks <repo>` | Suggest 3-5 high-priority tasks |
| `/run-tests <repo>` | Run lint, type-check, build, and tests |
| `/yes-commit <repo>` | Draft and commit changes |
| `/push <repo>` | Push commits to origin |
| `/update-docs <repo>` | Update CLAUDE.md, README, docs |
| `/review-project <repo>` | Technical review to docs/tech-review.md |
| `/add-repo <url>` | Clone repo and add to config |
| `/status [repo]` | Show status overview of all or one repo |
| `/sync [repo]` | Pull latest changes for all or one repo |
| `/switch <repo>` | Quick context switch to a repo |
| `/dev-rules` | Remind Claude of workspace rules |
| `/setup-plugins` | Install all recommended plugins |
| `/make-test <repo>` | Test Makefile targets interactively |
| `/list-commands` | List all available commands (this) |
| `/list-skills` | List available skills from plugins |
| `/linear <subcommand>` | Interact with Linear issues |

All repo commands support fuzzy matching via aliases (e.g., `/run-tests app`).

---

## Usage Tips

- **No argument**: Most commands show an interactive selection menu
- **Fuzzy match**: Use partial names or aliases (`app` → `my-nextjs-app`)
- **Flags**: Some commands accept flags (`/run-tests --fix`, `/make-test --dry-run`)

---

## Getting Help

- `/list-skills` - See available skills from installed plugins
- Read `_shared-repo-logic.md` for how repo selection works
- Check `config.yaml` for configured repositories
