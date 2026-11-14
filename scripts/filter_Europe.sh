#!/bin/bash
set -euo pipefail

INPUT_FILE="${1:-europe.osm.pbf}"

osmium tags-filter "$INPUT_FILE" "nw/highway=motorway,trunk,primary,secondary,tertiary,unclassified,residential,motorway_link,trunk_link,primary_link,secondary_link,tertiary_link" -o filtered.osm.pbf --overwrite
