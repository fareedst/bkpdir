#!/usr/bin/env bash
# Run tied-yaml MCP validation (yaml_index_validate + tied_validate_consistency) via stdio.
# See scripts/validate-tied-mcp.mjs for env vars and .cursor/mcp.json discovery.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
exec node "$ROOT/scripts/validate-tied-mcp.mjs" "$@"
