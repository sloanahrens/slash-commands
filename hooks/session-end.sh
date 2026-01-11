#!/usr/bin/env bash
# Stop hook for automatic knowledge capture
# Directs Claude to run capture commands before ending

set -euo pipefail

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

message="MANDATORY: Before completing this response, you MUST run \`/capture-session\` using the Skill tool to save this session's progress. If there were any errors, debugging, or non-obvious solutions during this session, also run \`/capture-hindsight\`. Do not ask — just run them."

escaped_message=$(escape_for_json "$message")

cat <<EOF
{
  "systemMessage": "${escaped_message}"
}
EOF

exit 0
