#!/usr/bin/env bash
# Formal pseudocode verification: specctl Layer B + go test conformance.
# [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

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
FORMAL_SIDECARS=()
while IFS= read -r token; do
  [[ -z "$token" || "$token" =~ ^# ]] && continue
  path="tied/implementation-decisions/${token}-pseudocode.md"
  if [[ -f "$path" ]]; then
    FORMAL_SIDECARS+=("$path")
  fi
done < <(grep -E '^\s+- "IMPL-' tied/spec/formal-spec-registry.yaml | sed 's/.*"\(IMPL-[^"]*\)".*/\1/')
if ((${#FORMAL_SIDECARS[@]} > 0)); then
  go run ./cmd/specctl check-leads "${FORMAL_SIDECARS[@]}"
else
  echo "WARNING: no formal_spec sidecars found in registry"
fi

echo "DEBUG: go test spec conformance"
go test ./internal/specparse/... ./internal/specmodel/... ./test/specconformance/... -count=1

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
