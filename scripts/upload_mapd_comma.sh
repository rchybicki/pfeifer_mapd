#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOST="${MAPD_SSH_HOST:-commawifi}"
REMOTE_PATH="${MAPD_REMOTE_PATH:-/data/openpilot/frogpilot/navigation/mapd}"

ssh "$HOST" "mkdir -p \"$(dirname "$REMOTE_PATH")\""
rsync "$SCRIPT_DIR/mapd-comma" "$HOST:$REMOTE_PATH"
ssh "$HOST" "chmod +x \"$REMOTE_PATH\""
echo "mapd-comma uploaded to $HOST:$REMOTE_PATH"
