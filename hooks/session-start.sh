#!/usr/bin/env bash
# SessionStart hook for knowledge-primed sessions
# Surfaces recent notes, patterns, and last session for continuity

set -euo pipefail

NOTES_DIR="${HOME}/.claude/notes"
PATTERNS_DIR="${HOME}/.claude/patterns"

escape_for_json() {
    local input="$1"
    local output=""
    local i char
    for (( i=0; i<${#input}; i++ )); do
        char="${input:$i:1}"
        case "$char" in
            '\') output+='\' ;;
            '"') output+='\"' ;;
            $'\n') output+='\n' ;;
            $'\r') output+='\r' ;;
            $'\t') output+='\t' ;;
            *) output+="$char" ;;
        esac
    done
    printf '%s' "$output"
}

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

# Find most recent session note
last_session=""
last_session_repo=""
if [ -d "${NOTES_DIR}/sessions" ]; then
    last_session=$(ls -t "${NOTES_DIR}/sessions"/*.md 2>/dev/null | head -1)
    if [ -n "$last_session" ]; then
        # Extract repo from filename (YYYY-MM-DD-<repo>.md)
        last_session_repo=$(basename "$last_session" .md | sed 's/^[0-9-]*-//')
    fi
fi

# Build context message
context_parts=""

if [ "$recent_hindsight" -gt 0 ]; then
    context_parts="${context_parts}**${recent_hindsight} recent hindsight note(s)** in ~/.claude/notes/hindsight/ (last 7 days)"$'\n'
fi

if [ "$pattern_count" -gt 0 ]; then
    context_parts="${context_parts}**${pattern_count} pattern(s)** available in ~/.claude/patterns/"$'\n'
fi

# Only output if there's something to report
if [ -n "$context_parts" ] || [ -n "$last_session" ]; then
    message="## Agent Memory Available"$'\n\n'

    if [ -n "$context_parts" ]; then
        message="${message}${context_parts}"$'\n'
    fi

    if [ -n "$last_session" ]; then
        message="${message}**Last session:** ${last_session_repo} — Run \`/find-tasks ${last_session_repo}\` to continue"$'\n\n'
    fi

    message="${message}Run \`/prime <repo>\` to load relevant notes before starting work."

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
    "additionalContext": "**Tip:** Use \`/capture-hindsight\` after encountering issues and \`/capture-session <repo>\` to track progress."
  }
}
EOF
fi

exit 0
