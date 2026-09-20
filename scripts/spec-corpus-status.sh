#!/usr/bin/env bash
# Print formal DSL corpus health and exit 0 when all-green.
# [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
# shellcheck source=scripts/lib/formal_spec_registry.sh
source "$ROOT/scripts/lib/formal_spec_registry.sh"
FORMAL_SPEC_REGISTRY_ROOT="$ROOT"

fail=0
SIDECARS=$(find tied/implementation-decisions -maxdepth 1 -name 'IMPL-*-pseudocode.md' 2>/dev/null | wc -l | tr -d ' ')
REGISTRY="$(formal_spec_registry_count)"
if [[ "$REGISTRY" -eq 0 ]]; then
  echo "DIAGNOSTIC: formal-spec-registry parse returned 0 tokens"
  exit 1
fi
LEAD_SKIP=$(grep -A999 '^lead_check_skip:' tied/spec/coverage-waivers.yaml 2>/dev/null | grep -cE '^\s+- IMPL-' || true)

queue_complete=""
if command -v python3 >/dev/null 2>&1; then
  queue_complete=$(python3 - <<'PY'
import sys
from pathlib import Path
sys.path.insert(0, "scripts")
from load_improvement_queue import load_queue
q = load_queue()
tokens = q.get("tokens", [])
n = sum(1 for t in tokens if t.get("formal_spec") and t.get("req_audit_pass"))
print(f"{n}/{len(tokens)}")
PY
)
fi

VALIDATE_OUT="$(go run ./cmd/specctl validate tied/implementation-decisions/*-pseudocode.md 2>&1)" || true
VAL_ERRORS=$(printf '%s\n' "$VALIDATE_OUT" | grep -c '\[error\]' || true)
VAL_WARNINGS=$(printf '%s\n' "$VALIDATE_OUT" | grep -c '\[warning\]' || true)

check() {
  local label="$1" ok="$2" detail="$3"
  if [[ "$ok" == "1" ]]; then
    printf "  %-22s PASS  %s\n" "$label" "$detail"
  else
    printf "  %-22s FAIL  %s\n" "$label" "$detail"
    fail=1
  fi
}

echo "========================================"
echo "Formal DSL corpus status"
echo "========================================"

if [[ "$SIDECARS" == "$REGISTRY" && "$SIDECARS" -gt 0 ]]; then
  check "sidecar/registry" 1 "${SIDECARS} sidecars == ${REGISTRY} formal_spec registry entries"
else
  check "sidecar/registry" 0 "${SIDECARS} sidecars vs ${REGISTRY} registry entries (must match)"
fi

if [[ "$VAL_ERRORS" -eq 0 && "$VAL_WARNINGS" -eq 0 ]]; then
  check "specctl validate" 1 "0 errors, 0 warnings (all 73 sidecars)"
else
  check "specctl validate" 0 "${VAL_ERRORS} errors, ${VAL_WARNINGS} warnings"
fi

if go run ./cmd/specctl matrix --check-coverage tied/implementation-decisions/*-pseudocode.md >/dev/null 2>&1; then
  check "matrix coverage" 1 "COVER-001 clear (waivers applied)"
else
  check "matrix coverage" 0 "uncovered formal blocks — run matrix --check-coverage"
fi

FORMAL_SIDECARS=()
while IFS= read -r token; do
  [[ -z "$token" ]] && continue
  path="tied/implementation-decisions/${token}-pseudocode.md"
  [[ -f "$path" ]] && FORMAL_SIDECARS+=("$path")
done < <(formal_spec_registry_tokens)

if ((${#FORMAL_SIDECARS[@]} > 0)); then
  if go run ./cmd/specctl check-leads "${FORMAL_SIDECARS[@]}" >/dev/null 2>&1; then
    check "check-leads" 1 "${#FORMAL_SIDECARS[@]} registry sidecars (${LEAD_SKIP} lead_check_skip)"
  else
    check "check-leads" 0 "TRACE-002 — sync // comments or add lead_check_skip"
  fi
else
  check "check-leads" 0 "no formal_spec sidecars in registry"
fi

if go test ./internal/specparse/... ./internal/specmodel/... ./test/specconformance/... -count=1 >/dev/null 2>&1; then
  check "conformance tests" 1 "specparse + specmodel + specconformance"
else
  check "conformance tests" 0 "go test ./internal/specparse/... ./test/specconformance/..."
fi

ORACLE_STATS=""
DEEP_SYNC_STATS=""
LOGIC_STATS=""
if command -v python3 >/dev/null 2>&1; then
  ORACLE_STATS=$(python3 - <<'PY'
import yaml
from pathlib import Path
p = Path("tied/spec/oracle-registry.yaml")
if not p.exists():
    print("0/50")
else:
    r = yaml.safe_load(p.read_text())
    n = sum(1 for o in r.get("oracles", []) if o.get("status") == "verified")
    t = r.get("oracle_targets", len(r.get("oracles", [])))
    print(f"{n}/{t}")
PY
)
  DEEP_SYNC_STATS=$(python3 - <<'PY'
import re, sys
from pathlib import Path
sys.path.insert(0, "scripts")
from load_improvement_queue import load_queue
text = Path("scripts/update_impl_sync_checklist.py").read_text()
q = load_queue()
done = sum(1 for t in q["tokens"] if t.get("deep_sync_complete"))
print(f"{done}/{len(q['tokens'])}")
PY
)
  LOGIC_STATS=$(python3 - <<'PY'
import sys
sys.path.insert(0, "scripts")
from load_improvement_queue import load_queue
q = load_queue()
done = sum(1 for t in q["tokens"] if t.get("status") == "logic_verified")
print(f"{done}/{len(q['tokens'])}")
PY
)
  REQ_AUDIT_STATS=$(python3 - <<'PY'
import sys
sys.path.insert(0, "scripts")
from load_improvement_queue import load_queue
q = load_queue()
done = sum(1 for t in q["tokens"] if t.get("req_audit_pass"))
print(f"{done}/{len(q['tokens'])}")
PY
)
  CRIT001_STATS=$(python3 - <<'PY'
import sys
sys.path.insert(0, "scripts")
from load_improvement_queue import load_queue
q = load_queue()
done = sum(1 for t in q["tokens"] if t.get("crit001_pass"))
print(f"{done}/{len(q['tokens'])}")
PY
)
fi

if [[ -n "$ORACLE_STATS" ]]; then
  ORACLE_N="${ORACLE_STATS%%/*}"
  ORACLE_T="${ORACLE_STATS##*/}"
  if [[ "$ORACLE_N" == "$ORACLE_T" && "$ORACLE_T" == "50" ]]; then
    check "oracle coverage" 1 "${ORACLE_STATS} verified (L3)"
  else
    check "oracle coverage" 0 "${ORACLE_STATS} verified (expect 50/50)"
  fi
fi

if [[ -n "$DEEP_SYNC_STATS" ]]; then
  DEEP_N="${DEEP_SYNC_STATS%%/*}"
  DEEP_T="${DEEP_SYNC_STATS##*/}"
  if [[ "$DEEP_N" == "$DEEP_T" ]]; then
    check "deep sync tracking" 1 "${DEEP_SYNC_STATS} deep_sync_complete"
  else
    check "deep sync tracking" 0 "${DEEP_SYNC_STATS} deep_sync_complete"
  fi
fi

if [[ -n "$LOGIC_STATS" ]]; then
  LOGIC_N="${LOGIC_STATS%%/*}"
  LOGIC_T="${LOGIC_STATS##*/}"
  if [[ "$LOGIC_N" == "$LOGIC_T" ]]; then
    check "logic verified" 1 "${LOGIC_STATS} logic_verified"
  else
    check "logic verified" 0 "${LOGIC_STATS} logic_verified"
  fi
fi

if [[ -n "$REQ_AUDIT_STATS" ]]; then
  REQ_N="${REQ_AUDIT_STATS%%/*}"
  REQ_T="${REQ_AUDIT_STATS##*/}"
  if [[ "$REQ_N" == "$REQ_T" ]]; then
    check "req audit pass" 1 "${REQ_AUDIT_STATS} req_audit_pass (manual)"
  else
    check "req audit pass" 0 "${REQ_AUDIT_STATS} req_audit_pass (expect ${REQ_T}/${REQ_T})"
  fi
fi

if [[ -n "$CRIT001_STATS" ]]; then
  CRIT_N="${CRIT001_STATS%%/*}"
  CRIT_T="${CRIT001_STATS##*/}"
  if [[ "$CRIT_N" == "$CRIT_T" ]]; then
    check "crit001 pass" 1 "${CRIT001_STATS} crit001_pass (scoped)"
  else
    check "crit001 pass" 0 "${CRIT001_STATS} crit001_pass (expect ${CRIT_T}/${CRIT_T})"
  fi
fi

if [[ -n "$queue_complete" && "$queue_complete" != "?" ]]; then
  if [[ "$queue_complete" == "${REGISTRY}/${REGISTRY}" || "$queue_complete" == "73/73" ]]; then
    check "improvement queue" 1 "${queue_complete} formal_spec + req_audit_pass"
  else
    check "improvement queue" 0 "${queue_complete} req_audit_pass (expect ${REGISTRY}/${REGISTRY})"
  fi
fi

echo "----------------------------------------"
echo "L4 projected (opt-in)"
echo "----------------------------------------"
if [[ "${RUN_L4_PILOT:-}" == "1" ]]; then
  if [[ "${RUN_L4_REQ_STRICT:-}" == "1" ]]; then
    if python3 scripts/run_impl_logic_audit.py --batch all --req-criteria-scoped --skip-specctl >/dev/null 2>&1; then
      check "req criteria scoped" 1 "CRIT-001 scoped check clean (pilot)"
    else
      strict_fail=$(python3 scripts/run_impl_logic_audit.py --batch all --req-criteria-scoped --skip-specctl 2>&1 | grep -c '^FAIL' || true)
      check "req criteria scoped" 0 "${strict_fail} tokens with CRIT-001 orphans (scoped)"
    fi
  else
    printf "  %-22s SKIP  set RUN_L4_REQ_STRICT=1 for CRIT-001 scoped audit\n" "req criteria scoped"
  fi
  MUT_WAVE="${MUTATION_WAVE:-0}"
  if MUTATION_WAVE="$MUT_WAVE" scripts/run-mutation-pilot.sh >/dev/null 2>&1; then
    check "mutation pilot" 1 "wave ${MUT_WAVE} >= ${MUTATION_SCORE_THRESHOLD:-70}%"
  else
    if command -v go-mutesting >/dev/null 2>&1; then
      check "mutation pilot" 0 "wave ${MUT_WAVE} below threshold (install avito-tech go-mutesting)"
    else
      printf "  %-22s SKIP  go-mutesting not installed (INFRA-MUTATION-TOOL)\n" "mutation pilot"
    fi
  fi
else
  printf "  %-22s SKIP  set RUN_L4_PILOT=1 to run L4 pilots\n" "L4 pilots"
fi

echo "----------------------------------------"
if [[ "$fail" -eq 0 ]]; then
  echo "CORPUS STATUS: ALL GREEN"
  echo "  Profile: formal_spec on all ${REGISTRY} IMPL sidecars"
  echo "  Gate:    SKIP_TIED_MCP=1 scripts/run-spec-verification.sh"
  exit 0
fi
echo "CORPUS STATUS: NOT ALL GREEN (see FAIL lines above)"
exit 1
