#!/usr/bin/env bash
# Shared formal-spec-registry parser for status/verification scripts.
# [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
set -euo pipefail

_formal_spec_registry_root() {
  if [[ -n "${FORMAL_SPEC_REGISTRY_ROOT:-}" ]]; then
    printf '%s' "$FORMAL_SPEC_REGISTRY_ROOT"
    return
  fi
  local lib_dir
  lib_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  cd "$lib_dir/../.." && pwd
}

_formal_spec_registry_py() {
  local root
  root="$(_formal_spec_registry_root)"
  FORMAL_SPEC_REGISTRY_PATH="${FORMAL_SPEC_REGISTRY_PATH:-$root/tied/spec/formal-spec-registry.yaml}" \
    python3 "$root/scripts/formal_spec_registry.py" "$@"
}

formal_spec_registry_tokens() {
  _formal_spec_registry_py tokens
}

formal_spec_registry_count() {
  _formal_spec_registry_py count
}
