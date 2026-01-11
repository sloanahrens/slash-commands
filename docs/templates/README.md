# Subagent Templates

Pre-configured prompts for the `/improve` command's Task subagents.

## Usage

Templates are referenced by name in `/improve`:

```bash
/improve --template=trace-error mango "TypeError in checkout"
```

Or used automatically based on task type detection.

## Template Format

Each template defines:
- **name** — Template identifier
- **description** — When to use this template
- **subagent_type** — Which Task agent type to use
- **prompt** — The prompt template with `{placeholders}`

## Available Templates

| Template | Purpose | Agent Type |
|----------|---------|------------|
| explore-codebase | Understand an area of code | Explore |
| find-similar-patterns | Find existing patterns | Explore |
| trace-error | Debug an error | general-purpose |
| check-test-coverage | Analyze test coverage | Explore |
| search-history | Check git history | Bash |
| search-notes | Search hindsight/session notes | general-purpose |

## Creating New Templates

1. Create `<name>.md` in this directory
2. Add YAML frontmatter with template metadata
3. Write the prompt template with placeholders
4. Document expected outputs

## Placeholders

| Placeholder | Replaced With |
|-------------|---------------|
| `{repo}` | Repository name |
| `{area}` | Area/module being investigated |
| `{error}` | Error message or description |
| `{keywords}` | Search keywords from task |
| `{files}` | Relevant file paths |
