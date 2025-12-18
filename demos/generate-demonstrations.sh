#!/bin/bash
# Comprehensive script to generate all visual demonstrations for bkpdir
# Generates: screen recordings for demos, diagrams, and before/after screenshots
# Usage: ./generate-demonstrations.sh [1|2]
#   - No argument: Generate both Demo 1 and Demo 2
#   - Argument "1": Generate only Demo 1
#   - Argument "2": Generate only Demo 2

set -e

# Parse optional argument
DEMO_NUM="${1:-}"
GENERATE_DEMO_config=true
GENERATE_DEMO_full=true
GENERATE_DEMO_help=true

if [ -n "${DEMO_NUM}" ]; then
    if [ "${DEMO_NUM}" = "help" ]; then
        GENERATE_DEMO_config=false
        GENERATE_DEMO_full=false
        GENERATE_DEMO_help=true
    elif [ "${DEMO_NUM}" = "full" ]; then
        GENERATE_DEMO_config=false
        GENERATE_DEMO_full=true
        GENERATE_DEMO_help=false
    elif [ "${DEMO_NUM}" = "config" ]; then
        GENERATE_DEMO_config=true
        GENERATE_DEMO_full=false
        GENERATE_DEMO_help=false
    else
        echo "Error: Invalid argument '${DEMO_NUM}'"
        echo "Usage: $0 [config|full|help]"
        echo "  No argument: Generate all"
        exit 1
    fi
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
DEMO_DIR_config="${HOME}/.bkpdir-demo-config"
DEMO_DIR_full="${HOME}/.bkpdir-demo-full"
DEMO_DIR_help="${HOME}/.bkpdir-demo-help"
IMAGES_DIR="${PROJECT_ROOT}/images"
# BKPDIR_BIN="${HOME}/.local/bin/bkpdir"
BKPDIR_BIN="bkpdir"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== BkpDir Visual Demonstrations Generator ===${NC}"
if [ -n "${DEMO_NUM}" ]; then
    echo -e "${CYAN}Generating Demo ${DEMO_NUM} only${NC}"
else
    echo -e "${CYAN}Generating both Demo 1 and Demo 2${NC}"
fi
echo ""

# Check prerequisites
echo -e "${YELLOW}Checking prerequisites...${NC}"
if ! command -v bkpdir &> /dev/null && [ ! -f "${BKPDIR_BIN}" ]; then
    echo "Error: bkpdir not found. Please install or set BKPDIR_BIN"
    exit 1
fi

if [ -f "${BKPDIR_BIN}" ]; then
    BKPDIR_BIN="${BKPDIR_BIN}"
elif command -v bkpdir &> /dev/null; then
    BKPDIR_BIN="bkpdir"
fi

if ! command -v asciinema &> /dev/null; then
    echo "Warning: asciinema not found. Screen recording will be skipped."
    echo "Install with: brew install asciinema or pip install asciinema"
    SKIP_RECORDING=true
else
    SKIP_RECORDING=false
fi

# Create directories
mkdir -p "${DEMO_DIR_config}"
mkdir -p "${DEMO_DIR_full}"
mkdir -p "${DEMO_DIR_help}"
mkdir -p "${IMAGES_DIR}"

cd "${PROJECT_ROOT}"

# Function to convert cast to GIF
convert_to_gif() {
    local cast_file="$1"
    local gif_file="$2"
    local demo_name="$3"
    
    if command -v agg &> /dev/null; then
        echo -e "${CYAN}Converting ${demo_name} to GIF...${NC}"
        ARGS="--theme asciinema"
        LFD="--last-frame-duration 1"
        if agg $ARGS $LFD "${cast_file}" "${gif_file}"; then
            echo -e "${GREEN}✓ GIF saved to ${gif_file}${NC}"
            
            # Optimize the GIF using gifsicle if available
            if command -v gifsicle &> /dev/null; then
                local optimized_file="${gif_file%.gif}-optimized.gif"
                gifsicle -O3 --colors 256 "${gif_file}" -o "${optimized_file}" && \
                    echo -e "${GREEN}✓ Optimized GIF saved to ${optimized_file}${NC}"
            fi
        else
            echo "Warning: GIF conversion failed. Install agg: pip install agg"
        fi
    else
        echo "To convert to GIF: pip install agg && agg ${cast_file} ${gif_file}"
    fi
}

gen_call () {
    local ID="$1"
    local varname="GENERATE_DEMO_$ID"
    local generate_demo="${!varname}"

    export DEMO_ID="$ID"

    # Step 1: Generate Demo 1 screen recording
    if [ "$generate_demo" = true ] && [ "$SKIP_RECORDING" = false ]; then
        echo -e "${GREEN}Step 1: Generating Demo 1 screen recording...${NC}"
        true
    elif [ "$generate_demo" = false ]; then
        echo -e "${YELLOW}Skipping Demo 1 (not requested)${NC}"
        false
    elif [ "$SKIP_RECORDING" = true ]; then
        echo -e "${YELLOW}Skipping screen recording (asciinema not available)${NC}"
        false
    fi
}

if gen_call help; then
    demo_dir_varname="DEMO_DIR_${DEMO_ID}"
    cd "${!demo_dir_varname}"
    
    fn_cast="bkpdir-demo-$DEMO_ID.cast"
    asciinema rec --overwrite $fn_cast -c 'echo "\n\033[1;32mCommand help and usage\033[0m\n"'
    asciinema rec --append $fn_cast -c 'stdbuf -oL -eL bkpdir --help | pv -q -L 4 --line-mode'

    mv $fn_cast "${IMAGES_DIR}/$fn_cast"
    echo -e "${GREEN}✓ Demo $DEMO_ID recording saved to ${IMAGES_DIR}/$fn_cast${NC}"

    convert_to_gif "${IMAGES_DIR}/$fn_cast" "${IMAGES_DIR}/bkpdir-demo-$DEMO_ID.gif" "Demo $DEMO_ID"

    cd "${PROJECT_ROOT}"
fi

if gen_call config; then
    demo_dir_varname="DEMO_DIR_${DEMO_ID}"
    cd "${!demo_dir_varname}"
    
    fn_cast="bkpdir-demo-$DEMO_ID.cast"
    asciinema rec --overwrite $fn_cast -c 'echo "\n\033[1;32mCurrent configuration\033[0m\n"'
    # Append the output of the command to the file
    asciinema rec --append $fn_cast -c 'echo "\n\033[2;32mCurrent configuration\033[0m\n"'
    asciinema rec --append $fn_cast -c 'stdbuf -oL -eL bkpdir config --format tree | pv -q -L 4 --line-mode'

    mv $fn_cast "${IMAGES_DIR}/$fn_cast"
    echo -e "${GREEN}✓ Demo $DEMO_ID recording saved to ${IMAGES_DIR}/$fn_cast${NC}"

    convert_to_gif "${IMAGES_DIR}/$fn_cast" "${IMAGES_DIR}/bkpdir-demo-$DEMO_ID.gif" "Demo $DEMO_ID"

    cd "${PROJECT_ROOT}"
fi

if gen_call full; then
    demo_dir_varname="DEMO_DIR_${DEMO_ID}"
    cd "${!demo_dir_varname}"
    
    fn_cast="bkpdir-demo-$DEMO_ID.cast"
    # Use expect script for more reliable automation
    if [ -f "${SCRIPT_DIR}/demo-automation-$DEMO_ID.exp" ]; then
        [ -f "$fn_cast" ] && rm $fn_cast
        # asciinema rec --append bkpdir-demo-2.cast -c 'echo "\033[1;36mThis demo will show advanced bkpdir features: comparison, restore, and configuration management\033[0m\n"'
        # asciinema rec --append bkpdir-demo-2.cast -c 'echo "\033[1;36mThis demo will show advanced bkpdir features: comparison, restore, and configuration management\033[0m\n"; sleep 2; expect -c send -- "bash -lc '\''cd \"$DEMO_DIR_2\" && mkdir -p myproject && cd myproject && git init'\''\r"'

        "${SCRIPT_DIR}/demo-automation-$DEMO_ID.exp" "${!demo_dir_varname}" "${BKPDIR_BIN}" "${PROJECT_ROOT}" "${fn_cast}" || {
            echo "Warning: Demo $DEMO_ID recording may have issues. Check $fn_cast"
        }
        
        set -x
        if [ -f "$fn_cast" ]; then
            # append summary
            asciinema rec --append $fn_cast -c 'echo "\n\033[1;32mSummary: Advanced features demonstrated\033[0m\n"; sleep 1; echo "\033[1;33m✓ Advanced configuration with exclusions\033[0m"; sleep 1; echo "\033[1;33m✓ Multiple incremental archives\033[0m"; sleep 1; echo "\033[1;33m✓ Archive listing and inspection\033[0m"; sleep 1; echo "\033[1;33m✓ Configuration management\033[0m"; sleep 2'

            mv $fn_cast "${IMAGES_DIR}/$fn_cast"
            echo -e "${GREEN}✓ Demo $DEMO_ID recording saved to ${IMAGES_DIR}/$fn_cast${NC}"

            convert_to_gif "${IMAGES_DIR}/$fn_cast" "${IMAGES_DIR}/bkpdir-demo-$DEMO_ID.gif" "Demo $DEMO_ID"
            echo "- ${IMAGES_DIR}/bkpdir-demo-$DEMO_ID.gif" 
            open "${IMAGES_DIR}/bkpdir-demo-$DEMO_ID.gif" 
        fi
    else
        echo "Warning: demo-automation-$DEMO_ID.exp not found"
    fi
    cd "${PROJECT_ROOT}"
fi

if false; then
    # Step 3: Create demo environment for screenshots (always generated)
    STEP_NUM=3
    echo -e "${GREEN}Step ${STEP_NUM}: Creating demo environment for screenshots...${NC}"
    cd "${DEMO_DIR}"
    rm -rf myapp
    mkdir -p myapp
    cd myapp

    # Initialize git
    git init -q
    git config user.name "Demo User"
    git config user.email "demo@example.com"

    # Create sample files
    echo '# My Application' > README.md
    echo 'package main' > main.go
    echo 'console.log("Hello")' > app.js
    mkdir -p src
    echo '// source file' > src/utils.js

    # Initial commit
    git add . > /dev/null 2>&1
    git commit -m 'Initial commit' > /dev/null 2>&1

    # Create .bkpdir.yml
    cat > .bkpdir.yml << 'EOF'
    # .bkpdir.yml
    archive_dir_path: "./archives"
    use_current_dir_name: true

    # Exclude patterns for files/directories to skip
    exclude_patterns:
      - ".git/"
      - "node_modules/"
      - "*.tmp"
      - "*.log"

    # Git integration settings
    include_git_info: true
    show_git_dirty_status: true
    include_branch: true
    EOF

    STEP_NUM=$((STEP_NUM + 1))
    echo -e "${GREEN}Step ${STEP_NUM}: Generating before/after screenshots...${NC}"

    # Before screenshot - directory structure
    echo "=== Before: Project Structure ===" > "${IMAGES_DIR}/before-structure.txt"
    find . -not -path '*/\.*' -not -path '*/archives/*' | sort >> "${IMAGES_DIR}/before-structure.txt"

    # Create full archive
    echo -e "${CYAN}Creating full archive...${NC}"
    "${BKPDIR_BIN}" . > /dev/null 2>&1 || true

    # Make changes
    echo '// Updated code' >> main.go
    echo '' >> README.md
    echo 'Updated README' >> README.md
    git add . > /dev/null 2>&1
    git commit -m 'Add new features' > /dev/null 2>&1

    # Create incremental archive
    "${BKPDIR_BIN}" inc > /dev/null 2>&1 || true

    # After screenshot - archive listing
    echo "=== After: Archive Listing ===" > "${IMAGES_DIR}/after-archives.txt"
    "${BKPDIR_BIN}" list >> "${IMAGES_DIR}/after-archives.txt" 2>&1 || true

    # Archive directory structure
    echo "" >> "${IMAGES_DIR}/after-archives.txt"
    echo "=== Archive Directory Structure ===" >> "${IMAGES_DIR}/after-archives.txt"
    find archives -name '*.zip' 2>/dev/null | sort >> "${IMAGES_DIR}/after-archives.txt" || true

    echo -e "${GREEN}✓ Screenshots saved to ${IMAGES_DIR}/before-structure.txt and ${IMAGES_DIR}/after-archives.txt${NC}"

    STEP_NUM=$((STEP_NUM + 1))
    echo -e "${GREEN}Step ${STEP_NUM}: Generating timeline data...${NC}"
    cd "${PROJECT_ROOT}"

    # Create a simple timeline text representation
    cat > "${IMAGES_DIR}/archive-timeline.txt" << 'EOF'
    BkpDir Archive Timeline
    =======================

    Day 1: Initial Full Archive
      ├─ Archive: myapp-YYYY-MM-DD-HH-MM=main=abc123.zip
      ├─ Size: ~50 KB (full backup)
      └─ Contains: All project files

    Day 2: Incremental Archive
      ├─ Archive: myapp_update=YYYY-MM-DD-HH-MM=main=def456.zip
      ├─ Size: ~5 KB (only changes)
      └─ Contains: Modified files only

    Day 3: Incremental Archive
      ├─ Archive: myapp_update=YYYY-MM-DD-HH-MM=main=ghi789.zip
      ├─ Size: ~3 KB (only changes)
      └─ Contains: Modified files only

    Disk Space Usage Over Time:
      Day 1: 50 KB
      Day 2: 55 KB (50 + 5)
      Day 3: 58 KB (50 + 5 + 3)

    Archive Relationships:
      Full Archive (Day 1)
        └─> Incremental (Day 2) [depends on Day 1]
        └─> Incremental (Day 3) [depends on Day 1]
    EOF

    echo -e "${GREEN}✓ Timeline data saved to ${IMAGES_DIR}/archive-timeline.txt${NC}"

    STEP_NUM=$((STEP_NUM + 1))
    echo -e "${GREEN}Step ${STEP_NUM}: Generating comparison output...${NC}"
    cd "${DEMO_DIR}/myapp"

    # Create comparison data
    cat > "${IMAGES_DIR}/archive-comparison.txt" << 'EOF'
    Archive Comparison
    =================

    Full Archive vs Incremental Archive

    Full Archive (Day 1):
      - Contains: All files in project
      - Size: Larger (~50 KB)
      - Use case: Initial backup, complete restore
      - Naming: myapp-YYYY-MM-DD-HH-MM=BRANCH=HASH.zip

    Incremental Archive (Day 2+):
      - Contains: Only changed files
      - Size: Smaller (~3-5 KB)
      - Use case: Regular backups, space-efficient
      - Naming: myapp_update=YYYY-MM-DD-HH-MM=BRANCH=HASH.zip
      - Dependency: Requires base full archive

    Benefits:
      ✓ Incremental archives save disk space
      ✓ Faster to create (only changed files)
      ✓ Maintains full history with base archive
      ✓ Git integration tracks branch and commit
    EOF

    echo -e "${GREEN}✓ Comparison data saved to ${IMAGES_DIR}/archive-comparison.txt${NC}"

    cd "${PROJECT_ROOT}"

    echo ""
    echo -e "${BLUE}=== Generation Complete ===${NC}"
    echo ""
    echo -e "${CYAN}Generated files:${NC}"
    # [ -f "${IMAGES_DIR}/bkpdir-demo-1.cast" ] && echo "  📹 Demo 1 recording: ${IMAGES_DIR}/bkpdir-demo-1.cast"
    # [ -f "${IMAGES_DIR}/bkpdir-demo-1.gif" ] && echo "  🎬 Demo 1 GIF: ${IMAGES_DIR}/bkpdir-demo-1.gif"
    # [ -f "${IMAGES_DIR}/bkpdir-demo-2.cast" ] && echo "  📹 Demo 2 recording: ${IMAGES_DIR}/bkpdir-demo-2.cast"
    # [ -f "${IMAGES_DIR}/bkpdir-demo-2.gif" ] && echo "  🎬 Demo 2 GIF: ${IMAGES_DIR}/bkpdir-demo-2.gif"
    echo "  📊 Timeline: ${IMAGES_DIR}/archive-timeline.txt"
    echo "  📋 Comparison: ${IMAGES_DIR}/archive-comparison.txt"
    echo "  📸 Before: ${IMAGES_DIR}/before-structure.txt"
    echo "  📸 After: ${IMAGES_DIR}/after-archives.txt"
    echo ""
    echo -e "${CYAN}Next steps:${NC}"
    echo "  1. Review generated files in ${IMAGES_DIR}/"
    echo "  2. Convert .cast to GIF: pip install agg && agg <cast-file> <gif-file>"
    echo "  3. Create visual diagrams from text files using your preferred tool"
    echo "  4. Update documentation with generated demonstrations"
    echo ""
fi
