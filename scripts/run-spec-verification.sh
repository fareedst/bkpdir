#!/usr/bin/env bash
# Formal pseudocode verification: specctl Layer B + go test conformance.
# [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
# shellcheck source=scripts/lib/formal_spec_registry.sh
source "$ROOT/scripts/lib/formal_spec_registry.sh"
FORMAL_SPEC_REGISTRY_ROOT="$ROOT"

echo "DEBUG: package-level lead hygiene audit"
python3 scripts/audit_package_level_leads.py --fail-on-suspect

echo "DEBUG: D15 redundant IMPL banner audit"
python3 scripts/strip_redundant_impl_banners.py --fail-on-redundant

echo "DEBUG: specctl validate"
go run ./cmd/specctl validate tied/implementation-decisions/*-pseudocode.md

echo "DEBUG: specctl matrix (report)"
go run ./cmd/specctl matrix tied/implementation-decisions/*-pseudocode.md

echo "DEBUG: specctl matrix --check-coverage (formal_spec only)"
go run ./cmd/specctl matrix --check-coverage tied/implementation-decisions/*-pseudocode.md

echo "DEBUG: specctl check-leads (all formal_spec registry sidecars)"
REGISTRY_COUNT="$(formal_spec_registry_count)"
if [[ "$REGISTRY_COUNT" -eq 0 ]]; then
  echo "DIAGNOSTIC: formal-spec-registry parse returned 0 tokens"
  exit 1
fi
FORMAL_SIDECARS=()
while IFS= read -r token; do
  [[ -z "$token" ]] && continue
  path="tied/implementation-decisions/${token}-pseudocode.md"
  if [[ -f "$path" ]]; then
    FORMAL_SIDECARS+=("$path")
  fi
done < <(formal_spec_registry_tokens)
if ((${#FORMAL_SIDECARS[@]} == 0)); then
  echo "DIAGNOSTIC: formal-spec-registry listed ${REGISTRY_COUNT} tokens but no sidecar files resolved"
  exit 1
fi
go run ./cmd/specctl check-leads "${FORMAL_SIDECARS[@]}"

echo "DEBUG: go test spec conformance"
go test ./internal/specparse/... ./internal/specmodel/... ./test/specconformance/... -count=1

if [[ "${RUN_LAYER_C:-}" == "1" ]]; then
  echo "DEBUG: Layer C pseudocode_analyze (Tier A; set RUN_LAYER_C=1)"
  "$ROOT/scripts/run_layer_c_gate.sh"
elif git diff --name-only HEAD 2>/dev/null | grep -qE '\-pseudocode\.md$'; then
  echo "DEBUG: changed sidecars detected — optional Layer C (set RUN_LAYER_C=1 to enforce)"
fi

if [[ "${SKIP_TIED_MCP:-}" != "1" && -x "$ROOT/scripts/validate-tied-mcp.sh" ]]; then
  echo "DEBUG: tied Layer A (set SKIP_TIED_MCP=1 to skip)"
  timeout 60 "$ROOT/scripts/validate-tied-mcp.sh" || echo "WARNING: tied_validate_consistency skipped or timed out"
fi

echo ""
if [[ -x "$ROOT/scripts/spec-corpus-status.sh" ]]; then
  "$ROOT/scripts/spec-corpus-status.sh"
else
  echo "DIAGNOSTIC: run-spec-verification completed"
fi
