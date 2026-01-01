---
description: Start brainstorming session with workspace context
---

# Super Command

Start a structured brainstorming session with full context about the workspace and selected repository. Leverages local MLX models for acceleration where appropriate.

**Arguments**: `$ARGUMENTS` - Optional repo name or task description. If repo recognized, selects it. Otherwise treated as brainstorm topic.

**Shared logic**: See `_shared-repo-logic.md` for repo discovery, selection, and commit rules.

---

## Process

### Step 0: Check Plugins

Before running the brainstorming session, verify required and optional plugins are available:

```bash
claude plugin list 2>/dev/null | grep -E "(superpowers|mlx)"
```

#### 0.1: Superpowers Plugin (Required)

The superpowers plugin provides the `/superpowers:brainstorming` skill used in Step 4.

**If superpowers is NOT installed:**

1. Inform the user: "The superpowers plugin is required for structured brainstorming but isn't currently installed."
2. Offer to install it:
   ```bash
   claude plugin add superpowers@superpowers-marketplace
   ```
3. If the user declines or installation fails, **stop here** - the `/super` command cannot complete without the brainstorming skill. Suggest alternatives:
   - Install the plugin and retry
   - Use basic context gathering without the structured brainstorming workflow

**If superpowers IS installed:** Continue to the next check.

#### 0.2: MLX-Hub Plugin (Optional)

The mlx-hub plugin enables local model acceleration for faster task processing.

**If mlx-hub is NOT installed:**

1. Inform the user: "The mlx-hub plugin enables local MLX model acceleration but isn't currently installed."
2. Offer to install it:
   ```bash
   claude plugin add https://github.com/sloanahrens/mlx-hub-claude-plugin
   ```
3. If the user declines or installation fails, **continue without local model acceleration** - the brainstorming session works fine using Claude alone, just without the speed boost from local models.

**If mlx-hub IS installed:** Proceed to Step 3 (Local Model Acceleration) when appropriate.

---

### Step 1: Resolve Repository

Follow repo selection from `_shared-repo-logic.md`:
1. Read `config.yaml` for base path and repo definitions
2. Match `$ARGUMENTS` to repo name or alias
3. If no repo recognized, ask which repo the task relates to
4. Confirm: "Brainstorming for: <repo-name>"

### Step 2: Load Repo Context

```bash
pwd  # Verify again before repo commands
cd <base_path>/<repo> && git status
cd <base_path>/<repo> && git log --oneline -5
```

Read: `<repo>/CLAUDE.md`, `README.md`, `docs/overview.md`

### Step 3: Local Model Acceleration

Use local MLX models to speed up tasks. Two tiers available:

| Tier | Model | Size | Speed | Use For |
|------|-------|------|-------|---------|
| **Fast** | `mlx-community/Llama-3.2-1B-Instruct-4bit` | 0.7GB | ~100 tok/s | Simple tasks, bulk operations |
| **Quality** | `mlx-community/Llama-3.3-70B-Instruct-8bit` | 70GB | ~15 tok/s | Complex reasoning, code generation |

**Routing rules:**

| Task | Tier | Review |
|------|------|--------|
| One-line summaries | Fast | No |
| List/enumerate items | Fast | No |
| Format/restructure text | Fast | No |
| File summaries (detailed) | Quality | No |
| Code generation | Quality | Yes - Claude reviews |
| Test stubs | Quality | Yes - Claude reviews |
| Doc drafts | Quality | Light review |
| Explore approaches | Quality | No |

**Keep on Claude:**
- Architectural decisions
- Security-sensitive code
- Complex debugging
- Final review of local-generated code
- Orchestration and synthesis

**Examples:**

```python
# Fast tier - simple extraction
mlx_infer(
  model_id="mlx-community/Llama-3.2-1B-Instruct-4bit",
  prompt="List the function names in this file:\n\n{content}",
  max_tokens=128
)

# Quality tier - code generation
mlx_infer(
  model_id="mlx-community/Llama-3.3-70B-Instruct-8bit",
  prompt="Write a TypeScript function that validates email format.",
  max_tokens=256
)
# Then: Claude reviews and refines
```

### Step 4: Run Brainstorming

Invoke `/superpowers:brainstorming` with:
- Selected repo name and path
- Task/topic from `$ARGUMENTS`
- Key context from repo's CLAUDE.md
- Current git status
- **Awareness of local model for acceleration**

---

## Documentation Location

When creating documentation:

| Type | Location |
|------|----------|
| Technical reviews | `<repo>/docs/tech-review.md` |
| Design docs | `<repo>/docs/plans/<date>-<topic>-design.md` |
| Implementation plans | `<repo>/docs/plans/<date>-<topic>-plan.md` |

**If unsure where docs belong, ASK the user.**

---

## Examples

```bash
/super my-app add user authentication   # Brainstorm for my-app repo
/super optimize database queries        # Prompts for repo selection
/super pulumi                           # Start brainstorming for infra repo
```

## Local Model Tips

- **Prompt tersely** - Llama responds well to direct instructions
- **Set appropriate max_tokens** - 256 for small functions, 1024 for larger drafts
- **Review code output** - Always have Claude review before committing
- **Use for parallelism** - Draft multiple approaches while Claude analyzes
