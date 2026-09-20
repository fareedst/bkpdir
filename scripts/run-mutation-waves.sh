#!/usr/bin/env bash
# Run L4 mutation pilot waves 0–4 sequentially (opt-in; not global CI).
# [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

THRESHOLD="${MUTATION_SCORE_THRESHOLD:-70}"
REGISTRY="${ROOT}/tied/spec/mutation-wave-registry.yaml"

if ! command -v go-mutesting >/dev/null 2>&1; then
  echo "WARNING: go-mutesting not installed — all waves skipped (INFRA-MUTATION-TOOL)"
  exit 0
fi

echo "DEBUG: L4 mutation wave rollout (threshold ${THRESHOLD}%)"
fail=0
for wave in 0 1 2 3 4; do
  echo "----------------------------------------"
  echo "Wave ${wave}"
  if MUTATION_WAVE="$wave" MUTATION_SCORE_THRESHOLD="$THRESHOLD" scripts/run-mutation-pilot.sh; then
    echo "Wave ${wave}: PASS"
    python3 - <<PY "$REGISTRY" "$wave"
import sys
from pathlib import Path
try:
    import yaml
except ImportError:
    sys.exit(0)
path = Path(sys.argv[1])
wave_num = int(sys.argv[2])
if not path.exists():
    sys.exit(0)
data = yaml.safe_load(path.read_text())
for w in data.get("waves", []):
    if w.get("wave") == wave_num:
        w["status"] = "pilot"
        note = w.get("notes") or ""
        w["notes"] = f"{note}; wave run PASS".strip("; ")
        break
path.write_text(yaml.dump(data, default_flow_style=False, sort_keys=False, allow_unicode=True))
PY
  else
    echo "Wave ${wave}: FAIL (stop rollout; fix tests or lower threshold for pilot)"
    fail=1
    break
  fi
done

if [[ "$fail" -ne 0 ]]; then
  exit 1
fi
echo "MUTATION WAVES: all selected waves PASS"
exit 0
