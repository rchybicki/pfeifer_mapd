#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SOURCE="${MAPD_OFFLINE_SOURCE:-$SCRIPT_DIR/offline-slovenia}"
HOST="${MAPD_SSH_HOST:-comma}"
REMOTE_PATH="${MAPD_OFFLINE_REMOTE_PATH:-/data/media/0/osm-v240/offline/}"

if [[ ! -d "$SOURCE" ]]; then
  echo "missing selective offline map directory: $SOURCE" >&2
  exit 1
fi

tile_count="$(find "$SOURCE" -type f | wc -l | tr -d ' ')"
if [[ "$tile_count" == "0" ]]; then
  echo "refusing to upload an empty selective map directory: $SOURCE" >&2
  exit 1
fi

echo "adding $tile_count tiles to $HOST:$REMOTE_PATH (existing files outside this set are preserved)"
ssh "$HOST" "mkdir -p \"$REMOTE_PATH\""
rsync -avcz --itemize-changes --stats --exclude='*.tar.gz' --exclude='._*' "$SOURCE/" "$HOST:$REMOTE_PATH"
echo "selective offline maps uploaded to $HOST:$REMOTE_PATH"
echo "runtime mapd binary was not uploaded"
