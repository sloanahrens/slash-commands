---
name: explore-codebase
description: Understand an area of code - structure, dependencies, patterns
subagent_type: Explore
---

# Explore Codebase Template

Explore {repo} to understand {area}.

## Your Task

Investigate the codebase area and report back with structured findings.

## Steps

1. Use Glob to find relevant files matching the area
2. Read key files to understand structure
3. Identify dependencies and imports
4. Note any patterns or conventions used

## Report Format

Provide your findings in this structure:

### Key Files
List the most important files with one-line descriptions:
- `path/to/file.ts` — Purpose of this file

### Data Flow
Describe how data moves through this area:
1. Entry point
2. Processing steps
3. Output/side effects

### Dependencies
- **Internal:** Other modules this area depends on
- **External:** Third-party libraries used

### Patterns Observed
Note any coding patterns, conventions, or abstractions used.

### Test Coverage
- Test files found: (list)
- Coverage gaps: (observed)

### Questions/Concerns
Any unclear areas or potential issues noted.
