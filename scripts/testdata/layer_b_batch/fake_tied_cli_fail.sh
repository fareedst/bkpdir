#!/usr/bin/env bash
# Test double: pseudocode_validate returns ok:false when token is IMPL-FAIL
set -euo pipefail
tool="${1:-}"
arg="${2:-}"
body=""
if [[ "$arg" == @* ]]; then
  body="$(cat "${arg#@}")"
else
  body="$arg"
fi
if [[ "$tool" == "pseudocode_validate" ]]; then
  if echo "$body" | grep -q 'IMPL-FAIL'; then
    echo '{"ok":false,"diagnostics":[{"severity":"error","code":"SHAPE-003","message":"stub fail"}]}'
  else
    echo '{"ok":true,"diagnostics":[]}'
  fi
  exit 0
fi
echo '{"ok":false,"error":"unknown tool"}' >&2
exit 1
