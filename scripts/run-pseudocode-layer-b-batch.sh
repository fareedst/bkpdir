#!/usr/bin/env bash
# Layer B batch driver via pseudocode_validate (MCP checklist authority).
# [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
TIED_CLI="${TIED_CLI:-$ROOT/.cursor/skills/tied-yaml/scripts/tied-cli.sh}"
export TIED_BASE_PATH="${TIED_BASE_PATH:-$ROOT/tied-project}"

MODE="structural"
OUT_DIR="${LAYER_B_REPORT_DIR:-$ROOT/tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/pseudocode-validate}"
declare -a TOKENS=()
declare -a ENFORCE=()
ENFORCE_ALL=0

usage() {
  echo "usage: run-pseudocode-layer-b-batch.sh [--mode structural|full] [--out-dir PATH] [--token IMPL-…]… [--enforce IMPL-…]… [--wave w0|w1|w2|w3]" >&2
  exit 2
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --mode) MODE="${2:-}"; shift 2 ;;
    --out-dir) OUT_DIR="${2:-}"; shift 2 ;;
    --token) TOKENS+=("${2:-}"); shift 2 ;;
    --enforce)
      if [[ $# -ge 2 && -n "${2:-}" && "${2:-}" != --* ]]; then
        ENFORCE+=("${2}"); shift 2
      else
        ENFORCE_ALL=1; shift
      fi
      ;;
    --wave)
      wave="${2:-}"
      shift 2
      wave_lc="$(echo "$wave" | tr '[:upper:]' '[:lower:]')"
      registry="$ROOT/tied-project/spec/layer-b-wave-registry.yaml"
      if [[ ! -f "$registry" ]]; then
        echo "DIAGNOSTIC: missing wave registry $registry" >&2
        exit 2
      fi
      mapfile -t _wave_tokens < <(python3 - <<PY
import sys, yaml
from pathlib import Path
reg = yaml.safe_load(Path("$registry").read_text())
w = reg.get("waves", {}).get("$wave_lc", {})
tokens = w.get("tokens") or []
if not tokens:
    print("DIAGNOSTIC: wave $wave_lc has no tokens in registry", file=sys.stderr)
    sys.exit(2)
for t in tokens:
    print(t)
PY
)
      if [[ ${#_wave_tokens[@]} -eq 0 ]]; then
        exit 2
      fi
      TOKENS+=("${_wave_tokens[@]}")
      ;;
    -h|--help) usage ;;
    *) echo "unknown arg: $1" >&2; usage ;;
  esac
done

if [[ ${#TOKENS[@]} -eq 0 ]]; then
  usage
fi

if [[ "$MODE" != "structural" && "$MODE" != "full" ]]; then
  echo "DIAGNOSTIC: invalid --mode $MODE" >&2
  exit 2
fi

mkdir -p "$OUT_DIR"

# Repo-local tied-project wins when env TIED_BASE_PATH points at another clone (AGENTS.md RISK-010).
_canonical_tied="$ROOT/tied-project"
if [[ ${#TOKENS[@]} -gt 0 && -d "$_canonical_tied/implementation-decisions" ]]; then
  _probe_rel="implementation-decisions/${TOKENS[0]}-pseudocode.md"
  if [[ ! -f "$TIED_BASE_PATH/$_probe_rel" && -f "$_canonical_tied/$_probe_rel" ]]; then
    echo "DIAGNOSTIC: TIED_BASE_PATH=$TIED_BASE_PATH missing $_probe_rel; using $_canonical_tied" >&2
    export TIED_BASE_PATH="$_canonical_tied"
  fi
fi

is_enforced() {
  local t="$1"
  if [[ "$ENFORCE_ALL" -eq 1 ]]; then
    return 0
  fi
  for e in "${ENFORCE[@]}"; do
    [[ "$e" == "$t" ]] && return 0
  done
  return 1
}

fail=0
for token in "${TOKENS[@]}"; do
  sidecar_rel="implementation-decisions/${token}-pseudocode.md"
  sidecar="$TIED_BASE_PATH/$sidecar_rel"
  if [[ ! -f "$sidecar" && -f "$_canonical_tied/$sidecar_rel" ]]; then
    sidecar="$_canonical_tied/$sidecar_rel"
  fi
  if [[ ! -f "$sidecar" ]]; then
    sidecar="$ROOT/tied/$sidecar_rel"
  fi
  if [[ ! -f "$sidecar" ]]; then
    echo "DIAGNOSTIC: missing sidecar for $token" >&2
    if is_enforced "$token"; then fail=1; fi
    continue
  fi
  report="$OUT_DIR/${token}.json"
  args_file="$(mktemp)"
  python3 - <<PY > "$args_file"
import json
print(json.dumps({
  "token": "$token",
  "essence_pseudocode_path": "$sidecar_rel",
  "mode": "$MODE",
}))
PY
  echo "DEBUG: Layer B validate $token mode=$MODE -> $report"
  if ! "$TIED_CLI" pseudocode_validate "@$args_file" > "$report"; then
    echo "FAIL: pseudocode_validate RPC failed for $token" >&2
    rm -f "$args_file"
    if is_enforced "$token"; then fail=1; fi
    continue
  fi
  rm -f "$args_file"
  token_errors=0
  if ! python3 - <<PY
import json, sys
r = json.load(open("$report"))
errs = [d for d in (r.get("diagnostics") or []) if d.get("severity") == "error"]
if r.get("ok") is not True or errs:
    print(f"DIAGNOSTIC: $token validation not clean ok={r.get('ok')} errors={len(errs)}", file=sys.stderr)
    sys.exit(1)
PY
  then
    token_errors=1
  fi
  if [[ "$token_errors" -ne 0 ]] && is_enforced "$token"; then
    fail=1
  fi
done

if [[ "$fail" -ne 0 ]]; then
  echo "DIAGNOSTIC: Layer B batch failed (--enforce)" >&2
  exit 1
fi
echo "DIAGNOSTIC: Layer B batch finished for ${#TOKENS[@]} token(s)"
