---
description: Run meta-agent improvement loop with parallel subagents
---

# Improve

Execute a meta-agent loop: spawn parallel subagents to explore, analyze findings, update notes, and iterate.

**Arguments**: `$ARGUMENTS` - `<repo> <task>` - Repository and task description

---

## Purpose

Automate the improvement cycle from the Confucius-inspired agent scaffolding. This command orchestrates multiple Task subagents in parallel, analyzes their results, and captures learnings as notes.

---

## Process

### Step 1: Prime Context

Run `/prime <repo>` to load relevant patterns and notes.

### Step 2: Parse Task

Extract the task from arguments:
- Identify task type: bug fix, feature, exploration, optimization
- Identify key terms for note searching

### Step 3: Spawn Parallel Subagents

Use the Task tool to launch 2-4 subagents in parallel based on task type:

**For bug fixes:**
```
1. Code Explorer: Search for relevant code, read surrounding context
2. History Analyzer: Check recent commits touching related files
3. Note Searcher: Search hindsight notes for similar issues
```

**For features:**
```
1. Pattern Finder: Find similar patterns in codebase
2. Dependency Mapper: Identify affected files and dependencies
3. Test Scanner: Find existing tests for the area
```

**For exploration:**
```
1. Structure Analyzer: Map directory structure and key files
2. Entry Point Finder: Identify main entry points and flows
3. Documentation Reader: Read README, CLAUDE.md, docs/
```

### Step 4: Launch Subagents

Use Task tool with appropriate subagent_type:

```python
# Launch in parallel (single message with multiple tool calls)
Task(
    subagent_type="Explore",
    description="Search codebase for <pattern>",
    prompt="Find all files related to <topic>. Report: file paths, key functions, dependencies."
)

Task(
    subagent_type="Bash",
    description="Check recent commits",
    prompt="Run: devbot log <repo> | look for commits mentioning <topic>"
)

Task(
    subagent_type="general-purpose",
    description="Search notes for <topic>",
    prompt="Search ~/.claude/notes/ for files mentioning <keywords>. Report relevant findings."
)
```

### Step 5: Collect Results

Wait for all subagents to complete. Gather their reports into a unified analysis.

### Step 6: Synthesize Findings

Analyze combined results:
- What did we learn?
- What patterns emerged?
- What failed or was unexpected?
- What's the recommended approach?

Present synthesis to user:

```
## Findings from 3 subagents

### Code Analysis
- Found 5 files related to <topic>
- Key entry point: src/handlers/auth.ts:45
- Related tests in __tests__/auth.test.ts

### Commit History
- Last modified 3 days ago (commit abc123)
- Previous fix attempt in commit def456 was reverted

### Note Search
- Found hindsight note: 2026-01-05-auth-timeout.md
  "Auth tokens expire silently - add explicit error handling"

## Recommended Approach
Based on findings, suggest:
1. Add token expiry check in auth.ts:89
2. Follow pattern from payment-service.ts:120
3. Add test case for expiry scenario
```

### Step 7: Attempt Task (if actionable)

If findings suggest a clear fix:

1. Ask user: "Proceed with fix? [Y/n]"
2. If yes, implement the fix
3. Run relevant tests via `devbot check <repo>`

### Step 8: Capture Learnings

Based on outcome:

**If fix succeeded:**
```
Run /capture-hindsight to save:
- What the issue was
- What fixed it
- Pattern for future reference
```

**If fix failed:**
```
Run /capture-hindsight to save:
- What was attempted
- Why it failed
- What to try differently
```

**If exploration only:**
```
Create session note with:
- What was learned
- Key files and patterns discovered
- Open questions
```

### Step 9: Iterate (if needed)

If task not complete:
1. Update notes with partial progress
2. Ask user if they want another iteration
3. If yes, return to Step 3 with refined focus

---

## Subagent Templates

Pre-configured prompts for common exploration patterns:

### explore-codebase
```
Explore <repo> to understand <area>.
Report:
1. Key files and their purposes
2. Data flow through the area
3. External dependencies
4. Test coverage status
```

### find-similar-patterns
```
Search <repo> for patterns similar to <description>.
Report:
1. File locations with similar code
2. How the pattern is implemented elsewhere
3. Commonalities and differences
```

### trace-error
```
Investigate <error> in <repo>.
Report:
1. Where the error originates
2. Call stack / control flow
3. Related error handling
4. Recent changes to affected code
```

### check-test-coverage
```
Analyze test coverage for <area> in <repo>.
Report:
1. Existing test files
2. What's covered vs missing
3. Test patterns used
4. Suggested test cases
```

---

## Options

| Flag | Effect |
|------|--------|
| `--dry-run` | Show what subagents would be spawned without running |
| `--no-capture` | Skip automatic hindsight capture |
| `--iterations=N` | Max improvement iterations (default: 3) |
| `--template=<name>` | Use specific subagent template |

---

## Output Format

```
/improve mango "fix flaky test in payment service"
================================================

## Phase 1: Priming
Loading context for: mango
- 2 patterns loaded
- 1 hindsight note found

## Phase 2: Spawning Subagents
Launching 3 parallel agents...
├── [explore] Searching payment service code
├── [bash] Checking recent commits
└── [general] Searching notes for "flaky" "payment"

## Phase 3: Results (12 seconds)
[explore] Found 3 relevant files, key function at payment.ts:234
[bash] Commit abc123 added retry logic 5 days ago
[general] Hindsight note mentions "race condition in mock"

## Phase 4: Synthesis
The flaky test is likely due to the race condition mentioned in
hindsight note 2026-01-06-payment-mock-race.md. The recent commit
attempted a fix but may have introduced timing issues.

## Recommended Fix
1. Add explicit wait in test setup (payment.test.ts:45)
2. Use deterministic mock from test-utils.ts

Proceed with fix? [Y/n]
```

---

## Examples

```bash
/improve mango fix flaky test in payment service
/improve atap-automation2 investigate timeout issues
/improve slash-commands add new devbot command
/improve --template=trace-error mango "TypeError in checkout"
/improve --dry-run fractals-nextjs optimize render performance
```

---

## Related Commands

- `/prime` — Load context before improvement loop
- `/capture-hindsight` — Save learnings from the session
- `/find-tasks` — Discover tasks to improve
