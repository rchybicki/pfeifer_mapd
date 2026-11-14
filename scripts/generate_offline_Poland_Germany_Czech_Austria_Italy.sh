#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

SOURCES=(
  "${MAPD_POLAND_SOURCE:-poland-latest.osm.pbf}"
  "${MAPD_GERMANY_SOURCE:-germany-latest.osm.pbf}"
  "${MAPD_CZECH_SOURCE:-czech-republic-latest.osm.pbf}"
  "${MAPD_AUSTRIA_SOURCE:-austria-latest.osm.pbf}"
  "${MAPD_ITALY_SOURCE:-italy-latest.osm.pbf}"
)
SOURCE_NAMES=(poland germany czech-republic austria italy)
MERGED_SOURCE="${MAPD_MERGED_SOURCE:-poland-germany-czech-austria-italy-latest.osm.pbf}"
OUTPUT_DIR="${MAPD_OUTPUT_DIR:-offline}"
UPLOAD="${MAPD_UPLOAD:-0}"

MIN_LON=5
MIN_LAT=35
MAX_LON=25
MAX_LAT=55

for idx in "${!SOURCES[@]}"; do
  if [[ ! -f "${SOURCES[$idx]}" ]]; then
    echo "missing ${SOURCES[$idx]}; run ./download_extracts.sh ${SOURCE_NAMES[*]} first" >&2
    exit 1
  fi
done

if ! command -v osmium >/dev/null 2>&1; then
  echo "missing required tool: osmium" >&2
  echo "Install osmium-tool before generating offline map tiles." >&2
  exit 1
fi

./build_mapd.sh

echo "merging ${SOURCES[*]} -> $MERGED_SOURCE"
osmium merge "${SOURCES[@]}" -o "$MERGED_SOURCE" --overwrite

./filter_Europe.sh "$MERGED_SOURCE"
rm -rf "$OUTPUT_DIR"

for (( i = MIN_LON; i < MAX_LON; i += 20 ))
do
  for (( j = MIN_LAT; j < MAX_LAT; j += 20 ))
  do
    max_lon=$((i + 20))
    max_lat=$((j + 20))
    if (( max_lon > MAX_LON )); then
      max_lon=$MAX_LON
    fi
    if (( max_lat > MAX_LAT )); then
      max_lat=$MAX_LAT
    fi
    echo "$i $j $max_lon $max_lat"
    ./extract_box.sh "$i" "$j" "$max_lon" "$max_lat"
    ./add_locations.sh
    ./mapd generate --input-file ./map.osm.pbf --output-directory "$OUTPUT_DIR" --minlat "$j" --minlon "$i" --maxlat "$max_lat" --maxlon "$max_lon"
  done
done

if [[ "$UPLOAD" == "1" ]]; then
  ./upload_small_offline_comma.sh
else
  echo "offline maps generated in $SCRIPT_DIR/$OUTPUT_DIR"
  echo "not uploaded; set MAPD_UPLOAD=1 to push the generated tiles to the device"
fi
