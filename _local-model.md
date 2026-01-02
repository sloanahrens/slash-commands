# Local Model Helper

Shared logic for invoking the local MLX model (Qwen2.5-Coder) in slash commands.

## Configuration

**Model**: `mlx-community/Qwen2.5-Coder-14B-Instruct-4bit`
**Size**: 7.7 GB
**Speed**: ~15-50 tok/s (cold/warm)
**Best for**: Commit messages, code explanations, type fixes, bug detection

## Invoking the Local Model

Use the mlx-hub MCP tool:
```
mcp__plugin_mlx-hub_mlx-hub__mlx_infer
  model_id: mlx-community/Qwen2.5-Coder-14B-Instruct-4bit
  prompt: "<prompt>"
  max_tokens: <limit>
```

## Output Format

When displaying local model output, always label it clearly:

```
[qwen] Commit message:
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
1. Log: `[qwen] Failed: <reason>`
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

## Supported Models

Models available via mlx-hub plugin. Check with:
```
mcp__plugin_mlx-hub_mlx-hub__mlx_list_local
```

Common choices:
- `mlx-community/Qwen2.5-Coder-14B-Instruct-4bit` (coding, recommended)
- `mlx-community/Llama-3.2-3B-Instruct-4bit` (general, faster)
