# BkpDir Demo Scripts

This directory contains scripts and documentation for creating demonstration GIFs and screen recordings of bkpdir functionality.

## Demo Scripts

### Demo 1: Basic Workflow (`demo-automation-1.exp`)

Demonstrates the complete basic workflow:
- Project setup and git initialization
- Configuration file creation
- Full archive creation
- Incremental archive creation
- Archive listing and inspection
- Help and configuration display

**Output**: `bkpdir-demo-1.cast` → `bkpdir-demo-1.gif`

### Demo 2: Advanced Features (`demo-automation-2.exp`)

Demonstrates advanced features and scenarios:
- [To be defined - will show different aspects]

**Output**: `bkpdir-demo-2.cast` → `bkpdir-demo-2.gif`

## Quick Start

### Prerequisites

```bash
# Install asciinema
brew install asciinema
# or
pip install asciinema

# Install agg for GIF conversion
pip install agg

# Install gifsicle for GIF optimization (optional)
brew install gifsicle

# Verify bkpdir is in PATH or set BKPDIR_BIN
export BKPDIR_BIN="${HOME}/.local/bin/bkpdir"
```

### Running Demo 1

```bash
cd demos
./demo-automation-1.exp \
  "${HOME}/.bkpdir-demo" \
  "${HOME}/.local/bin/bkpdir" \
  "$(cd .. && pwd)"
```

### Running Demo 2

```bash
cd demos
./demo-automation-2.exp \
  "${HOME}/.bkpdir-demo-2" \
  "${HOME}/.local/bin/bkpdir" \
  "$(cd .. && pwd)"
```

### Converting to GIF

After recording, convert the `.cast` file to GIF:

```bash
# Convert Demo 1
agg --theme asciinema --last-frame-duration 1 \
  bkpdir-demo-1.cast ../images/bkpdir-demo-1.gif

# Convert Demo 2
agg --theme asciinema --last-frame-duration 1 \
  bkpdir-demo-2.cast ../images/bkpdir-demo-2.gif

# Optimize GIFs (optional)
gifsicle -O3 --colors 256 ../images/bkpdir-demo-1.gif \
  -o ../images/bkpdir-demo-1-optimized.gif
gifsicle -O3 --colors 256 ../images/bkpdir-demo-2.gif \
  -o ../images/bkpdir-demo-2-optimized.gif
```

### Using the Generation Script

The `generate-demonstrations.sh` script can generate both demos:

```bash
./generate-demonstrations.sh
```

## Files

- `demo-automation-1.exp` - Expect script for Demo 1 (basic workflow)
- `demo-automation-1.applescript` - AppleScript alternative for Demo 1
- `demo-automation-2.exp` - Expect script for Demo 2 (advanced features)
- `record-demo-1.sh` - Shell script wrapper for Demo 1
- `generate-demonstrations.sh` - Comprehensive script to generate all demos
- `README.md` - This file

## Output Locations

- Recordings: `../images/bkpdir-demo-*.cast`
- GIFs: `../images/bkpdir-demo-*.gif`
- Optimized GIFs: `../images/bkpdir-demo-*-optimized.gif`

## Troubleshooting

### Terminal Not Responding
- Use expect scripts instead of AppleScript
- Increase delay values in script

### bkpdir Not Found
- Set `BKPDIR_BIN` environment variable
- Ensure bkpdir is in PATH
- Use full path in script

### Asciinema Issues
- Verify installation: `which asciinema`
- Check terminal permissions
- Ensure terminal is in focus

## See Also

- [Visual Demonstration Suggestions](../docs/visual-demonstration-suggestions.md) - Overview
- [Screen Recording Plan](../docs/screen-recording-plan.md) - Complete documentation
