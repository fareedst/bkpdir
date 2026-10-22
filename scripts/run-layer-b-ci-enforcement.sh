#!/usr/bin/env bash
# Opt-in Layer B enforcement for CI and local gates (changed sidecars or full wave).
# [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

BATCH="$ROOT/scripts/run-pseudocode-layer-b-batch.sh"
PLANNER="$ROOT/scripts/layer_b_ci_enforcement.py"
REGISTRY="${LAYER_B_WAVE_REGISTRY:-$ROOT/tied-project/spec/layer-b-wave-registry.yaml}"
MODE="${LAYER_B_VALIDATE_MODE:-full}"
OUT_DIR="${LAYER_B_REPORT_DIR:-$ROOT/tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/pseudocode-validate-ci}"
BASE_REF="${LAYER_B_DIFF_BASE:-HEAD~1}"
HEAD_REF="${LAYER_B_DIFF_HEAD:-HEAD}"

if [[ -n "${GITHUB_BASE_REF:-}" && "${GITHUB_EVENT_NAME:-}" == "pull_request" ]]; then
  git fetch --no-tags --depth=1 origin "${GITHUB_BASE_REF}" 2>/dev/null || true
  BASE_REF="origin/${GITHUB_BASE_REF}"
fi

planner_mode="if_changed"
wave=""
dry_run=0

usage() {
  echo "usage: run-layer-b-ci-enforcement.sh [--if-changed|--wave wN] [--dry-run]" >&2
  echo "  env: RUN_LAYER_B=1 RUN_LAYER_B_WAVE=w0|w1|w2|w3 RUN_LAYER_B_IF_CHANGED=1" >&2
  exit 2
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --if-changed) planner_mode="if_changed"; shift ;;
    --wave)
      wave="$(echo "${2:-}" | tr '[:upper:]' '[:lower:]')"
      [[ "$wave" =~ ^w[0-3]$ ]] || usage
      planner_mode="wave"
      shift 2
      ;;
    --dry-run) dry_run=1; shift ;;
    -h|--help) usage ;;
    *) echo "unknown arg: $1" >&2; usage ;;
  esac
done

if [[ "${RUN_LAYER_B:-}" == "1" ]]; then
  if [[ "${RUN_LAYER_B_WAVE:-}" =~ ^[wW][0-3]$ ]]; then
    wave="$(echo "$RUN_LAYER_B_WAVE" | tr '[:upper:]' '[:lower:]')"
    planner_mode="wave"
  else
    planner_mode="force_changed"
  fi
elif [[ "${RUN_LAYER_B_IF_CHANGED:-}" == "1" ]]; then
  planner_mode="if_changed"
fi

planner_args=(
  --registry "$REGISTRY"
  --base-ref "$BASE_REF"
  --head-ref "$HEAD_REF"
  --mode "$planner_mode"
)
if [[ -n "$wave" ]]; then
  planner_args+=(--wave "$wave")
fi

if [[ "$dry_run" -eq 1 ]]; then
  python3 "$PLANNER" "${planner_args[@]}" --dry-run --json
  exit 0
fi

mapfile -t ENFORCE_TOKENS < <(python3 "$PLANNER" "${planner_args[@]}")

if [[ ${#ENFORCE_TOKENS[@]} -eq 0 ]]; then
  echo "DEBUG: Layer B CI enforcement — no tokens selected (mode=$planner_mode wave=${wave:-none})"
  exit 0
fi

echo "DEBUG: Layer B CI enforcement — ${#ENFORCE_TOKENS[@]} token(s): ${ENFORCE_TOKENS[*]}"

batch_args=(--mode "$MODE" --out-dir "$OUT_DIR")
for t in "${ENFORCE_TOKENS[@]}"; do
  batch_args+=(--token "$t" --enforce "$t")
done

"$BATCH" "${batch_args[@]}"
