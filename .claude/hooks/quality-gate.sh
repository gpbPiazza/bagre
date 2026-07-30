#!/bin/bash
# Stop hook gate: run `just quality` only if Claude edited a .go file this session.
# Receives the hook payload JSON on stdin and inspects the session transcript
# for Edit/Write tool calls targeting .go files.

input=$(cat)
transcript_path=$(printf '%s' "$input" | jq -r '.transcript_path // empty')

# No transcript to inspect — nothing was edited by Claude, skip the suite.
[ -n "$transcript_path" ] && [ -f "$transcript_path" ] || exit 0

go_edits=$(jq -r '
  select(.type == "assistant")
  | .message.content[]?
  | select(.type == "tool_use")
  | select(.name == "Edit" or .name == "Write" or .name == "MultiEdit" or .name == "NotebookEdit")
  | .input.file_path // empty
' "$transcript_path" 2>/dev/null | grep -c '\.go$')

[ "${go_edits:-0}" -gt 0 ] || exit 0

just quality >&2 || exit 2
