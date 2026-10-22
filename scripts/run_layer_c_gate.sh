#!/usr/bin/env bash
# Layer C pseudocode_analyze gate (opt-in); tokens from layer-b-wave-registry.yaml.
# [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
export TIED_BASE_PATH="${TIED_BASE_PATH:-$ROOT/tied-project}"
TIED_CLI="$ROOT/.cursor/skills/tied-yaml/scripts/tied-cli.sh"
OUT_DIR="${LAYER_C_REPORT_DIR:-$ROOT/tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/pseudocode-analysis-layer-c}"
REGISTRY="${LAYER_B_WAVE_REGISTRY:-$ROOT/tied-project/spec/layer-b-wave-registry.yaml}"
mkdir -p "$OUT_DIR"

LAYER_C_WAVE="${LAYER_C_WAVE:-w0}"
LAYER_C_WAVE="$(echo "$LAYER_C_WAVE" | tr '[:upper:]' '[:lower:]')"

mapfile -t TOKENS < <(
  python3 - <<PY
import sys, yaml
from pathlib import Path
reg = yaml.safe_load(Path("$REGISTRY").read_text())
waves = reg.get("waves") or {}
wave = "$LAYER_C_WAVE"
if wave == "all":
    seen = []
    for w in ("w0", "w1", "w2", "w3"):
        for t in (waves.get(w) or {}).get("tokens") or []:
            if t not in seen:
                seen.append(t)
    tokens = seen
else:
    tokens = list((waves.get(wave) or {}).get("tokens") or [])
    if not tokens:
        print(f"DIAGNOSTIC: wave {wave} has no tokens", file=sys.stderr)
        sys.exit(2)
for t in tokens:
    print(t)
PY
)

fail=0
for token in "${TOKENS[@]}"; do
  sidecar="tied-project/implementation-decisions/${token}-pseudocode.md"
  [[ -f "$sidecar" ]] || continue
  report="$OUT_DIR/${token}-layer-c.json"
  args_file="$(mktemp)"
  python3 - <<PY > "$args_file"
import json
from pathlib import Path
import yaml
root = Path("$ROOT")
reg = yaml.safe_load((root / "tied-project/semantic-tokens.yaml").read_text())
known = sorted(k for k in reg if k.startswith(("REQ-", "ARCH-", "IMPL-")))
sidecar = root / "tied-project/implementation-decisions/${token}-pseudocode.md"
print(json.dumps({
  "token": "$token",
  "essence_pseudocode_path": str(sidecar),
  "gate_mode": True,
  "known_tokens": known,
  "tied_base_path": str(root / "tied-project"),
}))
PY
  echo "DEBUG: Layer C analyze $token -> $report"
  if ! env TIED_BASE_PATH="$TIED_BASE_PATH" "$TIED_CLI" pseudocode_analyze "@$args_file" > "$report"; then
    echo "FAIL: pseudocode_analyze RPC failed for $token"
    rm -f "$args_file"
    fail=1
    continue
  fi
  rm -f "$args_file"
  if ! REPORT="$report" TOKEN="$token" python3 - <<'PY'
import json, os, sys
r = json.load(open(os.environ["REPORT"]))
ok = r.get("ok") is True and r.get("gate_mode_applied") is True
if not ok:
    di = r.get("diagnostics") or []
    errs = [d for d in di if d.get("severity") == "error"]
    print(f"FAIL {os.environ['TOKEN']}: ok={r.get('ok')} errors={len(errs)}", file=sys.stderr)
    if r.get("error"):
        print(f"  error: {r.get('error')}", file=sys.stderr)
    for d in errs[:8]:
        print(f"  {d.get('code')}: {d.get('message')}", file=sys.stderr)
    sys.exit(1)
print(f"PASS {os.environ['TOKEN']} Layer C gate")
PY
  then
    fail=1
  fi
done

if [[ "$fail" -ne 0 ]]; then
  echo "DIAGNOSTIC: Layer C gate failed (see reports under $OUT_DIR)"
  exit 1
fi
echo "DIAGNOSTIC: Layer C gate passed for ${#TOKENS[@]} sidecar(s) (wave=$LAYER_C_WAVE)"
