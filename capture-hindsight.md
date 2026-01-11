---
description: Capture a failure or lesson learned as a hindsight note
---

# Capture Hindsight

Create a hindsight note after encountering an issue, error, or discovering a non-obvious pattern.

**Arguments**: `$ARGUMENTS` - Optional: brief description of what happened (e.g., "hookify blocked my cd command")

---

## Purpose

Capture failure modes so future sessions don't repeat them. This implements the "hindsight notes" phase of the Confucius-inspired agent scaffolding.

CCA research showed hindsight notes had clear measurable impact: reduced turns (64→61), reduced tokens (104k→93k), and improved task completion.

---

## Process

### Step 1: Gather Context

Ask the user (using AskUserQuestion tool) for missing information:

**Question 1: Which repo(s) does this apply to?**
- Options from config.yaml repos
- "all" for workspace-wide patterns
- Allow multiple selection

**Question 2: What tags describe this?**
- Suggest based on recent context: bash, hookify, git, testing, devbot, etc.
- Allow custom tags

### Step 2: Capture the Failure

If `$ARGUMENTS` provided, use that as starting point. Otherwise ask:

**Question 3: What happened?**
- "Paste the error message or describe what went wrong"
- Look at recent conversation for context

**Question 4: What was the fix?**
- "What command/approach actually worked?"
- Pull from recent successful commands if available

### Step 3: Generate Note

Create the hindsight note with this structure:

```markdown
---
type: hindsight
repos: [<selected-repos>]
tags: [<selected-tags>]
created: <today's date>
status: active
---

# <Brief title derived from description>

## What happened
<Error message or description of the failure>

## Why it failed
<Root cause analysis - not just symptoms>

## Correct approach
<The solution that worked, with code examples>

## Applies when
<Pattern matching criteria for future sessions>
```

### Step 4: Generate Filename

Create a slug from the title:
- Lowercase
- Replace spaces with hyphens
- Remove special characters
- Prefix with date

Example: `2026-01-11-hookify-blocked-cd-command.md`

### Step 5: Write the Note

Write to `~/.claude/notes/hindsight/<filename>`:

```bash
# Ensure directory exists
mkdir -p ~/.claude/notes/hindsight
```

Then use Write tool to create the file.

### Step 6: Confirm

```
✓ Hindsight captured: ~/.claude/notes/hindsight/2026-01-11-hookify-blocked-cd-command.md

  Repos: slash-commands, all
  Tags: hookify, bash, devbot

  This note will appear when you run /prime for these repos.
  If this pattern proves useful across sessions, run /promote-pattern to make it permanent.
```

---

## Quick Capture Mode

If context is clear from the conversation, skip questions and generate directly:

```bash
/capture-hindsight hookify blocked cd && npm test
```

Claude should:
1. Infer repos from current context
2. Infer tags from the description
3. Pull the error and fix from recent conversation
4. Generate and write the note
5. Show confirmation with option to edit

---

## Output Format

```
Capturing hindsight...

Repos: [atap-automation2]
Tags: [hookify, bash]

---
# Hookify blocks cd && compound commands

## What happened
Attempted `cd /path/to/repo && npm test` which hookify blocked with message:
"Compound commands with && are blocked. Use devbot exec instead."

## Why it failed
Hookify rule `block-dangerous` prevents compound bash commands to avoid
unintended side effects from partial execution on failure.

## Correct approach
Use `devbot exec <repo> <command>`:
```bash
devbot exec atap-automation2 npm test
```

## Applies when
- Running any command that requires being in a specific directory
- Any situation where you'd normally `cd` first
---

Save to ~/.claude/notes/hindsight/2026-01-11-hookify-cd-blocked.md? [Y/n]
```

---

## From Conversation Context

Claude should proactively suggest `/capture-hindsight` when:
- An error was encountered and resolved
- Multiple attempts were needed to find the right approach
- A hookify rule blocked a command
- A non-obvious pattern was discovered

Prompt: "That was tricky. Want me to capture this as a hindsight note for future sessions?"

---

## Examples

```bash
/capture-hindsight                           # Interactive mode
/capture-hindsight timeout in atap form      # Quick capture with description
/capture-hindsight devbot exec saved the day # Quick capture
```

---

## Related Commands

- `/prime <repo>` — Load notes before starting work
- `/promote-pattern` — Promote useful hindsight to versioned pattern
