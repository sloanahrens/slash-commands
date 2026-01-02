---
description: Quick code generation using local model
---

# Quick Gen

Generate simple code snippets using a local MLX model.

**Arguments**: `$ARGUMENTS` - Description of what to generate

**See also**: `_local-model.md` for model configuration

---

## Process

### Step 1: Parse Request

Extract from `$ARGUMENTS`:
- What to generate (function, type, test, etc.)
- Language (infer from context or current directory)
- Any constraints mentioned

### Step 2: Gather Context

If in a repo, check for style patterns:
```bash
# Get language from file extensions
ls <repo-path>/src/*.{ts,go,py} 2>/dev/null | head -1

# Sample existing code style
head -30 <similar-file>
```

### Step 3: Generate Code (Local Model)

```bash
mlx_lm.generate \
  --model mlx-community/DeepSeek-Coder-V2-Lite-Instruct-4bit-mlx \
  --max-tokens 200 \
  --prompt "Generate <language> code for: $ARGUMENTS

Follow this style:
<sample from repo>

Code:"
```

### Step 4: Display Result

```
[local] Generated code:
---
<code from local model>
---

(y) Use this  (c) Regenerate with Claude  (e) Edit
```

### Step 5: Apply If Accepted

If user accepts:
- Insert into appropriate file
- Or display for manual copy

---

## Examples

```bash
/quick-gen "validatePhone function for US numbers"
/quick-gen "test for the retry function"
/quick-gen "TypeScript interface for User with name, email, role"
```

---

## Output Format

Always show model label:

```
[local] Generated (0.6s):
---
export function validatePhone(phone: string): boolean {
  const pattern = /^\+?1?[-.\s]?\(?[0-9]{3}\)?[-.\s]?[0-9]{3}[-.\s]?[0-9]{4}$/;
  return pattern.test(phone);
}
---

(y) Use  (c) Claude  (e) Edit
```

---

## When to Use

| Use `/quick-gen` | Use Claude directly |
|-----------------|---------------------|
| Single function | Multi-file feature |
| Simple utility | Complex logic |
| Boilerplate | Architecture decisions |
| Quick prototype | Production code |

---

## Limitations

Local model works best for:
- Functions under 50 lines
- Well-defined inputs/outputs
- Common patterns

For complex generation, choose `(c)` for Claude.
