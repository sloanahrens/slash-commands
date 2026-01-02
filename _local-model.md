# Local Model Helper

Shared logic for invoking a local MLX model in slash commands.

## Configuration

Configure your local model in this section. Default: DeepSeek-Coder.

| Setting | Value |
|---------|-------|
| **Model** | `mlx-community/DeepSeek-Coder-V2-Lite-Instruct-4bit-mlx` |
| **Speed** | ~143 tok/s |
| **Memory** | ~9 GB |
| **Best for** | Commit messages, code explanations, simple code gen |

## Prerequisites

Install mlx-lm if not already installed:
```bash
pip3 install mlx-lm
```

Find the mlx_lm.generate path:
```bash
which mlx_lm.generate || python3 -c "import mlx_lm; print(mlx_lm.__file__)"
```

## Invoking the Local Model

```bash
mlx_lm.generate \
  --model mlx-community/DeepSeek-Coder-V2-Lite-Instruct-4bit-mlx \
  --max-tokens <limit> \
  --prompt "<prompt>"
```

If `mlx_lm.generate` is not on PATH, use the full path (e.g., `~/.local/bin/mlx_lm.generate` or similar).

## Output Format

When displaying local model output, always label it clearly:

```
[local] Commit message:
---
<output from model>
---

Options:
- (y) Accept - Use this output
- (c) Claude - Regenerate with Claude
- (e) Edit - Modify manually
```

## Error Handling

If the local model fails (not installed, memory issue, etc.):
1. Log: `[local] Failed: <reason>`
2. Fall back to Claude automatically
3. Label Claude output: `[claude] Commit message:`

## Task Guidelines

| Task | Max Tokens | Notes |
|------|------------|-------|
| Commit message | 100 | Short, focused |
| Code explanation | 150 | Concise summary |
| Simple function | 200 | Single function gen |

## When to Use Claude Instead

- Code review (needs high accuracy)
- Multi-file analysis
- Complex refactoring
- Security-sensitive code
- Architectural decisions

## Customizing the Model

To use a different MLX model:

1. Download via Hugging Face:
   ```bash
   mlx_lm.manage --model <model-name> download
   ```

2. Update the model name in commands that reference `_local-model.md`

3. Adjust max-tokens based on model capabilities

## Supported Models

Any MLX-compatible model from `mlx-community` on Hugging Face works. Popular choices:
- `DeepSeek-Coder-V2-Lite-Instruct-4bit-mlx` (coding, recommended)
- `Llama-3.2-3B-Instruct-4bit` (general, faster)
- `Mistral-7B-Instruct-v0.3-4bit` (general, balanced)
