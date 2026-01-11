# Confucius-Inspired Agent Scaffolding for mono-claude

**Date:** 2026-01-11
**Status:** Draft
**Reference:** [Confucius Code Agent Paper (arXiv:2512.10398)](https://arxiv.org/pdf/2512.10398)

---

## Executive Summary

This design describes how to incrementally add CCA-style capabilities to improve Claude Code's effectiveness in the mono-claude workspace. Rather than building a standalone orchestrator, we enhance Claude Code sessions with structured memory, failure capture, and pattern promotion using existing tools and filesystem conventions.

**Core insight from CCA:** The agent's "working memory" should be structured and distilled, not raw conversation history. Capture decisions, failures, and patterns in persistent notes that future sessions can reference.

---

## Goals

- Reduce repeated mistakes across sessions
- Shorter sessions for recurring tasks
- Visible learning: patterns extracted and preserved
- Leverage existing tools (devbot, slash commands, episodic memory)

## Non-Goals (for now)

- Building a standalone orchestrator
- Replacing Claude Code's execution model
- Formal extension framework with typed callbacks
- SWE-Bench evaluation or automated benchmarking

---

## Architecture

```
┌─────────────────────────────────────────────┐
│          Claude Code Session                │  ← Orchestrator (user + Claude)
│  • Spawns Task subagents for parallel work  │
│  • Analyzes findings, makes decisions       │
│  • Manages improvement loops                │
├─────────────────────────────────────────────┤
│  devbot (Go)         │  mlx-hub (TS)        │  ← Extensions
│  • Fast repo ops     │  • Local inference   │
│  • Workspace cmds    │  • Quick gen tasks   │
├─────────────────────────────────────────────┤
│  Patterns (versioned)│  Notes (local)       │  ← Memory
│  docs/patterns/      │  ~/.claude/notes/    │
├─────────────────────────────────────────────┤
│         episodic-memory plugin              │  ← Cross-session search
│         hookify rules                       │  ← Behavior guardrails
└─────────────────────────────────────────────┘
```

---

## Phased Approach

| Phase | Mechanism | What we build |
|-------|-----------|---------------|
| 0 | Foundation | Note conventions, directory structure |
| 1 | Hindsight Notes | Capture failure modes after errors |
| 2 | Session Notes | Summarize long sessions into reusable context |
| 3 | Context Priming | Auto-search notes at session start |
| 4 | Meta-Agent Loop | Spawn → analyze → improve cycle |

---

## Note Structure (Split Approach)

Two locations with different purposes:

### Versioned Patterns — `slash-commands/docs/patterns/`

```
slash-commands/docs/patterns/
├── README.md                    # Pattern conventions
├── bash-execution.md            # devbot exec vs cd &&
├── atap-field-selectors.md      # Label-based selectors
├── monorepo-navigation.md       # Subdir patterns
└── hookify-rules.md             # What's blocked and why
```

**Characteristics:**
- Proven knowledge (validated across 2+ sessions)
- Written in timeless style (not "today I learned...")
- Committed to git, travels with slash-commands
- Referenced by CLAUDE.md or slash commands

### Local Notes — `~/.claude/notes/`

```
~/.claude/notes/
├── hindsight/
│   ├── 2026-01-11-cd-compound-blocked.md
│   └── 2026-01-11-atap-timeout.md
├── sessions/
│   ├── 2026-01-11-devbot-prereq.md
│   └── 2026-01-11-zapier-integration.md
└── index.md                     # Optional quick reference
```

**Characteristics:**
- Temporal, dated, potentially verbose
- Not committed anywhere
- Searchable by `/prime` command
- Candidates for promotion to patterns

### Frontmatter Format

```yaml
---
type: hindsight | session | pattern
repos: [atap-automation2, mango]  # Affected repos, or [all]
tags: [hookify, bash, git]        # Searchable tags
created: 2026-01-11
updated: 2026-01-11               # Patterns only
status: active | promoted | stale # Local notes only
---
```

---

## Phase 1: Hindsight Notes

**Purpose:** Capture failure modes so future sessions don't repeat them.

**When to create:**
- Claude made an error that required backtracking
- A command was blocked by hookify
- Multiple attempts needed to solve something
- A pattern was discovered that wasn't obvious

**Template:**

```markdown
---
type: hindsight
repos: [slash-commands]
tags: [bash, hookify, devbot]
created: 2026-01-11
status: active
---

# Hookify blocks cd && compound commands

## What happened
Attempted `cd /path/to/repo && npm test` which hookify blocked.

## Why it failed
Hookify rule `block-dangerous` prevents compound bash commands to avoid
unintended side effects from chained execution.

## Correct approach
Use `devbot exec <repo> <command>` instead:
```bash
devbot exec atap-automation2 npm test
```

## Applies when
- Running commands in repo directories
- Any situation where you'd normally `cd` first
```

---

## Phase 2: Session Notes

**Purpose:** Distill long sessions into reusable context.

**When to create:**
- Completed a multi-step implementation
- Explored a codebase area thoroughly
- Made architectural decisions worth preserving
- Session exceeded ~30 minutes of substantive work

**Template:**

```markdown
---
type: session
repos: [devbot, slash-commands]
tags: [prereq, implementation, go]
created: 2026-01-11
status: active
---

# Implemented devbot prereq command

## Task
Add prerequisite validation before starting work on repos.

## Key decisions
1. Reuse `detect.ProjectStack()` for stack detection
2. Check tools, deps, and env vars (from .env.local.example)
3. Exit code 0 = pass, 1 = any failure

## What was built
- `internal/prereq/prereq.go` — Main orchestration
- `internal/prereq/tools.go` — Binary existence checks
- `internal/prereq/env.go` — Env file comparison

## Patterns discovered
- Parse .env.example for var names with regex `^[A-Z][A-Z0-9_]*=`
- Tool versions via `<tool> --version`, first line only

## Open items
- Could add `subprojects` config for monorepo-wide checks
```

---

## Phase 3: Context Priming

**Purpose:** Surface relevant notes at session start.

**Flow:**

```
Session Start
     │
     ▼
┌─────────────────────────────────┐
│ 1. User states task/repo        │
└─────────────────────────────────┘
     │
     ▼
┌─────────────────────────────────┐
│ 2. Search notes by:             │
│    • repo name in frontmatter   │
│    • tags matching task keywords│
│    • recent hindsight notes     │
└─────────────────────────────────┘
     │
     ▼
┌─────────────────────────────────┐
│ 3. Load 2-3 most relevant notes │
│    into working context         │
└─────────────────────────────────┘
     │
     ▼
┌─────────────────────────────────┐
│ 4. Proceed with task            │
└─────────────────────────────────┘
```

**Implementation:** `/prime <repo>` slash command that searches both locations.

---

## Phase 4: Meta-Agent Loop

**Purpose:** Automate the improvement cycle.

**Loop structure:**

```
┌──────────────────────────────────────────────────┐
│                 Claude Code Session              │
│                  (Orchestrator)                  │
└──────────────────────────────────────────────────┘
        │                              ▲
        │ 1. Spawn parallel            │ 4. Analyze
        │    Task subagents            │    results
        ▼                              │
┌──────────────────────────────────────────────────┐
│              Task Subagents                      │
│  • Explore codebase areas                        │
│  • Run tests/checks                              │
│  • Attempt implementations                       │
└──────────────────────────────────────────────────┘
        │
        │ 2. Execute & report
        ▼
┌──────────────────────────────────────────────────┐
│              Results & Failures                  │
│  • What worked                                   │
│  • What failed (and why)                         │
│  • Patterns observed                             │
└──────────────────────────────────────────────────┘
        │
        │ 3. Distill into notes
        ▼
┌──────────────────────────────────────────────────┐
│              Notes Update                        │
│  • New hindsight note if failure                 │
│  • Update session note with findings             │
│  • Promote pattern if reusable                   │
└──────────────────────────────────────────────────┘
```

**Example workflow:**

```
User: "Fix the flaky test in atap-automation2"

Claude (orchestrator):
1. /prime atap-automation2 → loads relevant notes
2. Spawn 3 parallel Task agents:
   - Agent A: grep for the test, read surrounding code
   - Agent B: check recent commits touching test files
   - Agent C: search hindsight notes for "flaky" or "timeout"
3. Analyze findings, form hypothesis
4. Attempt fix
5. If fix fails → write hindsight note
6. If fix succeeds → update session note with pattern
```

---

## Hello World Deliverables

| Item | Location | Description |
|------|----------|-------------|
| `~/.claude/notes/` | Local | Directory structure (hindsight/, sessions/) |
| `docs/patterns/` | slash-commands | Directory + README with conventions |
| `docs/patterns/bash-execution.md` | slash-commands | Seed pattern: devbot exec usage |
| `docs/patterns/hookify-rules.md` | slash-commands | Seed pattern: what's blocked |
| `/prime` | slash-commands | Search both locations, display relevant notes |
| `/capture-hindsight` | slash-commands | Create note in ~/.claude/notes/hindsight/ |
| `/promote-pattern` | slash-commands | Move local note → docs/patterns/ |

---

## Slash Command Specifications

### `/prime <repo>`

```markdown
1. Search docs/patterns/*.md for repo name or "all" in frontmatter
2. Search ~/.claude/notes/**/*.md for repo name in frontmatter
3. Display matches grouped by type:

   ## Patterns (versioned)
   - bash-execution.md — devbot exec vs cd &&

   ## Recent Hindsight (local)
   - 2026-01-11-atap-timeout.md — ATAP session recovery

   ## Recent Sessions (local)
   - 2026-01-11-zapier-integration.md
```

### `/capture-hindsight`

```markdown
1. Ask: What repo(s)? What tags?
2. Ask: What happened? (paste error or describe)
3. Ask: What was the fix?
4. Generate note with template, write to ~/.claude/notes/hindsight/
5. Confirm: "Hindsight captured: ~/.claude/notes/hindsight/2026-01-11-<slug>.md"
```

### `/promote-pattern`

```markdown
1. List recent hindsight notes
2. User selects one
3. Claude generalizes content (removes temporal language)
4. Write to docs/patterns/<slug>.md
5. Update original: status: promoted
6. Offer to commit the new pattern
```

---

## Promotion Workflow

```
Hindsight note created
        │
        ▼
Referenced in 2+ sessions?
        │
    yes │
        ▼
Generalize & clean up
        │
        ▼
Write to docs/patterns/
        │
        ▼
Mark local note: status: promoted
```

---

## Implementation Sequence

### Week 1: Foundation

1. Create `~/.claude/notes/{hindsight,sessions}/` structure
2. Create `docs/patterns/` with README.md
3. Write 2-3 seed patterns from existing knowledge
4. Test: can grep find them?

### Week 2: Slash Commands

5. Write `/prime` command
6. Write `/capture-hindsight` command
7. Write `/promote-pattern` command
8. Update `_shared-repo-logic.md` to mention note search

### Week 3: Validation

9. Use in real sessions for a week
10. Capture 5+ hindsight notes naturally
11. Promote at least 1 to pattern
12. Retrospective: what's working, what's friction?

---

## Future Enhancements

### Near-term (after hello world validated)

| Enhancement | Description | Trigger |
|-------------|-------------|---------|
| SessionStart hook | Auto-run `/prime` on session start | Manual priming feels repetitive |
| Tag-based search | `/prime --tag=hookify` across all notes | Repo-based search too narrow |
| Note aging | Mark notes `stale` after 30 days | Notes accumulating |
| devbot notes | `devbot notes search <query>` | Grep too slow |

### Medium-term (Meta-Agent)

| Enhancement | Description | Trigger |
|-------------|-------------|---------|
| `/improve` command | Full meta-agent loop | Want automated exploration |
| Subagent templates | Pre-configured Task prompts | Repeatedly writing similar prompts |
| Failure detection hook | Auto-trigger `/capture-hindsight` | Missing failure captures |
| Pattern suggestions | "Referenced 3x, promote?" | Manual promotion forgotten |

### Longer-term (CCA parity)

| Enhancement | Description | Trigger |
|-------------|-------------|---------|
| Context compression | Summarize mid-session | Long sessions degrading |
| Hierarchical memory | Session/task/subtask scopes | Flat notes insufficient |
| Note quality scoring | Track which notes helped | Can't tell useful from noise |
| Learning metrics | Turns reduced, tokens saved | Want quantified improvement |

### Not Planned

| Item | Reason |
|------|--------|
| Standalone orchestrator | Claude Code is the orchestrator |
| Formal extension framework | Slash commands + devbot sufficient |
| SWE-Bench evaluation | Focus on workspace effectiveness |
| RL-based training | Out of scope |

---

## Success Criteria

- [ ] Can prime a session with relevant notes in <5 seconds
- [ ] Hindsight notes capture failures within same session
- [ ] At least 3 patterns promoted within first month
- [ ] Subjective: fewer repeated mistakes on familiar repos
- [ ] Measurable (future): reduced turns/tokens for recurring tasks

---

## References

- [Confucius Code Agent Paper Summary](../../whitepaper-summaries/confucius-code-agent-2512.10398.md)
- [devbot README](../devbot/README.md)
- [Shared Repo Logic](../_shared-repo-logic.md)
