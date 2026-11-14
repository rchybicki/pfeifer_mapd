#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

POLAND_SOURCE="${MAPD_POLAND_SOURCE:-poland-latest.osm.pbf}"
GERMANY_SOURCE="${MAPD_GERMANY_SOURCE:-germany-latest.osm.pbf}"
CZECH_SOURCE="${MAPD_CZECH_SOURCE:-czech-republic-latest.osm.pbf}"
AUSTRIA_SOURCE="${MAPD_AUSTRIA_SOURCE:-austria-latest.osm.pbf}"
OUTPUT_DIR="${MAPD_OUTPUT_DIR:-offline-poland-update}"
WORK_DIR="${MAPD_WORK_DIR:-poland-update-work}"

# Quarter-degree tile boundaries enclosing the Poland download-menu bounds.
MIN_LON=14.0
MIN_LAT=49.0
MAX_LON=24.25
MAX_LAT=55.0

# Extra context retains roads crossing Poland's borders.
EXTRACT_BBOX="13.0,48.0,25.25,56.0"
HIGHWAY_FILTER="nw/highway=motorway,trunk,primary,secondary,tertiary,unclassified,residential,motorway_link,trunk_link,primary_link,secondary_link,tertiary_link"

SOURCES=("$POLAND_SOURCE" "$GERMANY_SOURCE" "$CZECH_SOURCE" "$AUSTRIA_SOURCE")
SOURCE_NAMES=(poland germany czech-republic austria)

for source_index in "${!SOURCES[@]}"; do
  if [[ ! -f "${SOURCES[$source_index]}" ]]; then
    echo "missing ${SOURCE_NAMES[$source_index]} source: ${SOURCES[$source_index]}" >&2
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

EXTRACTS=()
for source_index in "${!SOURCES[@]}"; do
  extract="$WORK_DIR/${SOURCE_NAMES[$source_index]}-box.osm.pbf"
  echo "extracting ${SOURCE_NAMES[$source_index]} border context"
  osmium extract --bbox "$EXTRACT_BBOX" "${SOURCES[$source_index]}" -o "$extract" --overwrite
  EXTRACTS+=("$extract")
done

echo "merging updated Poland with unchanged neighboring sources"
osmium merge "${EXTRACTS[@]}" -o "$WORK_DIR/combined-box.osm.pbf" --overwrite

# Border extracts overlap and can contain different versions of the same OSM
# object. Collapse that temporary history to the newest available version.
echo "resolving overlapping border object versions"
osmium time-filter "$WORK_DIR/combined-box.osm.pbf" -o "$WORK_DIR/current-box.osm.pbf" --overwrite
osmium check-refs "$WORK_DIR/current-box.osm.pbf"

echo "filtering routable roads"
osmium tags-filter "$WORK_DIR/current-box.osm.pbf" "$HIGHWAY_FILTER" -o "$WORK_DIR/filtered.osm.pbf" --overwrite
osmium add-locations-to-ways "$WORK_DIR/filtered.osm.pbf" -o "$WORK_DIR/map.osm.pbf" --overwrite

echo "generating Poland tiles with custom map overrides"
./mapd generate \
  --input-file "$WORK_DIR/map.osm.pbf" \
  --output-directory "$OUTPUT_DIR" \
  --generate-empty-files \
  --minlat "$MIN_LAT" \
  --minlon "$MIN_LON" \
  --maxlat "$MAX_LAT" \
  --maxlon "$MAX_LON"

tile_count="$(find "$OUTPUT_DIR" -type f | wc -l | tr -d ' ')"
tile_bytes="$(find "$OUTPUT_DIR" -type f -exec stat -f '%z' {} \; | awk '{total += $1} END {print total + 0}')"
echo "generated $tile_count Poland-area tiles ($tile_bytes bytes) in $SCRIPT_DIR/$OUTPUT_DIR"

if [[ "${MAPD_UPLOAD:-0}" == "1" ]]; then
  MAPD_OFFLINE_SOURCE="$OUTPUT_DIR" ./upload_selective_offline_comma.sh
else
  echo "not uploaded; set MAPD_UPLOAD=1 to update these tiles on the device"
fi
