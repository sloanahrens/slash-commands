---
name: trace-error
description: Debug an error by tracing its origin and context
subagent_type: general-purpose
---

# Trace Error Template

Investigate "{error}" in {repo}.

## Your Task

Trace the error to its source and gather context for fixing it.

## Steps

1. Search for the error message in the codebase
2. Identify where the error is thrown/generated
3. Trace the call stack / control flow
4. Check recent commits affecting the area
5. Look for related error handling

## Report Format

### Error Location
- **File:** path/to/file.ts
- **Line:** 123
- **Function:** functionName

### Call Flow
How execution reaches the error:
1. Entry point →
2. Intermediate calls →
3. Error site

### Error Context
```
<relevant code snippet around error>
```

### Related Error Handling
- Existing try/catch blocks
- Error boundaries
- Fallback mechanisms

### Recent Changes
Commits in last 2 weeks affecting this area:
- `abc123` — Description
- `def456` — Description

### Root Cause Hypothesis
Based on investigation, the error likely occurs because...

### Suggested Fix
Recommend approach to resolve:
1. Step one
2. Step two

### Related Notes
Any hindsight notes mentioning similar issues.
