#!/usr/bin/env bash
# Stop hook for suggesting hindsight capture
# Checks if session had errors/retries and suggests /capture-hindsight

set -euo pipefail

# This hook runs on Stop events
# For now, just provide a gentle reminder
# Future enhancement: analyze transcript for error patterns

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

message="If this session had any tricky issues or discoveries, consider running \`/capture-hindsight\` to save them for future sessions."
escaped_message=$(escape_for_json "$message")

cat <<EOF
{
  "hookSpecificOutput": {
    "hookEventName": "Stop",
    "additionalContext": "${escaped_message}"
  }
}
EOF

exit 0
