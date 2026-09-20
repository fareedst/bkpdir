#!/usr/bin/env bash
# Layer C pseudocode_analyze gate for Tier A runtime sidecars (opt-in).
# [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
TIED_CLI="$ROOT/.cursor/skills/tied-yaml/scripts/tied-cli.sh"
OUT_DIR="${LAYER_C_REPORT_DIR:-$ROOT/working/REQ-PSEUDOCODE_FORMAL_VERIFICATION/tied-canon/pseudocode-analysis}"
mkdir -p "$OUT_DIR"

TIER_A=(
  IMPL-ATOMIC_OPS
  IMPL-CONFIG_STRUCT
  IMPL-CFG_006
  IMPL-FILE_OPERATIONS
  IMPL-ZIP_FORMAT
)

fail=0
for token in "${TIER_A[@]}"; do
  sidecar="tied/implementation-decisions/${token}-pseudocode.md"
  [[ -f "$sidecar" ]] || continue
  report="$OUT_DIR/${token}-layer-c.json"
  args_file="$(mktemp)"
  python3 - <<PY > "$args_file"
import json
from pathlib import Path
import yaml
root = Path("$ROOT")
reg = yaml.safe_load((root / "tied/semantic-tokens.yaml").read_text())
known = sorted(k for k in reg if k.startswith(("REQ-", "ARCH-", "IMPL-")))
print(json.dumps({
  "token": "$token",
  "essence_pseudocode_path": "implementation-decisions/${token}-pseudocode.md",
  "gate_mode": True,
  "known_tokens": known,
}))
PY
  echo "DEBUG: Layer C analyze $token -> $report"
  if ! "$TIED_CLI" pseudocode_analyze "@$args_file" > "$report"; then
    echo "FAIL: pseudocode_analyze RPC failed for $token"
    rm -f "$args_file"
    fail=1
    continue
  fi
  rm -f "$args_file"
  if ! python3 - <<PY
import json, sys
r = json.load(open("$report"))
ok = r.get("ok") is True and r.get("gate_mode_applied") is True
if not ok:
    di = r.get("diagnostics") or []
    errs = [d for d in di if d.get("severity") == "error"]
    print(f"FAIL $token: ok={r.get('ok')} errors={len(errs)}", file=sys.stderr)
    for d in errs[:8]:
        print(f"  {d.get('code')}: {d.get('message')}", file=sys.stderr)
    sys.exit(1)
print(f"PASS $token Layer C gate")
PY
  then
    fail=1
  fi
done

if [[ "$fail" -ne 0 ]]; then
  echo "DIAGNOSTIC: Layer C gate failed (see reports under $OUT_DIR)"
  exit 1
fi
echo "DIAGNOSTIC: Layer C gate passed for ${#TIER_A[@]} Tier A sidecars"
