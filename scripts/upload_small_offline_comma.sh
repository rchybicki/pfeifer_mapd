#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOST="${MAPD_SSH_HOST:-commawifi}"
REMOTE_PATH="${MAPD_OFFLINE_REMOTE_PATH:-/data/media/0/osm/offline/}"

ssh "$HOST" "mkdir -p \"$REMOTE_PATH\""
rsync -av --delete --exclude='*.tar.gz' --exclude='._*' "$SCRIPT_DIR/offline/" "$HOST:$REMOTE_PATH"
echo "offline maps uploaded to $HOST:$REMOTE_PATH"
