#!/usr/bin/env bash
set -euo pipefail

# Check for import cycles or build errors by invoking `go list ./...`
if ! go list ./... >/dev/null 2>&1; then
  echo "ERROR: go list failed — possible import cycle or build error"
  go list ./... 2>&1 | sed -n '1,200p'
  exit 1
fi

# Ensure diagram exists
if [ ! -f docs/images/package-interdependency-mapping.svg ]; then
  echo "ERROR: Missing diagram: docs/images/package-interdependency-mapping.svg"
  exit 2
fi

echo "check_import_cycles: OK"
