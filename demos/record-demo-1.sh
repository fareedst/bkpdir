#!/bin/bash
# Script to automate bkpdir demonstration recording - Demo 1
# Uses osascript to control Terminal.app and asciinema to record

set -e

# Configuration
DEMO_DIR="${HOME}/.bkpdir-demo"
BKPDIR_BIN="${HOME}/.local/bin/bkpdir"  # Adjust path as needed
TERMINAL_APP="Terminal"
RECORDING_FILE="bkpdir-demo-1.cast"
RECORDING_OUTPUT="../images/bkpdir-demo-1.gif"  # Optional: convert to GIF

# Check prerequisites
if ! command -v asciinema &> /dev/null; then
    echo "Error: asciinema is not installed"
    echo "Install with: brew install asciinema or pip install asciinema"
    exit 1
fi

if ! command -v osascript &> /dev/null; then
    echo "Error: osascript is not available (macOS required)"
    exit 1
fi

# Check if bkpdir exists
if [ ! -f "${BKPDIR_BIN}" ] && ! command -v bkpdir &> /dev/null; then
    echo "Error: bkpdir binary not found at ${BKPDIR_BIN} and not in PATH"
    echo "Please set BKPDIR_BIN environment variable or install bkpdir"
    exit 1
fi

# Use bkpdir from PATH if binary not found at specified path
if [ ! -f "${BKPDIR_BIN}" ]; then
    BKPDIR_BIN="bkpdir"
fi

# Cleanup function
cleanup() {
    echo "Cleaning up..."
    rm -rf "${DEMO_DIR}"
    # Note: Terminal window cleanup handled by osascript
}

trap cleanup EXIT

# Create demo directory
echo "Setting up demo environment..."
mkdir -p "${DEMO_DIR}"
cd "${DEMO_DIR}"

# Create images directory if it doesn't exist
mkdir -p "$(dirname "${RECORDING_OUTPUT}")"

# Get script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Run automation script
echo "Starting recording automation for Demo 1..."
osascript "${SCRIPT_DIR}/demo-automation-1.applescript" "${DEMO_DIR}" "${BKPDIR_BIN}" "${PROJECT_ROOT}"

echo ""
echo "Recording complete!"
echo "Recording file: ${RECORDING_FILE}"
echo ""
echo "To convert to GIF (requires agg):"
echo "  pip install agg"
echo "  agg ${RECORDING_FILE} ${RECORDING_OUTPUT}"
echo ""
echo "To upload to asciinema.org:"
echo "  asciinema upload ${RECORDING_FILE}"
