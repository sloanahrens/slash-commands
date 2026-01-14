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

# Count recent insights notes (last 7 days)
recent_insights=0
if [ -d "${NOTES_DIR}/insights" ]; then
    recent_insights=$(find "${NOTES_DIR}/insights" -name "*.md" -mtime -7 2>/dev/null | wc -l | tr -d ' ')
fi

# Count patterns
pattern_count=0
if [ -d "${PATTERNS_DIR}" ]; then
    pattern_count=$(find "${PATTERNS_DIR}" -name "*.md" ! -name "README.md" 2>/dev/null | wc -l | tr -d ' ')
fi

# Helper: relative time from file modification
relative_time() {
    local file="$1"
    local now=$(date +%s)
    local mtime=$(stat -f %m "$file" 2>/dev/null || stat -c %Y "$file" 2>/dev/null)
    local diff=$((now - mtime))

    if [ "$diff" -lt 3600 ]; then
        echo "$((diff / 60))m ago"
    elif [ "$diff" -lt 86400 ]; then
        echo "$((diff / 3600))h ago"
    elif [ "$diff" -lt 172800 ]; then
        echo "yesterday"
    else
        echo "$((diff / 86400))d ago"
    fi
}

# Find recent session notes (top 3)
recent_sessions=""
session_details=""
if [ -d "${NOTES_DIR}/sessions" ]; then
    while IFS= read -r session_file; do
        [ -z "$session_file" ] && continue
        repo=$(basename "$session_file" .md | sed 's/^[0-9-]*-//')
        age=$(relative_time "$session_file")
        if [ -z "$recent_sessions" ]; then
            recent_sessions="${repo} (${age})"
            session_details="**${repo}** (${age})"
        else
            recent_sessions="${recent_sessions}, ${repo} (${age})"
            session_details="${session_details}, **${repo}** (${age})"
        fi
    done < <(ls -t "${NOTES_DIR}/sessions"/*.md 2>/dev/null | head -3)
fi

# Build context message
context_parts=""

if [ "$recent_insights" -gt 0 ]; then
    context_parts="${context_parts}**${recent_insights} recent insights note(s)** in ~/.claude/notes/insights/ (last 7 days)"$'\n'
fi

if [ "$pattern_count" -gt 0 ]; then
    context_parts="${context_parts}**${pattern_count} pattern(s)** available in ~/.claude/patterns/"$'\n'
fi

# Only output if there's something to report
if [ -n "$context_parts" ] || [ -n "$recent_sessions" ]; then
    # Short message for user visibility
    user_message="📚 "

    if [ "$pattern_count" -gt 0 ]; then
        user_message="${user_message}${pattern_count} patterns"
    fi

    if [ "$recent_insights" -gt 0 ]; then
        if [ "$pattern_count" -gt 0 ]; then
            user_message="${user_message}, "
        fi
        user_message="${user_message}${recent_insights} recent notes"
    fi

    if [ -n "$recent_sessions" ]; then
        if [ "$pattern_count" -gt 0 ] || [ "$recent_insights" -gt 0 ]; then
            user_message="${user_message}. "
        fi
        user_message="${user_message}Sessions: ${recent_sessions}"
    fi

    # Full context for Claude
    claude_context="## Agent Memory Available"$'\n\n'

    if [ -n "$context_parts" ]; then
        claude_context="${claude_context}${context_parts}"$'\n'
    fi

    if [ -n "$session_details" ]; then
        claude_context="${claude_context}**Recent sessions:** ${session_details}"$'\n\n'
    fi

    claude_context="${claude_context}Run \`/prime <repo>\` to load relevant notes before starting work."

    escaped_user=$(escape_for_json "$user_message")
    escaped_claude=$(escape_for_json "$claude_context")

    cat <<EOF
{
  "systemMessage": "${escaped_user}",
  "hookSpecificOutput": {
    "hookEventName": "SessionStart",
    "additionalContext": "${escaped_claude}"
  }
}
EOF
else
    # No notes yet - just remind about the system
    cat <<EOF
{
  "systemMessage": "💡 Use \`/capture-insight\` after issues and \`/capture-session <repo>\` to track progress."
}
EOF
fi

exit 0
