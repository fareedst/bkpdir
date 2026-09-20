#!/usr/bin/env bash
# L4 mutation pilot: per-package scores from oracle registry or MUTATION_PKG.
# [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

THRESHOLD="${MUTATION_SCORE_THRESHOLD:-70}"
MUTEST_INSTALL="github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest"
WAVE="${MUTATION_WAVE:-all}"

usage() {
  cat <<EOF
Usage: run-mutation-pilot.sh [options]

Environment:
  MUTATION_PKG              Single package path (default: pkg/fileops)
  MUTATION_FROM_REGISTRY=1  All unique pkg_paths from tied/spec/oracle-registry.yaml
  MUTATION_WAVE             Wave filter: 0|1|2|3|4|all (default all)
  MUTATION_SCORE_THRESHOLD  Pass threshold percent (default 70)

Waves (oracle pkg_paths):
  0  pkg/fileops
  1  pkg/errors
  2  pkg/git pkg/cli
  3  pkg/resources pkg/processing pkg/testutil
  4  pkg/config pkg/formatter
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

if ! command -v go-mutesting >/dev/null 2>&1; then
  echo "WARNING: go-mutesting not installed; run: go install ${MUTEST_INSTALL}"
  echo "Skipping mutation gate (pilot waiver INFRA-MUTATION-TOOL)"
  exit 0
fi

MUTEST_BIN="$(command -v go-mutesting)"
MUTEST_FLAVOR="unknown"
if go version -m "$MUTEST_BIN" 2>/dev/null | grep -q 'github.com/avito-tech/go-mutesting'; then
  MUTEST_FLAVOR="avito-tech"
elif go version -m "$MUTEST_BIN" 2>/dev/null | grep -q 'github.com/zimmski/go-mutesting'; then
  MUTEST_FLAVOR="zimmski"
fi

if [[ "$MUTEST_FLAVOR" == "zimmski" ]]; then
  echo "ERROR: go-mutesting on PATH is the unmaintained zimmski build (incompatible with Go 1.22+)."
  echo "Reinstall: go install ${MUTEST_INSTALL}"
  exit 2
fi

parse_score() {
  local tmp="$1"
  python3 - <<'PY' "$tmp"
import re, sys
path = sys.argv[1]
text = open(path, encoding="utf-8", errors="replace").read()
panicked = "panic:" in text and "go/ast.Walk" in text
score = None
m = re.search(r"The mutation score is ([0-9.]+)", text)
if m:
    score = float(m.group(1)) * 100.0
else:
    m = re.search(r"Mutation score[^0-9]*([0-9.]+)\s*%", text, re.I)
    if m:
        score = float(m.group(1))
if score is None:
    score = 0.0
print(f"{score:.4f}|{int(panicked)}")
PY
}

run_mutesting_pkg() {
  local pkg="$1"
  local tmp
  tmp="$(mktemp)"
  echo "DEBUG: mutation pilot ${pkg}/..."
  go-mutesting "./${pkg}/..." >"$tmp" 2>&1 || true
  local result score panicked
  result="$(parse_score "$tmp")"
  score="${result%%|*}"
  panicked="${result##*|}"
  if [[ "$panicked" == "1" ]]; then
    echo "ERROR: go-mutesting crashed on ${pkg}"
    cat "$tmp"
    rm -f "$tmp"
    return 2
  fi
  echo "  ${pkg}: ${score}% (threshold ${THRESHOLD}%)"
  cat "$tmp"
  rm -f "$tmp"
  python3 - <<PY
import sys
score = float("${score}" or "0")
thresh = float("${THRESHOLD}")
sys.exit(0 if score >= thresh else 1)
PY
}

wave_packages() {
  case "$1" in
    0) echo "pkg/fileops" ;;
    1) echo "pkg/errors" ;;
    2) echo "pkg/git"; echo "pkg/cli" ;;
    3) echo "pkg/resources"; echo "pkg/processing"; echo "pkg/testutil" ;;
    4) echo "pkg/config"; echo "pkg/formatter" ;;
    all)
      echo "pkg/fileops"
      echo "pkg/errors"
      echo "pkg/git"
      echo "pkg/cli"
      echo "pkg/resources"
      echo "pkg/processing"
      echo "pkg/testutil"
      echo "pkg/config"
      echo "pkg/formatter"
      ;;
    *) echo "ERROR: unknown MUTATION_WAVE=$1" >&2; return 1 ;;
  esac
}

collect_registry_packages() {
  python3 - <<'PY'
import yaml
from pathlib import Path
p = Path("tied/spec/oracle-registry.yaml")
if not p.exists():
    raise SystemExit(0)
data = yaml.safe_load(p.read_text())
seen = set()
for o in data.get("oracles", []):
    for pkg in o.get("pkg_paths") or []:
        if pkg and pkg not in seen:
            seen.add(pkg)
            print(pkg)
PY
}

PKGS=()
if [[ "${MUTATION_FROM_REGISTRY:-}" == "1" ]]; then
  while IFS= read -r pkg; do
    [[ -n "$pkg" ]] && PKGS+=("$pkg")
  done < <(collect_registry_packages)
elif [[ -n "${MUTATION_PKG:-}" ]]; then
  PKGS=("$MUTATION_PKG")
else
  while IFS= read -r pkg; do
    [[ -n "$pkg" ]] && PKGS+=("$pkg")
  done < <(wave_packages "$WAVE")
fi

if ((${#PKGS[@]} == 0)); then
  echo "WARNING: no mutation packages selected"
  exit 0
fi

echo "DEBUG: L4 mutation pilot packages: ${PKGS[*]}"
fail=0
for pkg in "${PKGS[@]}"; do
  if ! run_mutesting_pkg "$pkg"; then
    fail=1
  fi
done

if [[ "$fail" -ne 0 ]]; then
  echo "MUTATION PILOT: FAIL (see packages below threshold ${THRESHOLD}%)"
  exit 1
fi
echo "MUTATION PILOT: PASS (${#PKGS[@]} packages >= ${THRESHOLD}%)"
exit 0
