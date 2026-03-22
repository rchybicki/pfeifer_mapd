#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

wget -O europe.osm.pbf https://download.bbbike.org/osm/planet/sub-planet-daily/europe.osm.pbf
./generate_and_update_Poland.sh
