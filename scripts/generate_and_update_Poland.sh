#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

SOURCE_FILE="${1:-poland-latest.osm.pbf}"

MIN_LON=14
MIN_LAT=49
MAX_LON=24
MAX_LAT=54

./filter_Europe.sh "$SOURCE_FILE"
rm -rf offline

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
    ./mapd generate --input-file ./map.osm.pbf --output-directory ./offline --minlat "$j" --minlon "$i" --maxlat "$max_lat" --maxlon "$max_lon"
    ./upload_small_offline_comma.sh
    rm -rf offline
  done
done
