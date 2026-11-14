#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -- "$SCRIPT_DIR/.." && pwd)"
LOCAL_OUTPUT="${MAPD_LOCAL_BINARY:-$SCRIPT_DIR/mapd}"

echo "building local mapd binary -> $LOCAL_OUTPUT"
go build -o "$LOCAL_OUTPUT" "$REPO_ROOT"
