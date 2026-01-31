#!/usr/bin/env bash
# SessionStart hook - reminds about session workflow
# Session notes now live in each repo's .claude/sessions/ directory

set -euo pipefail

# Simple reminder about the workflow
cat <<EOF
{
  "systemMessage": "💡 Use \`/prime <repo>\` to load previous session context before starting work."
}
EOF

exit 0
