# AGENTS.md

## FrogPilot mapd workflow
- This repo is the source of truth for map generation logic and offline tile contents.
- For FrogPilot runtime deployment, do not push `mapd` directly to the device from this repo.
- The runtime `mapd` binary that the device launches lives in the sibling repo at `../openpilot-rch/frogpilot/navigation/mapd`.
- `../openpilot-rch/system/manager/process_config.py` starts `mapd` from `frogpilot/navigation`.
- If runtime `mapd` behavior changes here, rebuild the ARM64 binary and update `../openpilot-rch/frogpilot/navigation/mapd`, then deploy that repo to the device using its documented workflow.
- The scripts in this repo should only rebuild the local host `mapd` used for offline map generation and upload offline tiles.
