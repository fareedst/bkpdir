#!/usr/bin/env bash
# Run tied-yaml MCP validation (yaml_index_validate + tied_validate_consistency) via stdio,
# then formal pseudocode verification (specctl Layer B).
# See scripts/validate-tied-mcp.mjs for env vars and .cursor/mcp.json discovery.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TIED_EXIT=0
node "$ROOT/scripts/validate-tied-mcp.mjs" "$@" || TIED_EXIT=$?
if [[ -f "$ROOT/scripts/run-spec-verification.sh" ]]; then
  # Layer A already ran above; avoid run-spec-verification → validate-tied-mcp recursion.
  SKIP_TIED_MCP=1 bash "$ROOT/scripts/run-spec-verification.sh"
  SPEC_EXIT=$?
else
  SPEC_EXIT=0
fi
if [[ "$TIED_EXIT" -ne 0 || "$SPEC_EXIT" -ne 0 ]]; then
  exit 1
fi
