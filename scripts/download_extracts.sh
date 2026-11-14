#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

if [[ "$#" -eq 0 ]]; then
  countries=(poland germany czech-republic austria italy)
else
  countries=("$@")
fi

download_extract() {
  local country="$1"
  local output="$2"
  local url="$3"
  local tmp_output="${output}.download"

  if [[ -f "$output" && "${MAPD_FORCE_DOWNLOAD:-0}" != "1" ]]; then
    echo "using existing $output; set MAPD_FORCE_DOWNLOAD=1 to replace it"
    return
  fi

  if [[ "${MAPD_FORCE_DOWNLOAD:-0}" == "1" && "${MAPD_RESUME_DOWNLOAD:-0}" != "1" ]]; then
    rm -f "$tmp_output"
  fi

  echo "downloading $country extract -> $output"
  if command -v aria2c >/dev/null 2>&1; then
    aria2c \
      --continue=true \
      --max-connection-per-server="${MAPD_ARIA2_CONNECTIONS:-8}" \
      --split="${MAPD_ARIA2_CONNECTIONS:-8}" \
      --min-split-size=1M \
      --auto-file-renaming=false \
      --allow-overwrite=true \
      --dir . \
      --out "$tmp_output" \
      "$url"
  elif command -v wget >/dev/null 2>&1; then
    wget -c -O "$tmp_output" "$url"
  else
    curl --fail --location --continue-at - --output "$tmp_output" "$url"
  fi
  mv "$tmp_output" "$output"
}

for country in "${countries[@]}"; do
  case "$country" in
    PL | Poland | poland)
      download_extract "poland" "poland-latest.osm.pbf" "https://download.geofabrik.de/europe/poland-latest.osm.pbf"
      ;;
    DE | Germany | germany)
      download_extract "germany" "germany-latest.osm.pbf" "https://download.geofabrik.de/europe/germany-latest.osm.pbf"
      ;;
    CZ | Czechia | czechia | Czech-Republic | CzechRepublic | czech-republic | czech_republic)
      download_extract "czech-republic" "czech-republic-latest.osm.pbf" "https://download.geofabrik.de/europe/czech-republic-latest.osm.pbf"
      ;;
    AT | Austria | austria)
      download_extract "austria" "austria-latest.osm.pbf" "https://download.geofabrik.de/europe/austria-latest.osm.pbf"
      ;;
    IT | Italy | italy)
      download_extract "italy" "italy-latest.osm.pbf" "https://download.geofabrik.de/europe/italy-latest.osm.pbf"
      ;;
    SI | Slovenia | slovenia)
      download_extract "slovenia" "slovenia-latest.osm.pbf" "https://download.geofabrik.de/europe/slovenia-latest.osm.pbf"
      ;;
    Europe | europe)
      download_extract "europe" "europe-latest.osm.pbf" "https://download.geofabrik.de/europe-latest.osm.pbf"
      ;;
    *)
      echo "unknown extract '$country' (supported: poland, germany, czech-republic, austria, italy, slovenia, europe)" >&2
      exit 1
      ;;
  esac
done
