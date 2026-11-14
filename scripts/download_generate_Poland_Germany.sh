#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

MAPD_FORCE_DOWNLOAD="${MAPD_FORCE_DOWNLOAD:-1}" ./download_extracts.sh poland germany
./generate_offline_Poland_Germany.sh
