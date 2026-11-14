#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SOURCE="${MAPD_OFFLINE_SOURCE:-$SCRIPT_DIR/offline}"
HOST="${MAPD_SSH_HOST:-commawifi}"
REMOTE_PATH="${MAPD_OFFLINE_REMOTE_PATH:-/data/media/0/osm/offline/}"

if [[ ! -d "$SOURCE" ]]; then
  echo "missing offline map directory: $SOURCE" >&2
  exit 1
fi

# rsync --delete from a source with no tile to send would remove every tile on the device.
# Tiles are <lat>/<lon>/<minlat>_<minlon>_<maxlat>_<maxlon>; archives and dot entries are not sent.
tile_count="$(find "$SOURCE/" -mindepth 1 \( -name '*.tar.gz' -o -name '.*' \) -prune -o -type f -name '*_*_*_*' -print | wc -l | tr -d ' ')"
if [[ "$tile_count" == "0" ]]; then
  echo "refusing to upload an empty map directory: $SOURCE" >&2
  exit 1
fi

ssh "$HOST" "mkdir -p \"$REMOTE_PATH\""
rsync -av --delete --exclude='*.tar.gz' --exclude='._*' "$SOURCE/" "$HOST:$REMOTE_PATH"
echo "offline maps uploaded to $HOST:$REMOTE_PATH"
echo "runtime mapd binary is managed in ../openpilot-rch/frogpilot/navigation/mapd"
