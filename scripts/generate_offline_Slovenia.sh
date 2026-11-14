#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

BASE_SOURCE="${MAPD_BASE_SOURCE:-poland-germany-czech-austria-italy-latest.osm.pbf}"
SLOVENIA_SOURCE="${MAPD_SLOVENIA_SOURCE:-slovenia-latest.osm.pbf}"
OUTPUT_DIR="${MAPD_OUTPUT_DIR:-offline-slovenia}"
WORK_DIR="${MAPD_WORK_DIR:-slovenia-work}"

# Quarter-degree tile boundaries enclosing the full Geofabrik Slovenia extract.
MIN_LON=13.25
MIN_LAT=45.25
MAX_LON=16.75
MAX_LAT=47.0

# Extra source context preserves roads crossing the generated tile boundaries.
EXTRACT_BBOX="12.5,44.25,17.75,48.0"
HIGHWAY_FILTER="nw/highway=motorway,trunk,primary,secondary,tertiary,unclassified,residential,motorway_link,trunk_link,primary_link,secondary_link,tertiary_link"

for source in "$BASE_SOURCE" "$SLOVENIA_SOURCE"; do
  if [[ ! -f "$source" ]]; then
    echo "missing source: $source" >&2
    exit 1
  fi
done

if ! command -v osmium >/dev/null 2>&1; then
  echo "missing required tool: osmium" >&2
  exit 1
fi

source "$SCRIPT_DIR/require_deletable_dir.sh"
require_deletable_dir MAPD_OUTPUT_DIR "$OUTPUT_DIR"
require_deletable_dir MAPD_WORK_DIR "$WORK_DIR"

./build_mapd.sh

rm -rf -- "$WORK_DIR" "$OUTPUT_DIR"
mkdir -p "$WORK_DIR" "$OUTPUT_DIR"

echo "extracting existing five-country data around Slovenia"
osmium extract --bbox "$EXTRACT_BBOX" "$BASE_SOURCE" -o "$WORK_DIR/base-box.osm.pbf" --overwrite

echo "extracting Slovenia source data"
osmium extract --bbox "$EXTRACT_BBOX" "$SLOVENIA_SOURCE" -o "$WORK_DIR/slovenia-box.osm.pbf" --overwrite

# Extracts of different dates can hold different versions of one OSM object.
# Keep only the newest.
echo "merging existing data with Slovenia"
osmium merge "$WORK_DIR/base-box.osm.pbf" "$WORK_DIR/slovenia-box.osm.pbf" -f pbf -o - |
  osmium time-filter -F pbf - -o "$WORK_DIR/combined-box.osm.pbf" --overwrite

echo "filtering routable roads"
osmium tags-filter "$WORK_DIR/combined-box.osm.pbf" "$HIGHWAY_FILTER" -o "$WORK_DIR/filtered.osm.pbf" --overwrite
osmium add-locations-to-ways "$WORK_DIR/filtered.osm.pbf" -o "$WORK_DIR/map.osm.pbf" --overwrite

echo "generating Slovenia tiles with custom map overrides"
./mapd generate \
  --input-file "$WORK_DIR/map.osm.pbf" \
  --output-directory "$OUTPUT_DIR" \
  --minlat "$MIN_LAT" \
  --minlon "$MIN_LON" \
  --maxlat "$MAX_LAT" \
  --maxlon "$MAX_LON"

tile_count="$(find "$OUTPUT_DIR" -type f | wc -l | tr -d ' ')"
tile_bytes="$(find "$OUTPUT_DIR" -type f -exec stat -f '%z' {} \; | awk '{total += $1} END {print total + 0}')"
echo "generated $tile_count Slovenia-area tiles ($tile_bytes bytes) in $SCRIPT_DIR/$OUTPUT_DIR"

if [[ "${MAPD_UPLOAD:-0}" == "1" ]]; then
  MAPD_OFFLINE_SOURCE="$OUTPUT_DIR" ./upload_selective_offline_comma.sh
else
  echo "not uploaded; set MAPD_UPLOAD=1 to add these tiles to the device"
fi
