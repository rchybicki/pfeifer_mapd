#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

wget -c -O poland-latest.osm.pbf https://download.geofabrik.de/europe/poland-latest.osm.pbf
./generate_and_update_Poland.sh poland-latest.osm.pbf
