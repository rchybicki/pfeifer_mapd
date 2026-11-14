#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

POLAND_SOURCE="${MAPD_POLAND_SOURCE:-poland-latest.osm.pbf}"
GERMANY_SOURCE="${MAPD_GERMANY_SOURCE:-germany-latest.osm.pbf}"
MERGED_SOURCE="${MAPD_MERGED_SOURCE:-poland-germany-latest.osm.pbf}"
OUTPUT_DIR="${MAPD_OUTPUT_DIR:-offline}"
UPLOAD="${MAPD_UPLOAD:-0}"

case "$OUTPUT_DIR" in
  "" | "." | "/" | /* | ".." | ../* | */.. | */../*)
    echo "refusing unsafe MAPD_OUTPUT_DIR: $OUTPUT_DIR" >&2
    exit 1
    ;;
esac

MIN_LON=5
MIN_LAT=47
MAX_LON=25
MAX_LAT=55

if [[ ! -f "$POLAND_SOURCE" ]]; then
  echo "missing $POLAND_SOURCE; run ./download_extracts.sh poland germany first" >&2
  exit 1
fi

if [[ ! -f "$GERMANY_SOURCE" ]]; then
  echo "missing $GERMANY_SOURCE; run ./download_extracts.sh poland germany first" >&2
  exit 1
fi

if ! command -v osmium >/dev/null 2>&1; then
  echo "missing required tool: osmium" >&2
  echo "Install osmium-tool before generating offline map tiles." >&2
  exit 1
fi

./build_mapd.sh

echo "merging $POLAND_SOURCE + $GERMANY_SOURCE -> $MERGED_SOURCE"
osmium merge "$POLAND_SOURCE" "$GERMANY_SOURCE" -o "$MERGED_SOURCE" --overwrite

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
  MAPD_OFFLINE_SOURCE="$OUTPUT_DIR" ./upload_small_offline_comma.sh
else
  echo "offline maps generated in $SCRIPT_DIR/$OUTPUT_DIR"
  echo "not uploaded; set MAPD_UPLOAD=1 to push the generated tiles to the device"
fi
