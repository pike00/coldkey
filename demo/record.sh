#!/bin/bash
# =============================================================================
# Record the coldkey demo with asciinema
#
# Builds the binary, cleans previous output, and records a new demo.cast
# in the project root.
#
# Usage:
#   ./demo/record.sh
# =============================================================================

set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_ROOT"

# Build binary with a clean version tag (strip -dirty suffix for the recording)
VERSION=$(git describe --tags --always 2>/dev/null | sed 's/-dirty$//' || echo dev)
echo "Building coldkey ($VERSION)..."
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=$VERSION" -o coldkey ./cmd/coldkey

# Clean previous demo artifacts
rm -f output/demo-key* output/demo-backup* demo.cast

# Record (ensure TERM is set so asciinema doesn't write null into the cast header)
export TERM="${TERM:-xterm-256color}"
echo "Starting asciinema recording..."
asciinema rec --cols 110 --rows 38 --overwrite -c ./demo/demo.sh demo.cast

echo "Recording saved to demo.cast"
echo "Preview with: asciinema play demo.cast"
