#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

SOURCE="${MAPD_POLAND_SOURCE:-poland-latest.osm.pbf}"
BACKUP="${SOURCE}.pre-update"
UPDATE_SIZE_MB="${MAPD_UPDATE_SIZE_MB:-1000}"

for command_name in osmium pyosmium-up-to-date; do
  if ! command -v "$command_name" >/dev/null 2>&1; then
    echo "missing required tool: $command_name" >&2
    exit 1
  fi
done

if [[ ! -f "$SOURCE" ]]; then
  echo "missing Poland source: $SOURCE" >&2
  exit 1
fi

if [[ -e "$BACKUP" ]]; then
  echo "refusing to overwrite existing update backup: $BACKUP" >&2
  exit 1
fi

if ! cp -c "$SOURCE" "$BACKUP" 2>/dev/null; then
  cp "$SOURCE" "$BACKUP"
fi

while true; do
  set +e
  pyosmium-up-to-date -v --size "$UPDATE_SIZE_MB" "$SOURCE"
  update_status=$?
  set -e

  case "$update_status" in
    0)
      break
      ;;
    1)
      echo "more replication data is available; continuing"
      ;;
    *)
      echo "Poland update failed; original source remains in $BACKUP" >&2
      exit "$update_status"
      ;;
  esac
done

osmium check-refs "$SOURCE"
replication_timestamp="$(osmium fileinfo -e -g header.option.osmosis_replication_timestamp "$SOURCE")"
replication_sequence="$(osmium fileinfo -e -g header.option.osmosis_replication_sequence_number "$SOURCE")"
echo "Poland source updated to sequence $replication_sequence ($replication_timestamp)"

unlink "$BACKUP"
