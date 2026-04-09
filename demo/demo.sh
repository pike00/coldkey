#!/bin/bash
# =============================================================================
# coldkey demo recording script
#
# Designed for asciinema recording. Simulates realistic typing with pauses
# between sections so viewers can read the output before the next command.
#
# Usage:
#   asciinema rec --cols 110 --rows 38 -c ./demo/demo.sh demo.cast
#
# Prerequisites:
#   - ./coldkey binary built (just build)
#   - asciinema installed
# =============================================================================

set -euo pipefail

DEMO_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$DEMO_DIR")"
BINARY="$PROJECT_ROOT/coldkey"
OUTPUT_DIR="$PROJECT_ROOT/output"

# -- Tuning knobs --
CHAR_DELAY=0.055        # seconds between each typed character
POST_TYPE_DELAY=1     # pause after typing before execution
SECTION_PAUSE=2.5       # pause between demo sections
READ_PAUSE=3          # pause to let viewer read short output
LONG_READ_PAUSE=6     # pause for longer output blocks
INITIAL_PAUSE=1.0       # pause before the demo begins

# -- Helpers --

# Simulate typing a command character-by-character
type_cmd() {
    local cmd="$1"
    for (( i=0; i<${#cmd}; i++ )); do
        printf '%s' "${cmd:$i:1}"
        sleep "$CHAR_DELAY"
    done
    echo
    sleep "$POST_TYPE_DELAY"
}

# Print a colored shell prompt
prompt() {
    printf '\033[1;32m$\033[0m '
}

# Print a section comment (dimmed, like a # comment in the terminal)
section() {
    echo
    printf '\033[0;90m# %s\033[0m\n' "$1"
    sleep 0.6
}

# -- Setup --
mkdir -p "$OUTPUT_DIR"
rm -f "$OUTPUT_DIR"/demo-key* 2>/dev/null

# =============================================================================
# Demo begins
# =============================================================================

sleep "$INITIAL_PAUSE"

# --- Section 1: Show what coldkey does ---
section "See what coldkey can do"
prompt
type_cmd "coldkey --help"
"$BINARY" --help
sleep "$LONG_READ_PAUSE"

# --- Section 2: Generate a post-quantum age key ---
section "Generate a new post-quantum age key with paper backup"
prompt
type_cmd "coldkey generate -o output/demo-key.txt"
"$BINARY" generate -o "$OUTPUT_DIR/demo-key.txt"
sleep "$SECTION_PAUSE"

# --- Section 3: Show the generated files ---
section "Check what was created"
prompt
type_cmd "ls -lh output/"
ls -lh "$OUTPUT_DIR/"
sleep "$READ_PAUSE"

# --- Section 4: Peek at the key file ---
section "Peek at the key file (public recipient + private key)"
prompt
type_cmd "head -5 output/demo-key.txt"
head -5 "$OUTPUT_DIR/demo-key.txt"
sleep "$LONG_READ_PAUSE"

# --- Section 5: Create a standalone backup from an existing key ---
section "Create a paper backup from any existing age key"
prompt
type_cmd "coldkey backup -o output/demo-backup.html output/demo-key.txt"
"$BINARY" backup -o "$OUTPUT_DIR/demo-backup.html" "$OUTPUT_DIR/demo-key.txt"
sleep "$SECTION_PAUSE"

# --- Section 6: Show the backup file ---
section "Paper backup ready to print"
prompt
type_cmd "ls -lh output/demo-backup.html"
ls -lh "$OUTPUT_DIR/demo-backup.html"
sleep "$READ_PAUSE"

# --- Section 7: Version ---
section "Check the version"
prompt
type_cmd "coldkey version"
"$BINARY" version
sleep "$SECTION_PAUSE"

# --- End ---
echo
printf '\033[1;36m# Done! Print the HTML backup and store it somewhere safe.\033[0m\n'
sleep 2.0
