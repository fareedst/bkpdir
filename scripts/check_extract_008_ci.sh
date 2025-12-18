#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

# Run import cycle and diagram check
"${SCRIPT_DIR}/check_import_cycles.sh"

# Run token traceability validation if available
if [ -x "${SCRIPT_DIR}/validate-token-traceability.sh" ]; then
  "${SCRIPT_DIR}/validate-token-traceability.sh"
else
  echo "Warning: validate-token-traceability.sh not found or not executable; skipping token validation"
fi

echo "EXTRACT-008 CI checks: OK"
