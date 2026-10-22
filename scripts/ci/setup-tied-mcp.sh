#!/usr/bin/env bash
# Build TIED MCP for CI and export TIED_MCP_BIN / TIED_CLI (Path B Layer B enforcement).
# [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TIED_GIT_URL="${TIED_GIT_URL:-https://github.com/fareedst/tied.git}"
TIED_GIT_REF="${TIED_GIT_REF:-main}"
TIED_CHECKOUT="${TIED_CHECKOUT:-${RUNNER_TEMP:-/tmp}/tied-mcp-store}"

if [[ ! -d "$TIED_CHECKOUT/.git" ]]; then
  echo "DEBUG: cloning TIED store $TIED_GIT_URL (ref=$TIED_GIT_REF) -> $TIED_CHECKOUT"
  rm -rf "$TIED_CHECKOUT"
  git clone --depth 1 --branch "$TIED_GIT_REF" "$TIED_GIT_URL" "$TIED_CHECKOUT"
else
  echo "DEBUG: reusing TIED checkout at $TIED_CHECKOUT"
  git -C "$TIED_CHECKOUT" fetch --depth 1 origin "$TIED_GIT_REF"
  git -C "$TIED_CHECKOUT" checkout "$TIED_GIT_REF"
  git -C "$TIED_CHECKOUT" pull --ff-only origin "$TIED_GIT_REF" || true
fi

MCP_DIR="$TIED_CHECKOUT/mcp-server"
if [[ ! -f "$MCP_DIR/dist/index.js" ]]; then
  echo "DEBUG: building TIED MCP in $MCP_DIR"
  npm ci --prefix "$MCP_DIR"
  npm run build --prefix "$MCP_DIR"
fi

export TIED_REPO_ROOT="$TIED_CHECKOUT"
export TIED_MCP_BIN="$MCP_DIR/dist/index.js"
export TIED_YAML_MCP_JS="$TIED_MCP_BIN"
export TIED_BASE_PATH="${TIED_BASE_PATH:-$ROOT/tied-project}"
export TIED_METHODOLOGY_BUNDLE_PATH="${TIED_METHODOLOGY_BUNDLE_PATH:-$TIED_CHECKOUT/tied-bundle}"
export TIED_STORE_ROOT="$TIED_CHECKOUT"
export TIED_CLI="${TIED_CLI:-$TIED_CHECKOUT/tools/bundled-tied-yaml-skill/scripts/tied-cli.sh}"

if [[ ! -x "$TIED_CLI" ]]; then
  chmod +x "$TIED_CLI" 2>/dev/null || true
fi
if [[ ! -f "$TIED_MCP_BIN" ]]; then
  echo "DIAGNOSTIC: MCP build failed — missing $TIED_MCP_BIN" >&2
  exit 1
fi

echo "DEBUG: TIED MCP ready TIED_MCP_BIN=$TIED_MCP_BIN"
echo "DEBUG: TIED_BASE_PATH=$TIED_BASE_PATH"
echo "DEBUG: TIED_METHODOLOGY_BUNDLE_PATH=$TIED_METHODOLOGY_BUNDLE_PATH"

# Emit GITHUB_ENV when present (Actions).
if [[ -n "${GITHUB_ENV:-}" ]]; then
  {
    echo "TIED_REPO_ROOT=$TIED_REPO_ROOT"
    echo "TIED_MCP_BIN=$TIED_MCP_BIN"
    echo "TIED_YAML_MCP_JS=$TIED_YAML_MCP_JS"
    echo "TIED_BASE_PATH=$TIED_BASE_PATH"
    echo "TIED_METHODOLOGY_BUNDLE_PATH=$TIED_METHODOLOGY_BUNDLE_PATH"
    echo "TIED_STORE_ROOT=$TIED_STORE_ROOT"
    echo "TIED_CLI=$TIED_CLI"
  } >> "$GITHUB_ENV"
fi
