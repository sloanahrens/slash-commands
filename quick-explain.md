---
description: Quick code explanation using local model
---

# Quick Explain

Get a fast, concise explanation of code using a local MLX model.

**Arguments**: `$ARGUMENTS` - Code snippet, file path, or function name to explain

**See also**: `_local-model.md` for model configuration

---

## Process

### Step 1: Parse Input

Determine what to explain:
- If `$ARGUMENTS` is a file path → read the file
- If `$ARGUMENTS` looks like code → use directly
- If `$ARGUMENTS` is a function/class name → search codebase for it

### Step 2: Generate Explanation (Local Model)

```bash
mlx_lm.generate \
  --model mlx-community/DeepSeek-Coder-V2-Lite-Instruct-4bit-mlx \
  --max-tokens 150 \
  --prompt "Explain this code concisely (under 100 words):

<code>

Explanation:"
```

### Step 3: Display Result

```
[local] Explanation:
---
<explanation from local model>
---

(c) More detail with Claude  (done) Accept
```

### Step 4: Handle Follow-up

If user wants more detail `(c)`:
- Pass to Claude with full context
- Label output `[claude]`
- Provide deeper analysis

---

## Examples

```bash
/quick-explain "function retry<T>(fn, options)"
/quick-explain src/utils/async.ts
/quick-explain "What does validateEmail do?"
```

---

## Output Format

Always show model label:

```
[local] Explanation (0.4s):
---
The retry function attempts to execute an async function multiple times
with exponential backoff. It catches errors, waits with increasing delays,
and throws the last error if all attempts fail.
---
```

---

## When to Use

| Use `/quick-explain` | Use Claude directly |
|---------------------|---------------------|
| "What does this do?" | "Review this for bugs" |
| Single function | Multi-file analysis |
| Quick understanding | Deep architectural review |
| Repeated lookups | Security analysis |
