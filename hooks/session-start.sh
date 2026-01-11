#!/usr/bin/env bash
# SessionStart hook for Confucius-inspired agent scaffolding
# Surfaces recent hindsight notes and reminds about /prime command

set -euo pipefail

NOTES_DIR="${HOME}/.claude/notes"
PATTERNS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")/.." && pwd)/docs/patterns"

# Count recent hindsight notes (last 7 days)
recent_hindsight=0
if [ -d "${NOTES_DIR}/hindsight" ]; then
    recent_hindsight=$(find "${NOTES_DIR}/hindsight" -name "*.md" -mtime -7 2>/dev/null | wc -l | tr -d ' ')
fi

# Count patterns
pattern_count=0
if [ -d "${PATTERNS_DIR}" ]; then
    pattern_count=$(find "${PATTERNS_DIR}" -name "*.md" ! -name "README.md" 2>/dev/null | wc -l | tr -d ' ')
fi

# Build context message
context_parts=""

if [ "$recent_hindsight" -gt 0 ]; then
    context_parts="${context_parts}**${recent_hindsight} recent hindsight note(s)** in ~/.claude/notes/hindsight/ (last 7 days)\n"
fi

if [ "$pattern_count" -gt 0 ]; then
    context_parts="${context_parts}**${pattern_count} pattern(s)** available in docs/patterns/\n"
fi

# Only output if there's something to report
if [ -n "$context_parts" ]; then
    # Escape for JSON
    escape_for_json() {
        local input="$1"
        local output=""
        local i char
        for (( i=0; i<${#input}; i++ )); do
            char="${input:$i:1}"
            case "$char" in
                $'\\') output+='\\' ;;
                '"') output+='\"' ;;
                $'\n') output+='\n' ;;
                $'\r') output+='\r' ;;
                $'\t') output+='\t' ;;
                *) output+="$char" ;;
            esac
        done
        printf '%s' "$output"
    }

    message="## Agent Memory Available\n\n${context_parts}\nRun \`/prime <repo>\` to load relevant notes before starting work."
    escaped_message=$(escape_for_json "$message")

    cat <<EOF
{
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "${escaped_message}"
  }
}
EOF
else
    # No notes yet - just remind about the system
    cat <<EOF
{
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "**Tip:** Use \`/capture-hindsight\` after encountering issues to build knowledge for future sessions."
  }
}
EOF
fi

exit 0
